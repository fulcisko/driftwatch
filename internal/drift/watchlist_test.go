package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeWatchlistResults() []CompareResult {
	return []CompareResult{
		{Service: "alpha", Diffs: []DiffEntry{{Key: "k1"}, {Key: "k2"}}},
		{Service: "beta", Diffs: []DiffEntry{{Key: "k1"}}},
		{Service: "gamma", Diffs: []DiffEntry{}},
	}
}

func TestAddAndLoadWatchlist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watch.json")
	if err := AddToWatchlist(path, "alpha", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wl, err := LoadWatchlist(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(wl.Entries) != 1 || wl.Entries[0].Service != "alpha" {
		t.Errorf("expected alpha, got %+v", wl.Entries)
	}
}

func TestAddWatchlist_Duplicate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watch.json")
	_ = AddToWatchlist(path, "alpha", 1)
	err := AddToWatchlist(path, "alpha", 1)
	if err == nil {
		t.Error("expected error for duplicate service")
	}
}

func TestRemoveFromWatchlist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watch.json")
	_ = AddToWatchlist(path, "alpha", 1)
	_ = AddToWatchlist(path, "beta", 1)
	if err := RemoveFromWatchlist(path, "alpha"); err != nil {
		t.Fatalf("remove error: %v", err)
	}
	wl, _ := LoadWatchlist(path)
	if len(wl.Entries) != 1 || wl.Entries[0].Service != "beta" {
		t.Errorf("expected only beta, got %+v", wl.Entries)
	}
}

func TestRemoveFromWatchlist_NotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watch.json")
	_ = AddToWatchlist(path, "alpha", 1)
	err := RemoveFromWatchlist(path, "ghost")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestLoadWatchlist_NotFound(t *testing.T) {
	_, err := LoadWatchlist("/nonexistent/watch.json")
	if err == nil {
		t.Error("expected error")
	}
}

func TestMatchWatchlist_ThresholdFiltering(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watch.json")
	_ = AddToWatchlist(path, "alpha", 2) // threshold 2 — alpha has 2 diffs, should match
	_ = AddToWatchlist(path, "beta", 2)  // threshold 2 — beta has 1 diff, should not match
	wl, _ := LoadWatchlist(path)
	results := makeWatchlistResults()
	matched := MatchWatchlist(wl, results)
	if len(matched) != 1 || matched[0].Service != "alpha" {
		t.Errorf("expected only alpha, got %+v", matched)
	}
}

func TestMatchWatchlist_NotInList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "watch.json")
	_ = AddToWatchlist(path, "other", 1)
	wl, _ := LoadWatchlist(path)
	results := makeWatchlistResults()
	matched := MatchWatchlist(wl, results)
	if len(matched) != 0 {
		t.Errorf("expected no matches, got %+v", matched)
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
