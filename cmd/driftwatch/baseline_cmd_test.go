package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

func writeBaselineFile(t *testing.T, dir, name string, results []drift.CompareResult) string {
	t.Helper()
	b := drift.Baseline{
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
	data, _ := json.MarshalIndent(b, "", "  ")
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunSaveBaseline_Success(t *testing.T) {
	dir := t.TempDir()
	results := []drift.CompareResult{
		{Service: "svc-a", Diffs: nil},
	}
	src := writeBaselineFile(t, dir, "results.json", results)
	dst := filepath.Join(dir, "baseline.json")

	if err := runSaveBaseline(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Error("baseline file not created")
	}
}

func TestRunSaveBaseline_MissingSource(t *testing.T) {
	err := runSaveBaseline("/no/such/file.json", "/tmp/out.json")
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestRunDiffBaseline_NoChange(t *testing.T) {
	dir := t.TempDir()
	results := []drift.CompareResult{{Service: "svc-a", Diffs: nil}}
	base := writeBaselineFile(t, dir, "base.json", results)
	curr := writeBaselineFile(t, dir, "curr.json", results)

	if err := runDiffBaseline(base, curr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDiffBaseline_WithChanges(t *testing.T) {
	dir := t.TempDir()
	base := writeBaselineFile(t, dir, "base.json", []drift.CompareResult{
		{Service: "svc-a", Diffs: nil},
	})
	curr := writeBaselineFile(t, dir, "curr.json", []drift.CompareResult{
		{Service: "svc-a", Diffs: []drift.Diff{
			{Key: "replicas", Expected: "3", Actual: "1", Kind: drift.KindChanged},
		}},
	})
	if err := runDiffBaseline(base, curr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
