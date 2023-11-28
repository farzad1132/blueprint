package loadbalancer

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint/stringutil"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/coreplugins/service"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/ir"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gogen"
)

type LoadBalancerClient struct {
	golang.Service
	golang.GeneratesFuncs

	BalancerName   string
	Clients        []golang.Service
	ContainedNodes []ir.IRNode

	outputPackage string
}

func newLoadBalancerClient(name string, arg_nodes []ir.IRNode) (*LoadBalancerClient, error) {
	return &LoadBalancerClient{
		BalancerName:   name,
		ContainedNodes: arg_nodes,
		outputPackage:  "lb",
	}, nil
}

func (node *LoadBalancerClient) Name() string {
	return node.BalancerName
}

func (node *LoadBalancerClient) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%v = LoadBalancer() {\n", node.BalancerName))
	var children []string
	for _, child := range node.ContainedNodes {
		children = append(children, child.String())
	}
	b.WriteString(stringutil.Indent(strings.Join(children, "\n"), 2))
	b.WriteString("\n}")
	return b.String()
}

func (lb *LoadBalancerClient) GetInterface(ctx ir.BuildContext) (service.ServiceInterface, error) {
	// LoadBalancer doesn't modify the interface! As all clients must have the same interface, we can simply return the interface of any of the clients.
	return lb.Clients[0].GetInterface(ctx)
}

func (lb *LoadBalancerClient) AddInterfaces(module golang.ModuleBuilder) error {
	for _, node := range lb.ContainedNodes {
		if n, valid := node.(golang.ProvidesInterface); valid {
			if err := n.AddInterfaces(module); err != nil {
				return err
			}
		}
	}
	return nil
}

func (lb *LoadBalancerClient) AddInstantiation(builder golang.NamespaceBuilder) error {
	if builder.Visited(lb.BalancerName) {
		return nil
	}

	iface, err := golang.GetGoInterface(builder, lb.Clients[0])
	if err != nil {
		return err
	}

	constructor := &gocode.Constructor{
		Package: builder.Module().Info().Name + "/" + lb.outputPackage,
		Func: gocode.Func{
			Name: fmt.Sprintf("New_%v_LoadBalancer", iface.BaseName),
			Arguments: []gocode.Variable{
				{Name: "ctx", Type: &gocode.UserType{Package: "context", Name: "Context"}},
				{Name: "clients", Type: &gocode.Slice{SliceOf: iface}},
			},
		},
	}

	return builder.DeclareConstructor(lb.BalancerName, constructor, lb.ContainedNodes)
}

func (lb *LoadBalancerClient) GenerateFuncs(module golang.ModuleBuilder) error {
	if module.Visited(lb.BalancerName) {
		return nil
	}

	pkg, err := module.CreatePackage(lb.outputPackage)
	if err != nil {
		return err
	}

	iface, err := golang.GetGoInterface(module, lb.Clients[0])
	if err != nil {
		return err
	}

	args := &clientArgs{}
	args.LBName = lb.BalancerName
	args.PackageShortName = lb.outputPackage
	args.Imports = gogen.NewImports(pkg.Name)
	args.ServiceName = iface.BaseName
	args.Service = iface
	lbFileName := filepath.Join(module.Info().Path, args.PackageShortName)
	args.Imports.AddPackages("context", "gitlab.mpi-sws.org/cld/blueprint/runtime/plugins/loadbalancer")

	return gogen.ExecuteTemplateToFile("lb_client_constructor", lbTemplate, args, lbFileName)
}

type clientArgs struct {
	LBName           string
	ServiceName      string
	PackageShortName string
	Imports          *gogen.Imports
	Service          *gocode.ServiceInterface
}

var lbTemplate = `// This file is auto-generated by the Blueprint loadbalancer plugin
package {{.PackageShortName}}

{{.Imports}}

type {{.LBName}} struct {
	balancer *loadbalancer.LoadBalancer[{{NameOf .Service.UserType}}]
}

func New_{{.ServiceName}}_LoadBalancer(ctx context.Context, clients []{{NameOf .Service.UserType}}) (*{{.LBName}}, error) {
	handler := &{{.Name}}{}
	handler.balancer = loadbalancer.NewLoadBalancer[{{NameOf .Service.UserType}}](ctx, clients)
	return handler, nil
}

{{$service := .Service -}}
{{$receiver := .LBName -}}
{{ range $_, $f := .Service.Methods }}
func (lbalancer *{{$receiver}}) {{SignatureWithRetVars $f}} {
	client := lbalancer.balancer.PickClient(ctx)
	return client.{{$f.Name}}({{ArgVars $f "ctx"}})
}
{{end}}
`
