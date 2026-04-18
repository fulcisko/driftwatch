package drift

import (
	"strings"
	"testing"
)

func makeSummaryResults() map[string][]CompareResult {
	return map[string][]CompareResult{
		"svc-a": {
			{Key: "replicas", Expected: "3", Actual: "2", Kind: KindChanged},
		},
		"svc-b": {},
		"svc-c": {
			{Key: "image", Expected: "v1", Actual: "", Kind: KindMissing},
			{Key: "port", Expected: "8080", Actual: "9090", Kind: KindChanged},
		},
	}
}

func TestSummarize_Counts(t *testing.T) {
	results := makeSummaryResults()
	stats := Summarize(results)

	if stats.TotalServices != 3 {
		t.Errorf("expected 3 total services, got %d", stats.TotalServices)
	}
	if stats.DriftedServices != 2 {
		t.Errorf("expected 2 drifted services, got %d", stats.DriftedServices)
	}
	if stats.CleanServices != 1 {
		t.Errorf("expected 1 clean service, got %d", stats.CleanServices)
	}
	if stats.TotalDiffs != 3 {
		t.Errorf("expected 3 total diffs, got %d", stats.TotalDiffs)
	}
}

func TestSummarize_AllClean(t *testing.T) {
	results := map[string][]CompareResult{
		"svc-x": {},
		"svc-y": {},
	}
	stats := Summarize(results)
	if stats.DriftedServices != 0 {
		t.Errorf("expected 0 drifted, got %d", stats.DriftedServices)
	}
	if stats.CleanServices != 2 {
		t.Errorf("expected 2 clean, got %d", stats.CleanServices)
	}
}

func TestFormatSummary_ContainsFields(t *testing.T) {
	stats := SummaryStats{TotalServices: 4, DriftedServices: 2, CleanServices: 2, TotalDiffs: 5}
	out := FormatSummary(stats)
	for _, want := range []string{"4", "2", "5"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in summary output: %s", want, out)
		}
	}
}
