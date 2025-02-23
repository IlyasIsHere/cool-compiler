package semant

import (
	"cool-compiler/ast"
	"cool-compiler/lexer"
	"reflect"
	"testing"
)

func createTestClass(name string, parent string) *ast.Class {
	class := &ast.Class{
		Name: &ast.TypeIdentifier{
			Token: lexer.Token{Literal: name},
			Value: name,
		},
	}
	if parent != "" {
		class.Parent = &ast.TypeIdentifier{
			Token: lexer.Token{Literal: parent},
			Value: parent,
		}
	}
	return class
}

func TestInheritanceGraph_AddNodeAndEdge(t *testing.T) {
	ig := NewInheritanceGraph()
	class := createTestClass("A", "Object")

	ig.AddNode("A", class)
	if _, exists := ig.nodes["A"]; !exists {
		t.Error("Expected node A to be added to the graph")
	}

	ig.AddEdge("A", "Object")
	if parent, exists := ig.edges["A"]; !exists || parent != "Object" {
		t.Error("Expected edge A->Object to be added to the graph")
	}
}

func TestInheritanceGraph_DetectCycle(t *testing.T) {
	tests := []struct {
		name          string
		edges         map[string]string
		start         string
		expectCycle   bool
		expectedCycle []string
	}{
		{
			name: "No cycle",
			edges: map[string]string{
				"A": "Object",
				"B": "A",
				"C": "B",
			},
			start:       "C",
			expectCycle: false,
		},
		{
			name: "Simple cycle",
			edges: map[string]string{
				"A": "B",
				"B": "C",
				"C": "A",
			},
			start:         "A",
			expectCycle:   true,
			expectedCycle: []string{"A", "B", "C", "A"},
		},
		{
			name: "Self cycle",
			edges: map[string]string{
				"A": "A",
			},
			start:         "A",
			expectCycle:   true,
			expectedCycle: []string{"A", "A"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ig := NewInheritanceGraph()
			for class, parent := range tt.edges {
				ig.AddEdge(class, parent)
			}

			hasCycle, cycle := ig.detectCycle(tt.start)
			if hasCycle != tt.expectCycle {
				t.Errorf("Expected cycle detection to be %v, got %v", tt.expectCycle, hasCycle)
			}

			if tt.expectCycle && !reflect.DeepEqual(cycle, tt.expectedCycle) {
				t.Errorf("Expected cycle %v, got %v", tt.expectedCycle, cycle)
			}
		})
	}
}

func TestInheritanceGraph_IsConformant(t *testing.T) {
	tests := []struct {
		name     string
		edges    map[string]string
		type1    string
		type2    string
		expected bool
	}{
		{
			name: "Direct inheritance",
			edges: map[string]string{
				"A": "Object",
				"B": "A",
			},
			type1:    "B",
			type2:    "A",
			expected: true,
		},
		{
			name: "Indirect inheritance",
			edges: map[string]string{
				"A": "Object",
				"B": "A",
				"C": "B",
			},
			type1:    "C",
			type2:    "A",
			expected: true,
		},
		{
			name: "No conformance",
			edges: map[string]string{
				"A": "Object",
				"B": "Object",
			},
			type1:    "A",
			type2:    "B",
			expected: false,
		},
		{
			name: "Same type",
			edges: map[string]string{
				"A": "Object",
			},
			type1:    "A",
			type2:    "A",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ig := NewInheritanceGraph()
			for class, parent := range tt.edges {
				ig.AddEdge(class, parent)
			}

			result := ig.IsConformant(tt.type1, tt.type2)
			if result != tt.expected {
				t.Errorf("Expected conformance to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInheritanceGraph_FindLCA(t *testing.T) {
	tests := []struct {
		name     string
		edges    map[string]string
		type1    string
		type2    string
		expected string
	}{
		{
			name: "Common direct ancestor",
			edges: map[string]string{
				"B": "A",
				"C": "A",
				"A": "Object",
			},
			type1:    "B",
			type2:    "C",
			expected: "A",
		},
		{
			name: "One is ancestor of other",
			edges: map[string]string{
				"B": "A",
				"C": "B",
				"A": "Object",
			},
			type1:    "C",
			type2:    "A",
			expected: "A",
		},
		{
			name: "Same type",
			edges: map[string]string{
				"A": "Object",
			},
			type1:    "A",
			type2:    "A",
			expected: "A",
		},
		{
			name: "Only Object in common",
			edges: map[string]string{
				"A": "Object",
				"B": "Object",
			},
			type1:    "A",
			type2:    "B",
			expected: "Object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ig := NewInheritanceGraph()
			for class, parent := range tt.edges {
				ig.AddEdge(class, parent)
			}

			result := ig.FindLCA(tt.type1, tt.type2)
			if result != tt.expected {
				t.Errorf("Expected LCA to be %v, got %v", tt.expected, result)
			}
		})
	}
}
