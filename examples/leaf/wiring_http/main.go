package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/circuitbreaker"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/goproc"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/healthchecker"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/http"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/retries"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/simplecache"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/simplenosqldb"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/workflow"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/xtrace"
)

func serviceDefaults(wiring blueprint.WiringSpec, serviceName string) string {
	procName := fmt.Sprintf("p%s", serviceName)
	retries.AddRetries(wiring, serviceName, 10)
	healthchecker.AddHealthCheckAPI(wiring, serviceName)
	circuitbreaker.AddCircuitBreaker(wiring, serviceName, 1000, 0.1, "1s")
	xtrace.Instrument(wiring, serviceName)
	http.Deploy(wiring, serviceName)
	return goproc.CreateProcess(wiring, procName, serviceName)
}

func main() {
	slog.Info("Constructing Wiring Spec")

	wiring := blueprint.NewWiringSpec("leaf_example")

	workflow.Init("../workflow")

	b_database := simplenosqldb.Define(wiring, "b_database")
	b_cache := simplecache.Define(wiring, "b_cache")
	b := workflow.Define(wiring, "b", "LeafServiceImpl", b_cache, b_database)

	a := workflow.Define(wiring, "a", "NonLeafService", b)
	pa := serviceDefaults(wiring, a)
	pb := serviceDefaults(wiring, b)

	slog.Info("Wiring Spec: \n" + wiring.String())

	bp, err := wiring.GetBlueprint()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	bp.Instantiate(pa, pb)

	application, err := bp.Build()
	if err != nil {
		slog.Error("Unable to build blueprint, exiting", "error", err)
		slog.Info("Application: \n" + application.String())
		os.Exit(1)
	}

	slog.Info("Application: \n" + application.String())

	nodes := make(map[string]blueprint.IRNode)
	for _, node := range application.Children {
		nodes[node.Name()] = node
	}

	err = nodes["pb"].(*goproc.Process).GenerateArtifacts("tmp")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	err = nodes["pa"].(*goproc.Process).GenerateArtifacts("tmp")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Exiting")
}
