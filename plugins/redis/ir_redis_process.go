package redis

import (
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/address"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/backend"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/core/service"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/docker"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/goparser"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/workflow"
)

type RedisProcess struct {
	docker.Container
	backend.Cache

	InstanceName string
	Addr         *address.Address[*RedisProcess]
	Iface        *goparser.ParsedInterface
}

type RedisInterface struct {
	service.ServiceInterface
	Wrapped service.ServiceInterface
}

func (r *RedisInterface) GetName() string {
	return "redis(" + r.Wrapped.GetName() + ")"
}

func (r *RedisInterface) GetMethods() []service.Method {
	return r.Wrapped.GetMethods()
}

func newRedisProcess(name string, addr *address.Address[*RedisProcess]) (*RedisProcess, error) {
	proc := &RedisProcess{}
	proc.InstanceName = name
	proc.Addr = addr
	err := proc.init(name)
	if err != nil {
		return nil, err
	}
	return proc, nil
}

func (node *RedisProcess) init(name string) error {
	workflow.Init("../../runtime")

	spec, err := workflow.GetSpec()
	if err != nil {
		return err
	}

	details, err := spec.Get("RedisCache")
	if err != nil {
		return err
	}
	node.Iface = details.Iface
	return nil
}

func (r *RedisProcess) String() string {
	return r.InstanceName + " = RedisProcess(" + r.Addr.Bind.Name() + ")"
}

func (r *RedisProcess) Name() string {
	return r.InstanceName
}

func (node *RedisProcess) GetInterface(ctx blueprint.BuildContext) (service.ServiceInterface, error) {
	iface := node.Iface.ServiceInterface(ctx)
	return &RedisInterface{Wrapped: iface}, nil
}

func (r *RedisProcess) GenerateArtifacts(outputDir string) error {
	return nil
}