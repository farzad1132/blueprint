package httpcodegen

import (
	"fmt"
	"path/filepath"

	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gogen"
	"golang.org/x/exp/slog"
)

/*
This function is used by the HTTP plugin to generate the server-side HTTP service.
*/
func GenerateServerHandler(builder golang.ModuleBuilder, service *gocode.ServiceInterface, outputPackage string) error {
	// TODO: Implement
	pkg, err := builder.CreatePackage(outputPackage)
	if err != nil {
		return err
	}

	server := &serverArgs{
		Package: pkg,
		Service: service,
		Name:    service.BaseName + "_HTTPServerHandler",
		Imports: gogen.NewImports(pkg.Name),
	}

	server.Imports.AddPackages("context", "encoding/json", "net/http", "github.com/gorilla/mux")

	slog.Info(fmt.Sprintf("Generating %v/%v_HTTPServer.go", server.Package.PackageName, service.Name))
	outputFile := filepath.Join(server.Package.Path, service.Name+"_HTTPServer.go")
	return gogen.ExecuteTemplateToFile("HTTPServer", serverTemplate, server, outputFile)
}

/*
Arguments to the template code
*/
type serverArgs struct {
	Package golang.PackageInfo
	Service *gocode.ServiceInterface
	Name    string         // Name of the generated wrapper class
	Imports *gogen.Imports // Manages imports for us
}

var serverTemplate = `// Blueprint: Auto-generated by HTTP Plugin
package {{.Package.ShortName}}

{{.Imports}}

type {{.Name}} struct {
	Service {{.Imports.NameOf .Service.UserType}}
	Address string
}

func New_{{.Name}}(ctx context.Context, service {{.Imports.NameOf .Service.UserType}}, serverAddress string) (*{{.Name}}, error) {
	handler := &{{.Name}}{}
	handler.Service = service
	handler.Address = serverAddress
	return handler, nil
}

// Blueprint: Run is called automatically in a separate goroutine by runtime/plugins/golang/di.go
func (handler *{{.Name}}) Run(ctx context.Context) error {
	router := mux.NewRouter()
	// Add paths for the mux router
	{{ range $_, $f := .Service.Methods }}
	router.Path("/{{$f.Name}}").HandlerFunc(handler.{{$f.Name}})
	{{end}}
	srv := &http.Server {
		Addr: handler.Address,
		Handler: router,
	}

	go func() {
		select {
		case <-ctx.Done():
			srv.Shutdown(ctx)
		}
	}()

	return srv.ListenAndServe()
}

{{$service := .Service.Name -}}
{{$receiver := .Name -}}
{{ range $_, $f := .Service.Methods }}
func (handler *{{$receiver}}) {{$f.Name -}}
	(w http.ResponseWriter, r *http.Request) {
	var err error
	defer r.Body.Close()
	{{range $_, $arg := $f.Arguments}}
	request_{{$arg.Name}} := r.URL.Query().Get("{{$arg.Name}}")
	var {{$arg.Name}} {{NameOf $arg.Type}}
	err = json.Unmarshal([]byte(request_{{$arg.Name}}), &{{$arg.Name}})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	{{end}}
	ctx := context.Background()
	{{RetVars $f "err"}} {{HasNewReturnVars $f}} handler.Service.{{$f.Name}}({{ArgVars $f "ctx"}})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response := struct {
		{{range $i, $arg := $f.Returns}}
		Ret{{$i}} {{NameOf $arg.Type}}
		{{end}}
	}{}
	{{range $i, $arg := $f.Returns}}
	response.Ret{{$i}} = ret{{$i}}
	{{end}}
	json.NewEncoder(w).Encode(response)
}
{{end}}
`
