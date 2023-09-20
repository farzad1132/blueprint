package opentelemetry

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/service"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"golang.org/x/exp/slog"
)

type OpenTelemetryServerWrapper struct {
	golang.Service
	golang.Instantiable
	golang.GeneratesInterfaces
	golang.GeneratesFuncs

	WrapperName string
	Wrapped     golang.Service
	Collector   *OpenTelemetryCollectorClient
}

func newOpenTelemetryServerWrapper(name string, server blueprint.IRNode, collector blueprint.IRNode) (*OpenTelemetryServerWrapper, error) {
	serverNode, is_callable := server.(golang.Service)
	if !is_callable {
		return nil, blueprint.Errorf("opentelemetry server wrapper requires %s to be a golang service but got %s", server.Name(), reflect.TypeOf(server).String())
	}

	collectorClient, is_collector_client := collector.(*OpenTelemetryCollectorClient)
	if !is_collector_client {
		return nil, blueprint.Errorf("opentelemetry server wrapper requires %s to be an opentelemetry collector client", collector.Name())
	}

	node := &OpenTelemetryServerWrapper{}
	node.WrapperName = name
	node.Wrapped = serverNode
	node.Collector = collectorClient
	return node, nil
}

func (node *OpenTelemetryServerWrapper) Name() string {
	return node.WrapperName
}

func (node *OpenTelemetryServerWrapper) String() string {
	return node.Name() + " = OTServerWrapper(" + node.Wrapped.Name() + ", " + node.Collector.Name() + ")"
}

func (node *OpenTelemetryServerWrapper) GetInterface() service.ServiceInterface {
	// TODO: extend wrapped interface with tracing stuff
	return node.Wrapped.GetInterface()
}

func (n *OpenTelemetryServerWrapper) GetGoInterface() *gocode.ServiceInterface {
	// TODO: return memcached interface
	return nil
}

// Part of code generation compilation pass; creates the interface definition code for the wrapper,
// and any new generated structs that are exposed and can be used by other IRNodes
func (node *OpenTelemetryServerWrapper) GenerateInterfaces(builder golang.ModuleBuilder) error {
	// Only generate instantiation code for this instance once
	if builder.Visited(node.WrapperName + ".GenerateInterfaces") {
		return nil
	}
	slog.Info(fmt.Sprintf("GenerateInterfaces %v\n", node))

	// TODO: Generate the extended service interface that includes extra arguments and any structs that are used in that interface

	return nil
}

// Part of code generation compilation pass; provides implementation of interfaces from GenerateInterfaces
func (node *OpenTelemetryServerWrapper) GenerateFuncs(builder golang.ModuleBuilder) error {
	// Only generate instantiation code for this instance once
	if builder.Visited(node.WrapperName + ".GenerateFuncs") {
		return nil
	}
	slog.Info(fmt.Sprintf("GenerateFuncs %v\n", node))

	// TODO: Generate the wrapper implementation

	return nil
}

var serverBuildFuncTemplate = `func(ctr golang.Container) (any, error) {

		// TODO: generated OT server constructor

		return nil, nil

	}`

// Part of code generation compilation pass; provides instantiation snippet
func (node *OpenTelemetryServerWrapper) AddInstantiation(builder golang.GraphBuilder) error {
	// Only generate instantiation code for this instance once
	if builder.Visited(node.WrapperName) {
		return nil
	}

	// TODO: generate the OT wrapper instantiation

	// Instantiate the code template
	t, err := template.New(node.WrapperName).Parse(serverBuildFuncTemplate)
	if err != nil {
		return err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	err = t.Execute(buf, node)
	if err != nil {
		return err
	}

	return builder.Declare(node.WrapperName, buf.String())
}

func (node *OpenTelemetryServerWrapper) ImplementsGolangNode()    {}
func (node *OpenTelemetryServerWrapper) ImplementsGolangService() {}
