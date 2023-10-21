package zipkin

import (
	"fmt"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/address"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/service"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/goparser"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/workflow"
	"golang.org/x/exp/slog"
)

type ZipkinCollectorClient struct {
	golang.Node
	golang.Instantiable
	ClientName string
	ServerDial *address.DialConfig

	InstanceName string
	Iface        *goparser.ParsedInterface
	Constructor  *gocode.Constructor
}

func newZipkinCollectorClient(name string, addr *address.DialConfig) (*ZipkinCollectorClient, error) {
	node := &ZipkinCollectorClient{}
	err := node.init(name)
	if err != nil {
		return nil, err
	}
	node.ClientName = name
	node.ServerDial = addr
	return node, nil
}

func (node *ZipkinCollectorClient) Name() string {
	return node.ClientName
}

func (node *ZipkinCollectorClient) String() string {
	return node.Name() + " = ZipkinClient(" + node.ServerDial.Name() + ")"
}

func (node *ZipkinCollectorClient) init(name string) error {
	workflow.Init("../../runtime")

	spec, err := workflow.GetSpec()
	if err != nil {
		return err
	}

	details, err := spec.Get("ZipkinTracer")
	if err != nil {
		return err
	}

	node.InstanceName = name
	node.Iface = details.Iface
	node.Constructor = details.Constructor.AsConstructor()
	return nil
}

func (node *ZipkinCollectorClient) AddInstantiation(builder golang.GraphBuilder) error {
	// Only generate instantiation code for this instance once
	if builder.Visited(node.ClientName) {
		return nil
	}

	slog.Info(fmt.Sprintf("Instantiating ZipkinClient %v in %v/%v", node.InstanceName, builder.Info().Package.PackageName, builder.Info().FileName))

	return builder.DeclareConstructor(node.InstanceName, node.Constructor, []blueprint.IRNode{node.ServerDial})
}

func (node *ZipkinCollectorClient) GetInterface(ctx blueprint.BuildContext) (service.ServiceInterface, error) {
	return node.Iface.ServiceInterface(ctx), nil
}

func (node *ZipkinCollectorClient) AddInterfaces(builder golang.WorkspaceBuilder) error {
	return golang.AddRuntimeModule(builder)
}

func (node *ZipkinCollectorClient) AddToWorkspace(builder golang.WorkspaceBuilder) error {
	return golang.AddRuntimeModule(builder)
}

func (node *ZipkinCollectorClient) ImplementsGolangNode() {}

func (node *ZipkinCollectorClient) ImplementsOTCollectorClient() {}
