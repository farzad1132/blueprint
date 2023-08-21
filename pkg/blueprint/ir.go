package blueprint

import (
	"strings"
)

// This file contains some of the core IR nodes for Blueprint

// The base IRNode type
type IRNode interface {
	Name() string
	String() string
}

type IRMetadata interface {
	ImplementsIRMetadata()
}

// The IR Node that represents the whole application
type ApplicationNode struct {
	IRNode

	name     string
	children map[string]IRNode
}

// For generating output artifacts (e.g. code)
type ArtifactGenerator interface {
	GenerateOutput(string) error
}

func (node *ApplicationNode) Name() string {
	return node.name
}

// Print the IR graph
func (node *ApplicationNode) String() string {
	var b strings.Builder
	b.WriteString(node.name)
	b.WriteString(" = BlueprintApplication() {\n")
	var children []string
	for _, node := range node.children {
		children = append(children, node.String())
	}
	b.WriteString(Indent(strings.Join(children, "\n"), 2))
	b.WriteString("\n}")
	return b.String()
}