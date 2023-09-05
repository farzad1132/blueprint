package grpc

import (
	"fmt"
	"reflect"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/service"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang"
)

// IR node for the GRPC server.  Once a service is exposed over GRPC,
// it no longer has an interface that is callable by other golang instances,
// so it is not a golang.Service node any more.  However, it is still a service,
// but now one exposed over GRPC.
type GolangServer struct {
	service.ServiceNode

	InstanceName string
	Addr         *GolangServerAddress
	Wrapped      golang.Service
}

// Represents a service that is exposed over GRPC
type GRPCInterface struct {
	service.ServiceInterface
	Wrapped service.ServiceInterface
}

func (grpc *GRPCInterface) GetName() string {
	return "grpc(" + grpc.Wrapped.GetName() + ")"
}

func (grpc *GRPCInterface) GetMethods() []service.Method {
	return grpc.Wrapped.GetMethods()
}

func newGolangServer(name string, serverAddr blueprint.IRNode, wrapped blueprint.IRNode) (*GolangServer, error) {
	addr, is_addr := serverAddr.(*GolangServerAddress)
	if !is_addr {
		return nil, fmt.Errorf("GRPC server %s expected %s to be an address, but got %s", name, serverAddr.Name(), reflect.TypeOf(serverAddr).String())
	}

	service, is_service := wrapped.(golang.Service)
	if !is_service {
		return nil, fmt.Errorf("GRPC server %s expected %s to be a golang service, but got %s", name, wrapped.Name(), reflect.TypeOf(wrapped).String())
	}

	node := &GolangServer{}
	node.InstanceName = name
	node.Addr = addr
	node.Wrapped = service
	return node, nil
}

func (n *GolangServer) String() string {
	return n.InstanceName + " = GRPCServer(" + n.Wrapped.Name() + ", " + n.Addr.Name() + ")"
}

func (n *GolangServer) Name() string {
	return n.InstanceName
}

func (node *GolangServer) AddInstantiation(builder golang.DICodeBuilder) error {
	// TODO
	return nil
}

func (node *GolangServer) GetInterface() service.ServiceInterface {
	return &GRPCInterface{Wrapped: node.Wrapped.GetInterface()}
}
func (node *GolangServer) ImplementsGolangNode() {}
