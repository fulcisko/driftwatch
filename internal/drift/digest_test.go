package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeDigestResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs:   map[string]any{"replicas": "3 → 5", "image": "v1 → v2"},
		},
		{
			Service: "beta",
			Diffs:   map[string]any{},
		},
		{
			Service: "gamma",
			Diffs:   map[string]any{"timeout": "30 → 60"},
		},
	}
}

func TestComputeDigest_Stable(t *testing.T) {
	r := CompareResult{Service: "svc", Diffs: map[string]any{"a": "1 → 2", "b": "x → y"}}
	h1 := ComputeDigest(r)
	h2 := ComputeDigest(r)
	if h1 != h2 {
		t.Errorf("expected stable hash, got %q vs %q", h1, h2)
	}
}

func TestComputeDigest_DifferentForDifferentDiffs(t *testing.T) {
	r1 := CompareResult{Service: "svc", Diffs: map[string]any{"key": "a → b"}}
	r2 := CompareResult{Service: "svc", Diffs: map[string]any{"key": "a → c"}}
	if ComputeDigest(r1) == ComputeDigest(r2) {
		t.Error("expected different hashes for different diffs")
	}
}

func TestBuildDigests_SkipsClean(t *testing.T) {
	results := makeDigestResults()
	entries := BuildDigests(results)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries (skipping clean), got %d", len(entries))
	}
	services := map[string]bool{}
	for _, e := range entries {
		services[e.Service] = true
	}
	if services["beta"] {
		t.Error("expected beta (no diffs) to be skipped")
	}
}

func TestSaveAndLoadDigests(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "digests.json")

	results := makeDigestResults()
	entries := BuildDigests(results)

	if err := SaveDigests(path, entries); err != nil {
		t.Fatalf("SaveDigests: %v", err)
	}
	loaded, err := LoadDigests(path)
	if err != nil {
		t.Fatalf("LoadDigests: %v", err)
	}
	if len(loaded) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(loaded))
	}
}

func TestLoadDigests_NotFound(t *testing.T) {
	entries, err := LoadDigests("/nonexistent/digests.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestDigestsChanged_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "digests.json")

	original := makeDigestResults()
	prev := BuildDigests(original)
	if err := SaveDigests(path, prev); err != nil {
		t.Fatal(err)
	}

	modified := []CompareResult{
		{Service: "alpha", Diffs: map[string]any{"replicas": "3 → 9"}},
		{Service: "gamma", Diffs: map[string]any{"timeout": "30 → 60"}},
	}
	curr := BuildDigests(modified)

	loaded, _ := LoadDigests(path)
	changed := DigestsChanged(loaded, curr)

	if len(changed) != 1 || changed[0] != "alpha" {
		t.Errorf("expected [alpha] changed, got %v", changed)
	}
}

func TestDigestsChanged_NoDifference(t *testing.T) {
	results := makeDigestResults()
	entries := BuildDigests(results)
	changed := DigestsChanged(entries, entries)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func init() {
	_ = os.Getenv // suppress unused import warning in minimal test env
}
