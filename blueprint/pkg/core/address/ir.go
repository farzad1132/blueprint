package address

import (
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
)

type Address interface {
	blueprint.IRNode
	Name() string
	GetDestination() blueprint.IRNode
	SetDestination(blueprint.IRNode) error
	ImplementsAddressNode()
}
