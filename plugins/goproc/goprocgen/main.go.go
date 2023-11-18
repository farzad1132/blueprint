package goprocgen

import (
	"fmt"
	"path/filepath"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/ir"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gogen"
	"golang.org/x/exp/slog"
)

/*
Generates a main.go file in the provided module.  The main method will
call the graphConstructor provided to create and instantiate nodes.
*/
func GenerateMain(
	name string,
	argNodes []ir.IRNode,
	nodesToInstantiate []ir.IRNode,
	module golang.ModuleBuilder,
	graphPackage string,
	graphConstructor string) error {

	// Generate the main.go
	mainArgs := mainTemplateArgs{
		Name:             name,
		GraphPackage:     graphPackage,
		GraphConstructor: graphConstructor,
		Args:             nil,
		Config:           make(map[string]string),
		Instantiate:      nil,
	}

	// Expect command-line arguments for all argNodes specified
	for _, arg := range argNodes {
		mainArgs.Args = append(mainArgs.Args, mainArg{
			Name: arg.Name(),
			Doc:  arg.String(),
			Var:  ir.CleanName(arg.Name()),
		})
	}

	// Instantiate the nodes specified
	for _, node := range ir.FilterNodes[golang.Instantiable](nodesToInstantiate) {
		mainArgs.Instantiate = append(mainArgs.Instantiate, node.Name())
	}

	// Materialize any configuration
	for _, node := range ir.Filter[ir.IRConfig](nodesToInstantiate) {
		if !node.HasValue() {
			return blueprint.Errorf("golang main method expects to instantiate config variable %v but no value is set for it", node.Name())
		}
		mainArgs.Config[node.Name()] = node.Value()
	}

	slog.Info(fmt.Sprintf("Generating %v/main.go", module.Info().Name))
	mainFileName := filepath.Join(module.Info().Path, "main.go")
	return gogen.ExecuteTemplateToFile("goprocMain", mainTemplate, mainArgs, mainFileName)
}

type mainArg struct {
	Name string
	Doc  string
	Var  string
}

type mainTemplateArgs struct {
	Name             string
	GraphPackage     string
	GraphConstructor string
	Args             []mainArg
	Config           map[string]string
	Instantiate      []string
}

var mainTemplate = `// {{.Name}} runs the {{.Name}} Golang process.
//
// {{.Name}} is auto-generated by Blueprint's goproc plugin (goproc/goprocgen/main.go.go)
//
// Usage:
//
//   go run main.go {{range $_, $arg := .Args}}--{{$arg.Name}}=value {{end}}
//
// {{.Name}} requires the following arguments are passed:
{{- range $_, $arg := .Args }}
//
//   --{{$arg.Name}}
//       Auto-generated by Blueprint IR node:
//       {{$arg.Doc}}
{{- end }}
//
// {{.Name}} will instantiate the following IR nodes:
{{- range $_, $name := .Instantiate }}
//   {{$name}}
{{- end }}
package main

import (
	"context"
	"os"

	"{{.GraphPackage}}"

	"golang.org/x/exp/slog"
)

func main() {
	slog.Info("Running {{.Name}}")
	n, err := {{.GraphConstructor}}("{{.Name}}").Build(context.Background())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	n.Await()
	slog.Info("{{.Name}} exiting")
}`
