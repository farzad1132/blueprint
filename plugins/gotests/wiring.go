// Package gotests provides a Blueprint plugin for automatically converting black-box workflow spec unit tests into
// tests that can run against a compiled Blueprint system.
//
// If you have developed a Blueprint application, then you have probably also written tests of your workflow services.
// The gotests plugin lets you leverage those tests and automatically convert them into tests that can run against
// the compiled applciation.
//
// In addition to the documentation here, the [Workflow Tests] page of the user manual has more information.
//
// # Wiring Spec Usage
//
// To use the gotests plugin in a wiring spec, specify which services you want to test:
//
//	gotests.Test(spec, "my_service")
//
// The gotests plugin will search the workflow spec module for any compatible black-box tests, then convert those
// tests into tests that use clients to the compiled Blueprint application.
//
// You will probably also need to ensure that the tests module of your application is on the workflow spec search
// path.  See for example the [SockShop Tests] or [Train Ticket Tests], which have separate tests and workflow modules.
//
//	workflow.Init("../workflow", "../tests")
//
// # Running Tests
//
// After compiling your application, navigate to `gotests` in the output directory.  Within this workspace will be
// a copy of your test module; let's assume it is called `tests`.
//
// Before running the compiled tests, make sure you have started the compiled application.  Then:
//
//	cd gotests/tests
//	go test .
//
// You can use the usual gotest command line arguments to specify individual tests to run.
//
// Depending on the configuration you compiled, the tests might fail due to missing environment variables / arguments.
// For example, the tests will probably need to know the addresses of services to contact.  Remedy this by providing
// those addresses as command line arguments as appropriate.
//
// # Writing Tests: Test Location
//
// The gotests plugin is intended for use with black-box tests of workflow logic.  Most Blueprint applications will
// implement a bunch of services as part of their workflow; ideally developers also write some tests for those services,
// that make API calls to the service and check the receive results.  These are *black-box* tests because they only
// look at the externally-visible inputs and outputs of the service.
//
// By convention we recommend putting black-box tests in a sibling module to the workflow.  See for example the
// [SockShop Tests] or [Train Ticket Tests], which have separate tests and workflow modules.
//
// When compiling a wiring spec, the gotests plugin will need to find the location of the tests.  Like with the
// [workflow] plugin, this is achieved by adding the test module to the workflow spec search path.
//
//	workflow.Init("../workflow", "../tests")
//
// # Writing Tests: Test Compatibility
//
// For a test to be compatible with the gotests plugin, it must make use of the [registry.ServiceRegistry] to
// acquire service instances.  The Train Ticket application's [User Service] demonstrates the use of the ServiceRegistry,
// which we also explain here:
//
// First, in your test file, declare a service registry as a var:
//
//	var userServiceRegistry = registry.NewServiceRegistry[user.UserService]("user_service")
//
// Second, add an init function that instantiates the User Service *locally*.  This local instantiation will
// enable you to run unit tests locally on the workflow spec without having to compile the application
//
//	    func init() {
//	    	userServiceRegistry.Register("local", func(ctx context.Context) (user.UserService, error) {
//	    		db, err := simplenosqldb.NewSimpleNoSQLDB(ctx)
//	    		if err != nil {
//	    			return nil, err
//		    	}
//		    	return user.NewUserServiceImpl(ctx, db)
//		    })
//	    }
//
// Next, you can write a black-box unit test by using the ServiceRegistry's Get method to acquire the userService instance
//
//	func TestUserService(t *testing.T) {
//		ctx := context.Background()
//		service, err := userServiceRegistry.Get(ctx)
//		require.NoError(t, err)
//	 	// ...
//	}
//
// Your tests are now ready to be used with the gotests plugin.  You can also just run the tests locally,
// using `go test .` from your test module.
//
// To summarize, tests are only compatible with the gotests plugin if they:
//   - are contained in a separate module from the workflow spec
//   - the test module is on the workflow spec search path (workflow.Init(...))
//   - the tests declare a [registry.ServiceRegistry] var
//   - the tests using ServiceRegistry.Get to get the service instance to test.
//
// # Edge Cases
//
// Your black-box tests should ideally be idempotent (e.g. if you call Create to instantiate something then
// you should call Delete to remove it too), since we do not automatically spin up / tear down new system
// instances on your behalf.  If you have a large number of services and tests in your application, it is possible
// that one test might leave behind state in the application that causes a different test to fail.  You can
// remedy this by either making the system idempotent, or manually restarting the system and running tests
// individually.
//
// [registry.ServiceRegistry]: https://github.com/blueprint-uservices/blueprint/tree/main/runtime/core/registry
// [SockShop Tests]: https://github.com/blueprint-uservices/blueprint/tree/main/examples/sockshop/tests
// [Train Ticket Tests]: https://github.com/blueprint-uservices/blueprint/tree/main/examples/train_ticket/tests
// [workflow]: https://github.com/blueprint-uservices/blueprint/tree/main/plugins/workflow
// [User Service]: https://github.com/Blueprint-uServices/blueprint/blob/main/examples/sockshop/tests/userservice_test.go
// [Workflow Tests]: https://github.com/blueprint-uservices/blueprint/tree/main/docs/manual/workflow_tests.md
//
// [workflow.Init]: https://github.com/blueprint-uservices/blueprint/tree/main/plugins/workflow
package gotests

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/namespaceutil"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
)

var prop_SERVICESTOTEST = "Services"

// [Test] can be used by wiring specs to convert existing black-box workflow tests into tests that use
// generated service clients and can be run against the compiled Blueprint application.
//
// After compilation, the output will contain a golang workspace called "gotests" that will
// include modified versions of the source tests.
//
// servicesToTest should be the names of golang services instantiated in the wiring spec.
//
// The gotests plugin searches for any workflow packages with tests that make use of [registry.ServiceRegistry].
// Matching modules are copied to an output golang workspace caled "tests".
// Matching packges in the output workspace will have a file blueprint_clients.go that registers
// a service client.
//
// Returns the name "gotests" which must be included when later calling [wiring.WiringSpec.BuildIR]
//
// For more information about tests see [Workflow Tests].
//
// [Workflow Tests]: https://github.com/blueprint-uservices/blueprint/tree/main/docs/manual/workflow_tests.md
// [registry.ServiceRegistry]: https://github.com/blueprint-uservices/blueprint/tree/main/runtime/core/registry
func Test(spec wiring.WiringSpec, servicesToTest ...string) string {

	name := "gotests"

	for _, serviceName := range servicesToTest {
		spec.AddProperty(name, prop_SERVICESTOTEST, serviceName)
	}

	spec.Define(name, &testLibrary{}, func(namespace wiring.Namespace) (ir.IRNode, error) {
		lib := newTestLibrary(name)
		libNamespace, err := namespaceutil.InstantiateNamespace(namespace, &gotests{lib})
		if err != nil {
			return nil, err
		}

		var servicesToTest []string
		if err := namespace.GetProperties(name, prop_SERVICESTOTEST, &servicesToTest); err != nil {
			return nil, err
		}

		for _, serviceName := range servicesToTest {
			var service ir.IRNode
			if err := libNamespace.Get(serviceName, &service); err != nil {
				return nil, err
			}
			lib.ServicesToTest[serviceName] = service
		}
		return lib, err
	})

	return name
}

// A [wiring.NamespaceHandler] used to build the test library
type gotests struct {
	*testLibrary
}

// Implements [wiring.NamespaceHandler]
func (*gotests) Accepts(nodeType any) bool {
	_, isGolangNode := nodeType.(golang.Node)
	return isGolangNode
}

// Implements [wiring.NamespaceHandler]
func (lib *gotests) AddEdge(name string, edge ir.IRNode) error {
	lib.Edges = append(lib.Edges, edge)
	return nil
}

// Implements [wiring.NamespaceHandler]
func (lib *gotests) AddNode(name string, node ir.IRNode) error {
	lib.Nodes = append(lib.Nodes, node)
	return nil
}
