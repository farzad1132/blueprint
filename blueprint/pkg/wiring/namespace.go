package wiring

import (
	"fmt"
	"reflect"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/blueprint/logging"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint/pkg/ir"
	"golang.org/x/exp/slog"
)

/*
A Namespace is used during the IR-building process to accumulate built nodes.

Blueprint has several basic out-of-the-box namespaces that are used when building applications.  A plugin can implement
its own custom namespace.  Implementing a custom Namespace is useful to achieve any of the following:
  - Namespaces are the mechanism for limiting the visibility and addressibility of nodes
  - Namespaces are the mechanism for templating nodes (e.g. to implement replication of nodes)

For example, to build a GoProcess that contains Golang object instances, there will be a Namespace that accumulates
Golang object instance nodes during the building process, and then creates a GoProcess namespace node.

Most namespace implementations should extend the BasicNamespace struct
*/
type Namespace interface {
	Name() string                                         // The name of this namespace
	Get(name string, dst any) error                       // Get an ir.IRNode from this namespace or a parent namespace, possibly building it.  Places the ir.IRNode in the pointer dst.  dst can be an ir.IRNode or any implementation of an ir.IRNode
	Instantiate(name string, dst any) error               // The same as Get, but without creating a dependency (an edge) into the current namespace.  Places the ir.IRNode in the pointer dst.  dst can be an ir.IRNode or any implementation of an ir.IRNode
	GetProperty(name string, key string, dst any) error   // Get a property from this namespace; dst should be a pointer to value
	GetProperties(name string, key string, dst any) error // Get a slice property from this namespace; dst should be a pointer to a slice
	Put(name string, node ir.IRNode) error                // Put a node into this namespace
	Defer(f func() error)                                 // Enqueue a function to be executed once finished building the current nodes

	Info(message string, args ...any)        // Logging
	Warn(message string, args ...any)        // Logging
	Error(message string, args ...any) error // Logging
}

/*
A SimpleNamespace implements all of the Namespace methods and only requires users to implement a SimpleNamespaceHandler interface.
Most plugins will want to use SimpleNamespace rather than directly implementing Namespace.

See the documentation of SimpleNamespaceHandler for methods to override.
*/
type SimpleNamespace struct {
	Namespace

	NamespaceName   string                 // A name for this namespace
	NamespaceType   string                 // The type of this namespace
	ParentNamespace Namespace              // The parent namespace that created this namespace; can be nil
	Wiring          WiringSpec             // The wiring spec
	Handler         SimpleNamespaceHandler // User-provided handler
	Seen            map[string]ir.IRNode   // Cache of built nodes
	Added           map[string]any         // Nodes that have been passed to the handler
	Deferred        []func() error         // Deferred functions to execute

	stack []*WiringDef // Used when building; the stack of wiring defs currently being built
}

/*
Has four methods with default implementations that callers can override with custom logic:
  - LookupDef(name) - look up a WiringDef; default implementation directly consults the WiringSpec.
    callers can override this if they want to restrict, modify, or wrap definitions
    that get instantiated within this namespace.
  - Accepts(nodeType) - should return true if the specified node type should be built within this namespace,
    or false if we should ask the parent to build it instead.  Most namespace implementations will only
    accept certain node types, and will thus want to override this method.  For example, a golang process
    will only accept golang nodes
  - AddNode(name, ir.IRNode) - this is called when a node is created within this namespace.  The SimpleNamespace
    internally saves the node for future lookups; callers might want to save the node e.g. as a child within
    a node that is being created.
  - AddEdge(name, ir.IRNode) - this is called when a node was created by a parent namespace but referenced within
    this namespace.  The SimpleNamespace internally saves the node for future lookups; callers might want to save the
    node e.g. as an argument to the node that is being created
*/
type SimpleNamespaceHandler interface {
	Init(*SimpleNamespace)
	LookupDef(string) (*WiringDef, error)
	Accepts(any) bool
	AddEdge(string, ir.IRNode) error
	AddNode(string, ir.IRNode) error
}

type DefaultNamespaceHandler struct {
	SimpleNamespaceHandler
	Namespace *SimpleNamespace

	Nodes []ir.IRNode
	Edges []ir.IRNode
}

func (handler *DefaultNamespaceHandler) Init(namespace *SimpleNamespace) {
	handler.Namespace = namespace
}

/*
Look up a WiringDef; default implementation directly consults the WiringSpec.

	callers can override this if they want to restrict, modify, or wrap definitions
	that get instantiated within this namespace.
*/
func (handler *DefaultNamespaceHandler) LookupDef(name string) (*WiringDef, error) {
	def := handler.Namespace.Wiring.GetDef(name)
	if def == nil {
		return nil, blueprint.Errorf("%s does not exist in the wiring spec of namespace %s", name, handler.Namespace.Name())
	}
	return def, nil
}

/*
should return true if the specified node type should be built within this namespace, or false if we should ask the parent to build it instead.  Most namespace implementations will only

	accept certain node types, and will thus want to override this method.  For example, a golang process
	will only accept golang nodes
*/
func (handler *DefaultNamespaceHandler) Accepts(nodeType any) bool {
	return true
}

// This is called after getting a node from the parent namespace.  By default it just saves the node
// as an edge.  Namespace implementations can override this method to do other things.
func (handler *DefaultNamespaceHandler) AddEdge(name string, node ir.IRNode) error {
	handler.Edges = append(handler.Edges, node)
	return nil
}

// This is called after building a node in the current namespace.  By default it just saves the node
// on the namespace.  Namespace implementations can override this method to do other things.
func (handler *DefaultNamespaceHandler) AddNode(name string, node ir.IRNode) error {
	handler.Nodes = append(handler.Nodes, node)
	return nil
}

func (namespace *SimpleNamespace) Init(name, namespacetype string, parent Namespace, wiring WiringSpec, handler SimpleNamespaceHandler) {
	namespace.NamespaceName = name
	namespace.NamespaceType = namespacetype
	namespace.ParentNamespace = parent
	namespace.Wiring = wiring
	namespace.Handler = handler
	namespace.Seen = make(map[string]ir.IRNode)
	namespace.Added = make(map[string]any)
}
func (namespace *SimpleNamespace) Name() string {
	return namespace.NamespaceName
}

func (namespace *SimpleNamespace) Instantiate(name string, dst any) error {
	return namespace.get(name, false, dst)
}

func (namespace *SimpleNamespace) Get(name string, dst any) error {
	return namespace.get(name, true, dst)
}

func (namespace *SimpleNamespace) get(name string, addEdge bool, dst any) error {
	// If it already exists, return it
	if node, ok := namespace.Seen[name]; ok {
		return copyResult(node, dst)
	}

	// Look up the definition
	def, err := namespace.Handler.LookupDef(name)
	if err != nil {
		return err
	}

	// Track the defs being built
	namespace.stack = append(namespace.stack, def)
	defer func() {
		namespace.stack = namespace.stack[:len(namespace.stack)-1]
	}()

	// If it's an alias, get the aliased node
	if def.Name != name {
		namespace.Info("Resolved %s to %s", name, def.Name)
		var node ir.IRNode
		err := namespace.get(def.Name, addEdge, &node)
		namespace.Seen[name] = node
		if err != nil {
			return err
		}
		return copyResult(node, dst)
	}

	// See if the node should be created here or in the parent
	if !namespace.Handler.Accepts(def.NodeType) {
		if namespace.ParentNamespace == nil {
			return namespace.Error("Namespace does not accept node %s of type %s but there is no parent namespace to get them from", name, reflect.TypeOf(def.NodeType).String())
		}
		namespace.Info("Getting %s of type %s from parent namespace %s", name, reflect.TypeOf(def.NodeType).String(), namespace.ParentNamespace.Name())
		var node ir.IRNode
		if addEdge {
			err = namespace.ParentNamespace.Get(name, &node)
		} else {
			err = namespace.ParentNamespace.Instantiate(name, &node)
		}
		if err != nil {
			return err
		}
		if _, already_added := namespace.Added[node.Name()]; !already_added {
			if _, is_metadata := node.(ir.IRMetadata); !is_metadata && addEdge {
				// Don't bother adding edges for metadata
				namespace.Handler.AddEdge(name, node)
			}
			namespace.Added[node.Name()] = true
		}
		namespace.Seen[name] = node
		return copyResult(node, dst)
	}

	if def.Name == name {
		namespace.Info("Building %s of type %s", name, reflect.TypeOf(def.NodeType).String())
	} else {
		namespace.Info("Building %s (alias %s) of type %s", def.Name, name, reflect.TypeOf(def.NodeType).String())
	}

	// Build the node
	node, err := def.Build(namespace)
	if err != nil {
		namespace.Error("Unable to build %v: %s", name, err.Error())
		return err
	}

	if _, already_added := namespace.Added[node.Name()]; !already_added {
		namespace.Handler.AddNode(name, node)
		namespace.Added[node.Name()] = true
	}
	namespace.Info("Finished building %s of type %s", name, reflect.TypeOf(node).String())
	namespace.Seen[name] = node
	return copyResult(node, dst)
}

func (namespace *SimpleNamespace) Put(name string, node ir.IRNode) error {
	namespace.Seen[name] = node

	if namespace.Handler.Accepts(node) {
		namespace.Handler.AddNode(name, node)
		namespace.Info("%s of type %s added to namespace", name, reflect.TypeOf(node).Elem().Name())
		return nil
	}

	if namespace.ParentNamespace != nil {
		return namespace.Error("%s of type %s does not belong in this namespace, but cannot push to parent namespace because no parent namespace exists", name, reflect.TypeOf(node).Elem().Name())
	}

	namespace.Info("%s of type %s does not belong in this namespace; pushing to parent namespace %s", name, reflect.TypeOf(node).Elem().Name(), namespace.ParentNamespace)
	err := namespace.ParentNamespace.Put(name, node)
	if err != nil {
		return err
	}
	namespace.Handler.AddEdge(name, node)
	return err
}

func (namespace *SimpleNamespace) Defer(f func() error) {
	if namespace.ParentNamespace == nil {
		namespace.Deferred = append(namespace.Deferred, f)
	} else {
		namespace.ParentNamespace.Defer(f)
	}
}

func (namespace *SimpleNamespace) GetProperty(name string, key string, dst any) error {
	def, err := namespace.Handler.LookupDef(name)
	if err != nil {
		return err
	}
	return def.GetProperty(key, dst)
}

func (namespace *SimpleNamespace) GetProperties(name string, key string, dst any) error {
	def, err := namespace.Handler.LookupDef(name)
	if err != nil {
		return err
	}
	return def.GetProperties(key, dst)
}

// Augments debug messages with information about the namespace
func (namespace *SimpleNamespace) Info(message string, args ...any) {
	if len(namespace.stack) > 0 {
		src := namespace.stack[len(namespace.stack)-1]
		callstack := src.Properties["callsite"][0].(*logging.Callstack)
		slog.Info(fmt.Sprintf(fmt.Sprintf("%s %s: %s (%s)", namespace.NamespaceType, namespace.Name(), message, callstack.Stack[0].String()), args...))
	} else {
		slog.Info(fmt.Sprintf(fmt.Sprintf("%s %s: %s", namespace.NamespaceType, namespace.Name(), message), args...))
	}
}

// Augments debug messages with information about the namespace
func (namespace *SimpleNamespace) Debug(message string, args ...any) {
	if len(namespace.stack) > 0 {
		src := namespace.stack[len(namespace.stack)-1]
		callstack := src.Properties["callsite"][0].(*logging.Callstack)
		slog.Info(callstack.String())
		slog.Debug(fmt.Sprintf(fmt.Sprintf("%s %s: %s (%s)", namespace.NamespaceType, namespace.Name(), message, callstack.Stack[0].String()), args...))
	} else {
		slog.Debug(fmt.Sprintf(fmt.Sprintf("%s %s: %s", namespace.NamespaceType, namespace.Name(), message), args...))
	}
}

// Augments debug messages with information about the namespace
func (namespace *SimpleNamespace) Error(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	if len(namespace.stack) > 0 {
		src := namespace.stack[len(namespace.stack)-1]
		callstack := src.Properties["callsite"][0].(*logging.Callstack)
		slog.Error(fmt.Sprintf("%s %s: %s (%s)", namespace.NamespaceType, namespace.Name(), formattedMessage, callstack.Stack[0].String()))
	} else {
		slog.Error(fmt.Sprintf("%s %s: %s", namespace.NamespaceType, namespace.Name(), formattedMessage))
	}
	return fmt.Errorf(formattedMessage)
}
