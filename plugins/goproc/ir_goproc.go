package goproc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/irutil"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/process"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gogen"
	"golang.org/x/exp/slog"
)

/*
This file contains the implementation of the golang.Process IRNode.

The `GenerateArtifacts` method generates the main method based on the process's contained nodes.

Most of the heavy lifting of code generation is done by the following:
* gogen/workspacebuilder
* gogen/modulebuilder
* gogen/graphbuilder

*/

var generatedModulePrefix = "gitlab.mpi-sws.org/cld/blueprint/plugins/golang/process"

// An IRNode representing a golang process.
// This is Blueprint's main implementation of Golang processes
type Process struct {
	blueprint.IRNode
	process.ProcessNode
	process.ArtifactGenerator

	InstanceName   string
	ArgNodes       []blueprint.IRNode
	ContainedNodes []blueprint.IRNode
}

// A Golang Process Node can either be given the child nodes ahead of time, or they can be added using AddArtifactNode / AddCodeNode
func newGolangProcessNode(name string) *Process {
	node := Process{}
	node.InstanceName = name
	return &node
}

func (node *Process) Name() string {
	return node.InstanceName
}

func (node *Process) String() string {
	var b strings.Builder
	b.WriteString(node.InstanceName)
	b.WriteString(" = GolangProcessNode(")
	var args []string
	for _, arg := range node.ArgNodes {
		args = append(args, arg.Name())
	}
	b.WriteString(strings.Join(args, ", "))
	b.WriteString(") {\n")
	var children []string
	for _, child := range node.ContainedNodes {
		children = append(children, child.String())
	}
	b.WriteString(blueprint.Indent(strings.Join(children, "\n"), 2))
	b.WriteString("\n}")
	return b.String()
}

func (node *Process) AddArg(argnode blueprint.IRNode) {
	node.ArgNodes = append(node.ArgNodes, argnode)
}

func (node *Process) AddChild(child blueprint.IRNode) error {
	node.ContainedNodes = append(node.ContainedNodes, child)
	return nil
}

type mainArg struct {
	Name string
	Doc  string
	Var  string
}

type mainTemplateArgs struct {
	GraphPackage     string
	GraphConstructor string
	Args             []mainArg
	Instantiate      []string
}

var mainTemplate = `// This file is auto-generated by the Blueprint goproc plugin
package main

import (
	"flag"
	"os"
	"golang.org/x/exp/slog"
	"context"
	"{{.GraphPackage}}"
)

func checkArg(name, value string) {
	if value == "" {
		slog.Error("No value set for required cmd line argument " + name)
		os.Exit(1)
	} else {
		slog.Info(fmt.Sprintf("Arg %v = %v", name, value))
	}
}

func main() {
	{{- range $i, $arg := .Args}}
	{{$arg.Var}} := flag.String("{{$arg.Name}}", "", "Argument automatically generated from Blueprint IR: {{$arg.Doc}}")
	{{end}}

	flag.Parse()

	{{range $i, $arg := .Args -}}
	checkArg("{{$arg.Name}}", *{{$arg.Var}})
	{{end}}
	
	graphArgs := map[string]string{
		{{- range $i, $arg := .Args}}
		"{{$arg.Name}}": *{{$arg.Var}},
		{{- end}}
	}

	ctx, cancel := context.WithCancel(context.Background())
	graph, err := {{.GraphConstructor}}(ctx, cancel, graphArgs)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	{{range $i, $node := .Instantiate -}}
	_, err = graph.Get("{{$node}}")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	{{end}}

	graph.WaitGroup().Wait()
}`

func (node *Process) GenerateArtifacts(outputDir string) error {
	err := gogen.CheckDir(outputDir, true)
	if err != nil {
		return fmt.Errorf("unable to create %s for process %s due to %s", outputDir, node.Name(), err.Error())
	}

	// TODO: might end up building multiple times which is OK, so need a check here that we haven't already built this artifact, even if it was by a different (but identical) node

	slog.Info(fmt.Sprintf("Building %s to %s\n", node.Name(), outputDir))

	// Generate the workspace and copy all local artifacts
	cleanName := irutil.Clean(node.Name())
	workspaceDir := filepath.Join(outputDir, cleanName)
	workspace, err := gogen.NewWorkspaceBuilder(workspaceDir)
	if err != nil {
		return err
	}

	err = workspace.Visit(node.ContainedNodes)
	if err != nil {
		return err
	}

	// Generate the module and add all dependencies
	moduleName := generatedModulePrefix + "/" + cleanName
	module, err := gogen.NewModuleBuilder(workspace, cleanName, moduleName)
	if err != nil {
		return err
	}
	module.Require("golang.org/x/exp", "v0.0.0-20230728194245-b0cb94b80691")

	err = module.Visit(node.ContainedNodes)
	if err != nil {
		return err
	}

	// Generate the graph of gonodes contained in this process
	graphFileName := strings.ToLower(cleanName) + ".go"
	packagePath := "goproc"
	constructorName := "New" + strings.ToTitle(cleanName)

	graph, err := gogen.NewGraphBuilder(module, graphFileName, packagePath, constructorName)
	if err != nil {
		return err
	}

	err = graph.Visit(node.ContainedNodes)
	if err != nil {
		return err
	}

	// Generate the main.go
	mainFileArgs := mainTemplateArgs{
		GraphPackage:     fmt.Sprintf("%s/%s", module.Name, packagePath),
		GraphConstructor: fmt.Sprintf("%s.%s", packagePath, constructorName),
	}
	for _, arg := range node.ArgNodes {
		mainFileArgs.Args = append(mainFileArgs.Args, mainArg{
			Name: arg.Name(),
			Doc:  arg.String(),
			Var:  irutil.Clean(arg.Name()),
		})
	}
	// For now explicitly instantiate every child node
	for _, child := range node.ContainedNodes {
		mainFileArgs.Instantiate = append(mainFileArgs.Instantiate, child.Name())
	}
	mainFileName := filepath.Join(module.ModuleDir, "main.go")
	t, err := template.New("main.go").Parse(mainTemplate)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(mainFileName, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	err = t.Execute(f, mainFileArgs)
	if err != nil {
		return err
	}

	// Build workspace, module, and graph
	err = graph.Build()
	if err != nil {
		return err
	}

	err = module.Finish()
	if err != nil {
		return err
	}

	err = workspace.Finish()
	if err != nil {
		return err
	}

	return nil
}