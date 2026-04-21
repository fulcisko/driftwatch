package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func depPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "deps.json")
}

func TestAddAndLoadDependency(t *testing.T) {
	p := depPath(t)
	if err := AddDependency(p, "api", "db", "uses"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	graph, err := LoadDependencies(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(graph.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(graph.Edges))
	}
	e := graph.Edges[0]
	if e.From != "api" || e.To != "db" || e.Label != "uses" {
		t.Errorf("unexpected edge: %+v", e)
	}
}

func TestAddDependency_Duplicate(t *testing.T) {
	p := depPath(t)
	_ = AddDependency(p, "api", "db", "uses")
	err := AddDependency(p, "api", "db", "uses")
	if err == nil {
		t.Fatal("expected duplicate error, got nil")
	}
}

func TestLoadDependencies_NotFound(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	graph, err := LoadDependencies(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graph.Edges) != 0 {
		t.Errorf("expected empty graph")
	}
}

func TestDependentsOf(t *testing.T) {
	p := depPath(t)
	_ = AddDependency(p, "api", "db", "")
	_ = AddDependency(p, "worker", "db", "")
	_ = AddDependency(p, "api", "cache", "")
	graph, _ := LoadDependencies(p)
	deps := DependentsOf(graph, "db")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependents of db, got %d", len(deps))
	}
}

func TestDependenciesOf(t *testing.T) {
	p := depPath(t)
	_ = AddDependency(p, "api", "db", "")
	_ = AddDependency(p, "api", "cache", "")
	graph, _ := LoadDependencies(p)
	deps := DependenciesOf(graph, "api")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies for api, got %d", len(deps))
	}
}

func TestDependenciesOf_Empty(t *testing.T) {
	graph := DependencyGraph{}
	result := DependenciesOf(graph, "orphan")
	if len(result) != 0 {
		t.Errorf("expected no dependencies")
	}
}

func TestAddDependency_MultipleEdges(t *testing.T) {
	p := depPath(t)
	services := []struct{ from, to string }{
		{"a", "b"}, {"b", "c"}, {"c", "d"},
	}
	for _, s := range services {
		if err := AddDependency(p, s.from, s.to, ""); err != nil {
			t.Fatalf("add error: %v", err)
		}
	}
	graph, _ := LoadDependencies(p)
	if len(graph.Edges) != 3 {
		t.Errorf("expected 3 edges, got %d", len(graph.Edges))
	}
	_ = os.Remove(p)
}
