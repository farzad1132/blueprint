package wiring

import (
	"testing"

	"gitlab.mpi-sws.org/cld/blueprint/plugins/goproc"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/workflow"

	"github.com/stretchr/testify/assert"
)

/*
Tests for correct IR layout from wiring spec helper functions for goproc

Primarily want visibility tests for nodes that are in separate processes but not addressible
*/

func TestServicesWithinSameProcess(t *testing.T) {
	spec := newWiringSpec("TestServicesWithinSameProcess")

	leaf := workflow.Service(spec, "leaf", "TestLeafServiceImpl")
	nonleaf := workflow.Service(spec, "nonleaf", "TestNonLeafService", leaf)

	myproc := goproc.CreateProcess(spec, "myproc", leaf, nonleaf)

	app := assertBuildSuccess(t, spec, myproc)

	assertIR(t, app,
		`TestServicesWithinSameProcess = BlueprintApplication() {
            leaf.handler.visibility
            nonleaf.handler.visibility
            myproc = GolangProcessNode() {
              leaf = TestLeafService()
              nonleaf = TestNonLeafService(leaf)
            }
          }`)
}

func TestSeparateServicesInSeparateProcesses(t *testing.T) {
	spec := newWiringSpec("TestSeparateServicesInSeparateProcesses")

	leaf1 := workflow.Service(spec, "leaf1", "TestLeafServiceImpl")
	leaf2 := workflow.Service(spec, "leaf2", "TestLeafServiceImpl")
	nonleaf := workflow.Service(spec, "nonleaf", "TestNonLeafService", leaf2)

	leaf1proc := goproc.CreateProcess(spec, "leaf1proc", leaf1)
	myproc := goproc.CreateProcess(spec, "myproc", leaf2, nonleaf)

	app := assertBuildSuccess(t, spec, leaf1proc, myproc)

	assertIR(t, app,
		`TestSeparateServicesInSeparateProcesses = BlueprintApplication() {
            leaf1.handler.visibility
            leaf1proc = GolangProcessNode() {
              leaf1 = TestLeafService()
            }
            leaf2.handler.visibility
            nonleaf.handler.visibility
            myproc = GolangProcessNode() {
              leaf2 = TestLeafService()
              nonleaf = TestNonLeafService(leaf2)
            }
          }`)
}

func TestAddChildrenToProcess(t *testing.T) {
	spec := newWiringSpec("TestAddChildrenToProcess")

	myproc := goproc.CreateProcess(spec, "myproc")

	leaf := workflow.Service(spec, "leaf", "TestLeafServiceImpl")
	nonleaf := workflow.Service(spec, "nonleaf", "TestNonLeafService", leaf)

	goproc.AddChildToProcess(spec, myproc, leaf)
	goproc.AddChildToProcess(spec, myproc, nonleaf)

	app := assertBuildSuccess(t, spec, myproc)

	assertIR(t, app,
		`TestAddChildrenToProcess = BlueprintApplication() {
            leaf.handler.visibility
            nonleaf.handler.visibility
            myproc = GolangProcessNode() {
              leaf = TestLeafService()
              nonleaf = TestNonLeafService(leaf)
            }
          }`)

}

func TestReachabilityErrorForSeparateProcesses(t *testing.T) {
	spec := newWiringSpec("TestReachabilityErrorForSeparateProcesses")

	leaf := workflow.Service(spec, "leaf", "TestLeafServiceImpl")
	nonleaf := workflow.Service(spec, "nonleaf", "TestNonLeafService", leaf)

	leafproc := goproc.CreateProcess(spec, "leafproc", leaf)
	nonleafproc := goproc.CreateProcess(spec, "nonleafproc", nonleaf)

	err := assertBuildFailure(t, spec, leafproc, nonleafproc)
	assert.Contains(t, err.Error(), "reachability error")
}

func TestClientWithinSameProcess(t *testing.T) {
	spec := newWiringSpec("TestClientWithinSameProcess")

	leaf := workflow.Service(spec, "leaf", "TestLeafServiceImpl")
	nonleaf := workflow.Service(spec, "nonleaf", "TestNonLeafService", leaf)

	nonleafclient := goproc.CreateClientProcess(spec, "nonleafclient", nonleaf)

	app := assertBuildSuccess(t, spec, nonleafclient)

	assertIR(t, app,
		`TestClientWithinSameProcess = BlueprintApplication() {
            nonleaf.handler.visibility
            leaf.handler.visibility
            nonleafclient = GolangProcessNode() {
              leaf = TestLeafService()
              nonleaf = TestNonLeafService(leaf)
            }
          }`)
}

func TestImplicitServicesWithinSameProcess(t *testing.T) {
	spec := newWiringSpec("TestImplicitServicesWithinSameProcess")

	leaf := workflow.Service(spec, "leaf", "TestLeafServiceImpl")
	nonleaf := workflow.Service(spec, "nonleaf", "TestNonLeafService", leaf)

	nonleafproc := goproc.CreateProcess(spec, "nonleafproc", nonleaf)

	app := assertBuildSuccess(t, spec, nonleafproc)

	assertIR(t, app,
		`TestImplicitServicesWithinSameProcess = BlueprintApplication() {
            nonleaf.handler.visibility
            leaf.handler.visibility
            nonleafproc = GolangProcessNode() {
              leaf = TestLeafService()
              nonleaf = TestNonLeafService(leaf)
            }
          }`)
}
