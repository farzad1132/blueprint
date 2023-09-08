package grpccodegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gogen"
)

type clientWrapperArgs struct {
	Builder          golang.ModuleBuilder
	PackageName      string // fully qualified package name
	PackageShortName string // package shortname
	FilePath         string
	Service          *gocode.ServiceInterface
	Name             string         // Name of the generated wrapper class
	Imports          *gogen.Imports // Manages imports for us
}

/*
It is assumed that outputPackage is the same as the one where the .proto is generated to
*/
func GenerateClientWrapper(builder golang.ModuleBuilder, service *gocode.ServiceInterface, outputPackage string) error {
	client := &clientWrapperArgs{}
	client.Builder = builder
	splits := strings.Split(outputPackage, "/")
	outputPackageName := splits[len(splits)-1]
	client.PackageName = builder.Info().Name + "/" + outputPackage
	client.PackageShortName = outputPackageName
	client.Service = service
	client.Name = service.Name + "_GRPCClientWrapper"

	outputDir := filepath.Join(builder.Info().Path, filepath.Join(splits...))
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("unable to create grpc output dir %v due to %v", outputDir, err.Error())
	}

	err = client.initImports()
	if err != nil {
		return err
	}

	outputFilename := service.Name + "_GRPCClientWrapper.go"
	return client.GenerateCode(filepath.Join(outputDir, outputFilename))
}

var requiredModules = map[string]string{
	"google.golang.org/grpc": "v1.41.0",
}
var importedPackages = []string{
	"context", "errors", "time",
	"google.golang.org/grpc",
	"google.golang.org/grpc/credentials/insecure",
}

func (client *clientWrapperArgs) importType(t gocode.TypeName) error {
	client.Imports.AddType(t)
	return client.Builder.RequireType(t)
}

func (client *clientWrapperArgs) initImports() error {
	// In addition to a few GRPC-related requirements,
	// we also depend on the modules that define the
	// argument types to the RPC methods
	client.Imports = gogen.NewImports(client.PackageName)

	for name, version := range requiredModules {
		err := client.Builder.Require(name, version)
		if err != nil {
			return err
		}
	}

	for _, pkg := range importedPackages {
		client.Imports.AddPackage(pkg)
	}

	for _, f := range client.Service.Methods {
		for _, v := range f.Arguments {
			err := client.importType(v.Type)
			if err != nil {
				return err
			}
		}
		for _, v := range f.Returns {
			err := client.importType(v.Type)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var clientWrapperTemplate = `// Blueprint: Auto-generated by GRPC Plugin
package {{.PackageShortName}}

{{.Imports}}

type {{.Name}} struct {
	Client {{.Service.Name}}Client // The GRPC-generated client
	Timeout time.Duration
}

func New_{{.Name}}(serverAddress string) (*{{.Name}}, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	duration, err := time.ParseDuration("1s")
	if err != nil {
		return nil, err
	}
	opts = append(opts, grpc.WithTimeout(duration))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		return nil, err
	}

	wrapper := &{{.Name}}{}
	wrapper.Client = New{{.Service.Name}}Client(conn)
	wrapper.Timeout = duration
	return wrapper, nil
}

{{$service := .Service.Name}}
{{$receiver := .Name}}
{{$imports := .Imports}}
{{ range $_, $f := .Service.Methods }}
func (client *{{$receiver}}) {{$f.Name}}(
	{{- range $i, $arg := $f.Arguments -}}
		{{if $i}}, {{end}}{{$arg.Name}} {{$imports.NameOf $arg.Type}}
	{{- end -}}
) (
	{{- range $i, $ret := $f.Returns -}}
	{{if $i}}, {{end}}{{$imports.NameOf $ret.Type}}
	{{- end -}}
) {
	{{$ctx := (index $f.Arguments 0).Name}}

	request := &{{$service}}_{{$f.Name}}_Request{}
	{{$ctx}}, cancel := context.WithTimeout({{$ctx}},client.Timeout)
	defer cancel()

	// TODO: arg marshalling for args 1:n

	response, err := client.Client.{{$f.Name}}({{$ctx}},request)

	{{range $i, $ret := $f.Returns -}}
	var ret{{$i}} {{$imports.NameOf $ret.Type}}
	{{end}}

	// TODO: returns marshalling

	if err == nil {
		err = ctx.Err()
	}
	return {{range $i, $ret := $f.Returns}}{{if $i}}, {{end}}ret{{$i}}{{end}}
}
{{end}}
`

/*
Generates the file within its module
*/
func (client *clientWrapperArgs) GenerateCode(outputFilePath string) error {
	t, err := template.New("GRPCClientWrapperTemplate").Parse(clientWrapperTemplate)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(outputFilePath, os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	return t.Execute(f, client)
}
