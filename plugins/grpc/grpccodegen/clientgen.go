package grpccodegen

import (
	"fmt"
	"path/filepath"

	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/golang/gocode"
	"github.com/blueprint-uservices/blueprint/plugins/golang/gogen"
	"golang.org/x/exp/slog"
)

// Generates a gRPC client for the specified service
func GenerateClient(builder golang.ModuleBuilder, service *gocode.ServiceInterface, outputPackage string) error {
	pkg, err := builder.CreatePackage(outputPackage)
	if err != nil {
		return err
	}

	client := &clientArgs{
		Package: pkg,
		Service: service,
		Name:    service.BaseName + "_GRPCClient",
		Imports: gogen.NewImports(pkg.Name),
	}

	client.Imports.AddPackages(
		"context", "time",
		"google.golang.org/grpc",
		"google.golang.org/grpc/credentials/insecure",
	)

	slog.Info(fmt.Sprintf("Generating %v/%v.go", client.Package.PackageName, client.Name))
	outputFile := filepath.Join(client.Package.Path, client.Name+".go")
	return gogen.ExecuteTemplateToFile("GRPCClient", clientTemplate, client, outputFile)
}

/*
Arguments to the template code
*/
type clientArgs struct {
	Package golang.PackageInfo
	Service *gocode.ServiceInterface
	Name    string         // Name of the generated client class
	Imports *gogen.Imports // Manages imports for us
}

var clientTemplate = `// Blueprint: Auto-generated by GRPC Plugin
package {{.Package.ShortName}}

{{.Imports}}

type {{.Name}} struct {
	{{.Imports.NameOf .Service.UserType}}
	Client {{.Service.Name}}Client // The actual GRPC-generated client
	Timeout time.Duration
}

func New_{{.Name}}(ctx context.Context, serverAddress string) (*{{.Name}}, error) {
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

	c := &{{.Name}}{}
	c.Client = New{{.Service.Name}}Client(conn)
	c.Timeout = duration
	return c, nil
}

{{$service := .Service.Name -}}
{{$receiver := .Name -}}
{{- range $_, $f := .Service.Methods }}
func (client *{{$receiver}}) {{SignatureWithRetVars $f}} {
	// Create and marshall the GRPC Request object
	req := &{{$service}}_{{$f.Name}}_Request{}
	req.marshall({{ArgVars $f}})

	// Configure the client-side request timeout
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()

	// Make the remote call
	rsp, err := client.Client.{{$f.Name}}(ctx, req)
	if err == nil {
		err = ctx.Err()
	}
	if err != nil {
		return
	}

	{{RetVarsEquals $f}} rsp.unmarshall()
	return
}
{{end}}
`
