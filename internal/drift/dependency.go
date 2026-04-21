package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// DependencyEdge represents a directional dependency between two services.
type DependencyEdge struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Label     string    `json:"label,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// DependencyGraph holds all edges in the dependency map.
type DependencyGraph struct {
	Edges []DependencyEdge `json:"edges"`
}

// AddDependency appends a dependency edge to the graph file, deduplicating by (from, to).
func AddDependency(path, from, to, label string) error {
	graph, err := LoadDependencies(path)
	if err != nil {
		return err
	}
	for _, e := range graph.Edges {
		if e.From == from && e.To == to {
			return fmt.Errorf("dependency %q -> %q already exists", from, to)
		}
	}
	graph.Edges = append(graph.Edges, DependencyEdge{
		From:      from,
		To:        to,
		Label:     label,
		CreatedAt: time.Now().UTC(),
	})
	return saveDependencies(path, graph)
}

// LoadDependencies reads the dependency graph from disk, returning an empty graph if not found.
func LoadDependencies(path string) (DependencyGraph, error) {
	var graph DependencyGraph
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return graph, nil
	}
	if err != nil {
		return graph, err
	}
	return graph, json.Unmarshal(data, &graph)
}

// DependentsOf returns all service names that depend on the given service.
func DependentsOf(graph DependencyGraph, service string) []string {
	var result []string
	for _, e := range graph.Edges {
		if e.To == service {
			result = append(result, e.From)
		}
	}
	return result
}

// DependenciesOf returns all services that the given service depends on.
func DependenciesOf(graph DependencyGraph, service string) []string {
	var result []string
	for _, e := range graph.Edges {
		if e.From == service {
			result = append(result, e.To)
		}
	}
	return result
}

func saveDependencies(path string, graph DependencyGraph) error {
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
