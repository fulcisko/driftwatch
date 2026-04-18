package drift

import (
	"path/filepath"
	"testing"
)

func makeSnapshotResults() []CompareResult {
	return []CompareResult{
		{Service: "alpha", Diffs: []Diff{{Key: "port", Expected: "8080", Actual: "9090"}}},
		{Service: "beta", Diffs: []Diff{}},
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	results := makeSnapshotResults()

	if err := SaveSnapshot(path, "v1", results); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if snap.Label != "v1" {
		t.Errorf("expected label v1, got %s", snap.Label)
	}
	if len(snap.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(snap.Results))
	}
}

func TestLoadSnapshot_NotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	snap := Snapshot{
		Label: "base",
		Results: []CompareResult{
			{Service: "alpha", Diffs: []Diff{{Key: "port", Expected: "8080", Actual: "9090"}}},
			{Service: "beta", Diffs: []Diff{}},
		},
	}
	// alpha now clean, beta now drifted
	current := []CompareResult{
		{Service: "alpha", Diffs: []Diff{}},
		{Service: "beta", Diffs: []Diff{{Key: "timeout", Expected: "30", Actual: "60"}}},
	}
	changed := DiffSnapshot(snap, current)
	if len(changed) != 2 {
		t.Errorf("expected 2 changed services, got %d: %v", len(changed), changed)
	}
}

func TestDiffSnapshot_NoDifference(t *testing.T) {
	results := makeSnapshotResults()
	snap := Snapshot{Label: "base", Results: results}
	changed := DiffSnapshot(snap, results)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func TestDiffSnapshot_NewService(t *testing.T) {
	snap := Snapshot{
		Label:   "base",
		Results: []CompareResult{{Service: "alpha", Diffs: []Diff{}}},
	}
	current := []CompareResult{
		{Service: "alpha", Diffs: []Diff{}},
		{Service: "gamma", Diffs: []Diff{{Key: "env", Expected: "prod", Actual: "staging"}}},
	}
	changed := DiffSnapshot(snap, current)
	if len(changed) != 1 || changed[0] != "gamma" {
		t.Errorf("expected [gamma], got %v", changed)
	}
}
