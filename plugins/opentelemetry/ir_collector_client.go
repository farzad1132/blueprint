package opentelemetry

import (
	"fmt"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/coreplugins/address"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/coreplugins/service"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/ir"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/goparser"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/workflow"
	"golang.org/x/exp/slog"
)

type OpenTelemetryCollectorInterface interface {
	golang.Node
	golang.Instantiable
	ImplementsOTCollectorClient()
}

type OpenTelemetryCollectorClient struct {
	golang.Node
	golang.Instantiable

	ClientName string
	ServerDial *address.DialConfig

	InstanceName string
	Iface        *goparser.ParsedInterface
	Constructor  *gocode.Constructor
}

func newOpenTelemetryCollectorClient(name string, addr *address.DialConfig) (*OpenTelemetryCollectorClient, error) {
	node := &OpenTelemetryCollectorClient{}
	err := node.init(name)
	if err != nil {
		return nil, err
	}
	node.ClientName = name
	node.ServerDial = addr
	return node, nil
}

func (node *OpenTelemetryCollectorClient) Name() string {
	return node.ClientName
}

func (node *OpenTelemetryCollectorClient) String() string {
	return node.Name() + " = OTClient(" + node.ServerDial.Name() + ")"
}

func (node *OpenTelemetryCollectorClient) init(name string) error {
	workflow.Init("../../runtime")

	spec, err := workflow.GetSpec()
	if err != nil {
		return err
	}

	details, err := spec.Get("StdoutTracer")
	if err != nil {
		return err
	}

	node.InstanceName = name
	node.Iface = details.Iface
	node.Constructor = details.Constructor.AsConstructor()
	return nil
}

func (node *OpenTelemetryCollectorClient) AddInstantiation(builder golang.NamespaceBuilder) error {
	// Only generate instantiation code for this instance once
	if builder.Visited(node.ClientName) {
		return nil
	}

	slog.Info(fmt.Sprintf("Instantiating OTCollectorClient %v in %v/%v", node.InstanceName, builder.Info().Package.PackageName, builder.Info().FileName))

	return builder.DeclareConstructor(node.InstanceName, node.Constructor, []ir.IRNode{node.ServerDial})
}

func (node *OpenTelemetryCollectorClient) GetInterface(ctx ir.BuildContext) (service.ServiceInterface, error) {
	return node.Iface.ServiceInterface(ctx), nil
}

func (node *OpenTelemetryCollectorClient) AddInterfaces(builder golang.ModuleBuilder) error {
	return node.AddToWorkspace(builder.Workspace())
}

func (node *OpenTelemetryCollectorClient) AddToWorkspace(builder golang.WorkspaceBuilder) error {
	return golang.AddRuntimeModule(builder)
}

func (node *OpenTelemetryCollectorClient) ImplementsGolangNode() {}

func (node *OpenTelemetryCollectorClient) ImplementsOTCollectorClient() {}
