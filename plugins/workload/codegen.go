package workload

import (
	"path/filepath"

	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gogen"
)

// Generates the workload generator client
func GenerateWorkloadgenCode(builder golang.ModuleBuilder, service *gocode.ServiceInterface, outputPackage string) error {
	pkg, err := builder.CreatePackage(outputPackage)
	if err != nil {
		return err
	}

	wlgen := &wlgenArgs{
		Name:    service.Name + "_WorkloadGenerator",
		Package: pkg,
		Service: service,
		Imports: gogen.NewImports(pkg.Name),
	}

	wlgen.Imports.AddPackages(
		"context", "fmt", "strings",
	)

	outputFile := filepath.Join(wlgen.Package.Path, service.Name+"_workloadgen.go")
	return gogen.ExecuteTemplateToFile("workloadgen", workloadClientTemplate, wlgen, outputFile)
}

type wlgenArgs struct {
	Name    string // Name of the generated workloadgen struct
	Package golang.PackageInfo
	Service *gocode.ServiceInterface
	Imports *gogen.Imports
}

var workloadClientTemplate = `// Blueprint: Auto-generated by Workload Plugin
package {{.Package.ShortName}}

{{.Imports}}

type {{.Name}} struct {
	Service {{NameOf .Service.UserType}}
}

func New_{{.Name}}(service {{NameOf .Service.UserType}}) (*{{.Name}}, error) {
	wlgen := &{{.Name}}{}
	wlgen.Service = service
	return wlgen, nil
}

// Blueprint: Run is called automatically in a separate goroutine by runtime/plugins/golang/di.go
func (wlgen *{{.Name}}) Run(ctx context.Context) error {
	{{ range $_, $f := .Service.Methods -}}
	err := wlgen.Call_{{$f.Name}}(ctx)
	if err != nil {
		return err
	}
	{{end}}

	return nil
}

// Utility method for printing values of args and retvals
func toString(values ...any) string {
	var s []string
	for _, v := range values {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return strings.Join(s, ", ")
}


{{$service := .Service -}}
{{$receiver := .Name -}}
{{$imports := .Imports -}}
{{ range $_, $f := .Service.Methods }}
func (wlgen *{{$receiver}}) Call_{{$f.Name}}(ctx context.Context) error {
	{{DeclareArgVars $f}}
	fmt.Printf("{{$service.UserType.Name}}.{{$f.Name}}(%v)", toString({{ArgVars $f}}))
	{{RetVars $f "err"}} :=  wlgen.Service.{{$f.Name}}({{ArgVars $f "ctx"}})
	if err != nil {
		fmt.Printf(" = error\n")
		return err
	} else {
		fmt.Printf(" = %v\n", toString({{RetVars $f}}))
	}
	return nil
}
{{end}}



`
