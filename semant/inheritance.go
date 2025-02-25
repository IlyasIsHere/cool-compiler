package semant

import (
	"cool-compiler/ast"
)

type InheritanceGraph struct {
	// Maps class name to its parent class name
	edges map[string]string
	// Maps class name to its Class AST node
	nodes map[string]*ast.Class
}

func NewInheritanceGraph() *InheritanceGraph {
	return &InheritanceGraph{
		edges: make(map[string]string),
		nodes: make(map[string]*ast.Class),
	}
}

func (ig *InheritanceGraph) AddEdge(class, parent string) {
	ig.edges[class] = parent
}

func (ig *InheritanceGraph) AddNode(className string, class *ast.Class) {
	ig.nodes[className] = class
}

func (ig *InheritanceGraph) GetNode(className string) *ast.Class {
	if node, ok := ig.nodes[className]; ok {
		return node
	}
	return nil
}

// detectCycle returns true and the cycle path if there's a cycle
func (ig *InheritanceGraph) detectCycle(start string) (bool, []string) {
	visited := make(map[string]bool)
	path := make(map[string]bool)
	cycle := []string{}

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		path[node] = true
		cycle = append(cycle, node)

		if parent, exists := ig.edges[node]; exists {
			if !visited[parent] {
				if dfs(parent) {
					return true
				}
			} else if path[parent] {
				cycle = append(cycle, parent)
				return true
			}
		}

		path[node] = false
		cycle = cycle[:len(cycle)-1]
		return false
	}

	hasCycle := dfs(start)
	return hasCycle, cycle
}

// IsConformant checks if type1 conforms to type2 (type1 ≤ type2)
func (ig *InheritanceGraph) IsConformant(type1, type2 string) bool {
	if type1 == type2 {
		return true
	}

	current := type1
	for {
		parent, exists := ig.edges[current]
		if !exists {
			return false
		}
		if parent == type2 {
			return true
		}
		current = parent
	}
}

// FindLCA finds the least common ancestor of two types
func (ig *InheritanceGraph) FindLCA(type1, type2 string) string {
	if type1 == type2 {
		return type1
	}

	// Get path from type1 to Object
	path1 := make(map[string]bool)
	current := type1
	path1[current] = true
	for {
		parent, exists := ig.edges[current]
		if !exists {
			break
		}
		path1[parent] = true
		current = parent
	}

	// Walk up type2's path until we find a common ancestor
	current = type2
	if path1[current] {
		return current
	}

	for {
		parent, exists := ig.edges[current]
		if !exists {
			return "Object" // Default to Object if no common ancestor found
		}
		if path1[parent] {
			return parent
		}
		current = parent
	}
}

// GetParent returns the parent class name of a given class, if it exists.
func (ig *InheritanceGraph) GetParent(className string) (string, bool) {
	parent, exists := ig.edges[className]
	return parent, exists
}
