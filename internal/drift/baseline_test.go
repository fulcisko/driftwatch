package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeBaselineResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs: []Diff{
				{Key: "replicas", Expected: "3", Actual: "2", Kind: KindChanged},
			},
		},
		{
			Service: "beta",
			Diffs:   nil,
		},
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	results := makeBaselineResults()

	if err := SaveBaseline(path, results); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if len(b.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(b.Results))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if b.Results[0].Service != "alpha" {
		t.Errorf("expected alpha, got %s", b.Results[0].Service)
	}
}

func TestLoadBaseline_NotFound(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDiffBaseline_DetectsNewDrift(t *testing.T) {
	baseline := &Baseline{
		CreatedAt: time.Now(),
		Results:   makeBaselineResults(),
	}
	current := []CompareResult{
		{Service: "alpha", Diffs: nil},          // drift resolved
		{Service: "beta", Diffs: []Diff{
			{Key: "image", Expected: "v1", Actual: "v2", Kind: KindChanged},
		}}, // new drift
	}
	changed := DiffBaseline(baseline, current)
	if len(changed) != 2 {
		t.Errorf("expected 2 changed services, got %d", len(changed))
	}
}

func TestDiffBaseline_NoDifference(t *testing.T) {
	results := makeBaselineResults()
	baseline := &Baseline{CreatedAt: time.Now(), Results: results}
	changed := DiffBaseline(baseline, results)
	if len(changed) != 0 {
		t.Errorf("expected 0 changed, got %d", len(changed))
	}
}

func TestDiffBaseline_NewService(t *testing.T) {
	baseline := &Baseline{CreatedAt: time.Now(), Results: makeBaselineResults()}
	current := append(makeBaselineResults(), CompareResult{Service: "gamma", Diffs: []Diff{
		{Key: "port", Expected: "8080", Actual: "9090", Kind: KindChanged},
	}})
	changed := DiffBaseline(baseline, current)
	if len(changed) != 1 || changed[0].Service != "gamma" {
		t.Errorf("expected gamma as new service, got %+v", changed)
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
