package drift

import (
	"strings"
	"testing"
	"time"
)

func makeWindowHistory() []HistoryEntry {
	now := time.Now().UTC()
	return []HistoryEntry{
		{Service: "api", RecordedAt: now.Add(-1 * time.Hour)},
		{Service: "worker", RecordedAt: now.Add(-48 * time.Hour)},
		{Service: "gateway", RecordedAt: now.Add(-10 * time.Minute)},
	}
}

func makeWindowResults() []CompareResult {
	return []CompareResult{
		{Service: "api", Diffs: []Diff{{Key: "replicas", Expected: "3", Actual: "2"}}},
		{Service: "worker", Diffs: []Diff{{Key: "timeout", Expected: "30", Actual: "60"}}},
		{Service: "gateway", Diffs: []Diff{}},
	}
}

func TestApplyWindow_IncludesRecentServices(t *testing.T) {
	now := time.Now().UTC()
	opts := WindowOptions{From: now.Add(-2 * time.Hour), To: now}
	results := makeWindowResults()
	history := makeWindowHistory()

	filtered := ApplyWindow(results, history, opts)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 results, got %d", len(filtered))
	}
	services := map[string]bool{}
	for _, r := range filtered {
		services[r.Service] = true
	}
	if !services["api"] || !services["gateway"] {
		t.Errorf("expected api and gateway in window results")
	}
}

func TestApplyWindow_ExcludesOldServices(t *testing.T) {
	now := time.Now().UTC()
	opts := WindowOptions{From: now.Add(-2 * time.Hour), To: now}
	results := makeWindowResults()
	history := makeWindowHistory()

	filtered := ApplyWindow(results, history, opts)
	for _, r := range filtered {
		if r.Service == "worker" {
			t.Errorf("worker should be excluded from window")
		}
	}
}

func TestApplyWindow_EmptyHistory(t *testing.T) {
	now := time.Now().UTC()
	opts := WindowOptions{From: now.Add(-1 * time.Hour), To: now}
	results := makeWindowResults()

	filtered := ApplyWindow(results, nil, opts)
	if len(filtered) != 0 {
		t.Errorf("expected 0 results with empty history, got %d", len(filtered))
	}
}

func TestFormatWindow_ContainsDates(t *testing.T) {
	opts := NewWindowOptions(24 * time.Hour)
	s := FormatWindow(opts)
	if !strings.Contains(s, "from") || !strings.Contains(s, "to") {
		t.Errorf("expected formatted window to contain 'from' and 'to', got: %s", s)
	}
}

func TestNewWindowOptions_Duration(t *testing.T) {
	before := time.Now().UTC()
	opts := NewWindowOptions(6 * time.Hour)
	after := time.Now().UTC()

	if opts.To.Before(before) || opts.To.After(after) {
		t.Errorf("window To should be approximately now")
	}
	expectedFrom := opts.To.Add(-6 * time.Hour)
	if opts.From.After(expectedFrom.Add(time.Second)) || opts.From.Before(expectedFrom.Add(-time.Second)) {
		t.Errorf("window From should be 6h before To")
	}
}
