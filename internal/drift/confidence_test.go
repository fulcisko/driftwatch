package drift

import (
	"strings"
	"testing"
	"time"
)

func makeConfidenceResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api-gateway",
			Diffs: []DiffEntry{
				{Key: "replicas", Expected: "3", Actual: "1"},
				{Key: "timeout", Expected: "30s", Actual: "10s"},
			},
		},
		{
			Service: "auth-service",
			Diffs: []DiffEntry{
				{Key: "env", Expected: "prod", Actual: "staging"},
			},
		},
		{
			Service: "clean-service",
			Diffs:   []DiffEntry{},
		},
	}
}

func makeConfidenceHistory(service string, count int) []HistoryEntry {
	var entries []HistoryEntry
	for i := 0; i < count; i++ {
		entries = append(entries, HistoryEntry{
			Timestamp: time.Now().Add(-time.Duration(i) * time.Hour),
			Results: []CompareResult{
				{Service: service, Diffs: []DiffEntry{{Key: "k", Expected: "a", Actual: "b"}}},
			},
		})
	}
	return entries
}

func TestScoreConfidence_SkipsClean(t *testing.T) {
	results := makeConfidenceResults()
	out := ScoreConfidence(results, nil)
	for _, r := range out {
		if r.Service == "clean-service" {
			t.Errorf("clean-service should be excluded from confidence results")
		}
	}
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestScoreConfidence_HighWithHistory(t *testing.T) {
	results := makeConfidenceResults()
	history := makeConfidenceHistory("api-gateway", 5)
	out := ScoreConfidence(results, history)
	var found *ConfidenceResult
	for i := range out {
		if out[i].Service == "api-gateway" {
			found = &out[i]
		}
	}
	if found == nil {
		t.Fatal("api-gateway not found in results")
	}
	if found.Level != ConfidenceHigh {
		t.Errorf("expected high confidence, got %s (score=%.2f)", found.Level, found.Score)
	}
}

func TestScoreConfidence_LowWithNoHistory(t *testing.T) {
	results := []CompareResult{
		{Service: "new-svc", Diffs: []DiffEntry{{Key: "x", Expected: "1", Actual: "2"}}},
	}
	out := ScoreConfidence(results, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Level != ConfidenceLow {
		t.Errorf("expected low confidence for new service, got %s", out[0].Level)
	}
}

func TestScoreConfidence_SortedByScore(t *testing.T) {
	results := makeConfidenceResults()
	history := makeConfidenceHistory("auth-service", 6)
	out := ScoreConfidence(results, history)
	for i := 1; i < len(out); i++ {
		if out[i].Score > out[i-1].Score {
			t.Errorf("results not sorted descending by score at index %d", i)
		}
	}
}

func TestFormatConfidence_ContainsService(t *testing.T) {
	results := makeConfidenceResults()
	out := ScoreConfidence(results, nil)
	formatted := FormatConfidence(out)
	if !strings.Contains(formatted, "api-gateway") {
		t.Errorf("expected api-gateway in formatted output")
	}
	if !strings.Contains(formatted, "score=") {
		t.Errorf("expected score= in formatted output")
	}
}

func TestFormatConfidence_Empty(t *testing.T) {
	formatted := FormatConfidence(nil)
	if !strings.Contains(formatted, "no drifted services") {
		t.Errorf("expected empty message, got: %s", formatted)
	}
}
