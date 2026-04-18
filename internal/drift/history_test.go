package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeHistoryResults() []CompareResult {
	return []CompareResult{
		{
			Service: "svc-a",
			Diffs: []Diff{
				{Key: "replicas", Expected: "3", Actual: "2"},
			},
		},
	}
}

func TestAppendAndLoadHistory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	if err := AppendHistory(path, makeHistoryResults()); err != nil {
		t.Fatalf("first append: %v", err)
	}
	if err := AppendHistory(path, makeHistoryResults()); err != nil {
		t.Fatalf("second append: %v", err)
	}

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Results[0].Service != "svc-a" {
		t.Errorf("unexpected service: %s", entries[0].Results[0].Service)
	}
}

func TestLoadHistory_NotFound(t *testing.T) {
	_, err := LoadHistory("/nonexistent/history.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !os.IsNotExist(err) {
		t.Logf("got non-IsNotExist error (acceptable): %v", err)
	}
}

func TestLatestHistory_ReturnsLast(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	first := []CompareResult{{Service: "old"}}
	second := []CompareResult{{Service: "new"}}

	_ = AppendHistory(path, first)
	_ = AppendHistory(path, second)

	entry, err := LatestHistory(path)
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Results[0].Service != "new" {
		t.Errorf("expected 'new', got %s", entry.Results[0].Service)
	}
}

func TestLatestHistory_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	// Write empty array
	os.WriteFile(path, []byte("[]"), 0644)

	entry, err := LatestHistory(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry for empty history")
	}
}
