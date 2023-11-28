package loadbalancer

import (
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/coreplugins/pointer"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/ir"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/wiring"
)

// Creates a client-side load-balancer for multiple instances of a service. The list of services must be provided as an argument at compile-time when using this plugin.
func Create(spec wiring.WiringSpec, services []string, serviceType string) string {
	loadbalancer_name := serviceType + ".lb"
	spec.Define(loadbalancer_name, &LoadBalancerClient{}, func(namespace wiring.Namespace) (ir.IRNode, error) {
		var arg_nodes []ir.IRNode
		for _, arg_name := range services {
			var arg ir.IRNode
			if err := namespace.Get(arg_name, &arg); err != nil {
				return nil, err
			}
			arg_nodes = append(arg_nodes, arg)
		}

		return newLoadBalancerClient(serviceType, arg_nodes)
	})

	dstName := loadbalancer_name + ".dst"
	spec.Alias(dstName, loadbalancer_name)

	pointer.CreatePointer(spec, loadbalancer_name, &LoadBalancerClient{}, dstName)
	return loadbalancer_name
}
