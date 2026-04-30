package drift

import (
	"strings"
	"testing"
	"time"
)

func makeDecayResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []DiffEntry{
				{Key: "replicas", Expected: "3", Actual: "2"},
				{Key: "image", Expected: "v1", Actual: "v2"},
			},
		},
		{
			Service: "worker",
			Diffs: []DiffEntry{
				{Key: "memory", Expected: "512Mi", Actual: "256Mi"},
			},
		},
		{
			Service: "clean-svc",
			Diffs:   []DiffEntry{},
		},
	}
}

func makeDecayHistory(service string, daysAgo float64) HistoryEntry {
	return HistoryEntry{
		Service:   service,
		Timestamp: time.Now().UTC().Add(-time.Duration(daysAgo*24) * time.Hour),
	}
}

func TestApplyDecay_SkipsCleanServices(t *testing.T) {
	results := makeDecayResults()
	entries := ApplyDecay(results, nil, DefaultDecayOptions())
	for _, e := range entries {
		if e.Service == "clean-svc" {
			t.Errorf("clean service should not appear in decay entries")
		}
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestApplyDecay_FullScoreWhenNoHistory(t *testing.T) {
	results := []CompareResult{
		{Service: "api", Diffs: []DiffEntry{{Key: "k", Expected: "a", Actual: "b"}}},
	}
	entries := ApplyDecay(results, nil, DefaultDecayOptions())
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	// With age ~0, decayed score should be very close to raw score (1.0).
	if entries[0].Score < 0.99 {
		t.Errorf("expected score near 1.0 for fresh service, got %.4f", entries[0].Score)
	}
}

func TestApplyDecay_OldHistoryReducesScore(t *testing.T) {
	results := []CompareResult{
		{Service: "api", Diffs: []DiffEntry{{Key: "k", Expected: "a", Actual: "b"}}},
	}
	history := []HistoryEntry{makeDecayHistory("api", 14)} // 2 half-lives
	opts := DefaultDecayOptions()                           // half-life = 7 days
	entries := ApplyDecay(results, history, opts)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	// After 2 half-lives, score should be ~0.25 of raw (1 diff).
	if entries[0].Score > 0.30 {
		t.Errorf("expected decayed score < 0.30, got %.4f", entries[0].Score)
	}
}

func TestApplyDecay_DecayedFlagSet(t *testing.T) {
	results := []CompareResult{
		{Service: "old", Diffs: []DiffEntry{{Key: "x", Expected: "1", Actual: "2"}}},
	}
	history := []HistoryEntry{makeDecayHistory("old", 60)} // very old
	entries := ApplyDecay(results, history, DefaultDecayOptions())
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if !entries[0].Decayed {
		t.Errorf("expected Decayed=true for very old service")
	}
}

func TestApplyDecay_SortedByScoreDescending(t *testing.T) {
	results := makeDecayResults()
	// Give "api" an old timestamp so its score is lower than "worker"'s fresh score.
	history := []HistoryEntry{
		makeDecayHistory("api", 30),
		makeDecayHistory("worker", 0),
	}
	entries := ApplyDecay(results, history, DefaultDecayOptions())
	if len(entries) < 2 {
		t.Fatalf("expected at least 2 entries")
	}
	if entries[0].Score < entries[1].Score {
		t.Errorf("entries not sorted descending by score")
	}
}

func TestFormatDecay_ContainsHeaders(t *testing.T) {
	results := makeDecayResults()
	entries := ApplyDecay(results, nil, DefaultDecayOptions())
	out := FormatDecay(entries)
	for _, hdr := range []string{"SERVICE", "SCORE", "AGE", "DECAYED"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("FormatDecay output missing header %q", hdr)
		}
	}
}

func TestFormatDecay_Empty(t *testing.T) {
	out := FormatDecay(nil)
	if !strings.Contains(out, "no active") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
