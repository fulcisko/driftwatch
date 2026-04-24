package drift

import (
	"strings"
	"testing"
)

func makeImpactResults() []CompareResult {
	return []CompareResult{
		{
			Service: "auth-service",
			Diffs: []DiffEntry{
				{Key: "secret_key", Expected: "abc", Actual: "xyz"},
				{Key: "password", Expected: "old", Actual: "new"},
				{Key: "token", Expected: "t1", Actual: "t2"},
			},
		},
		{
			Service: "api-gateway",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
		{
			Service: "clean-service",
			Diffs:   []DiffEntry{},
		},
	}
}

func TestAssessImpact_SkipsClean(t *testing.T) {
	results := makeImpactResults()
	reports := AssessImpact(results)
	for _, r := range reports {
		if r.Service == "clean-service" {
			t.Errorf("expected clean-service to be skipped")
		}
	}
}

func TestAssessImpact_SortedByScore(t *testing.T) {
	results := makeImpactResults()
	reports := AssessImpact(results)
	if len(reports) < 2 {
		t.Fatalf("expected at least 2 reports, got %d", len(reports))
	}
	if reports[0].Score < reports[1].Score {
		t.Errorf("reports not sorted descending by score: %d < %d", reports[0].Score, reports[1].Score)
	}
}

func TestAssessImpact_CriticalForHighKeys(t *testing.T) {
	results := []CompareResult{
		{
			Service: "svc",
			Diffs: []DiffEntry{
				{Key: "secret_key", Expected: "a", Actual: "b"},
				{Key: "password", Expected: "a", Actual: "b"},
				{Key: "token", Expected: "a", Actual: "b"},
			},
		},
	}
	reports := AssessImpact(results)
	if len(reports) != 1 {
		t.Fatalf("expected 1 report")
	}
	if reports[0].Impact != ImpactCritical {
		t.Errorf("expected critical impact, got %s (score=%d)", reports[0].Impact, reports[0].Score)
	}
}

func TestAssessImpact_TopKeysLimit(t *testing.T) {
	diffs := make([]DiffEntry, 8)
	for i := range diffs {
		diffs[i] = DiffEntry{Key: fmt.Sprintf("key_%d", i), Expected: "a", Actual: "b"}
	}
	results := []CompareResult{{Service: "svc", Diffs: diffs}}
	reports := AssessImpact(results)
	if len(reports[0].TopKeys) > 5 {
		t.Errorf("expected at most 5 top keys, got %d", len(reports[0].TopKeys))
	}
}

func TestFormatImpact_ContainsServiceName(t *testing.T) {
	reports := AssessImpact(makeImpactResults())
	out := FormatImpact(reports)
	if !strings.Contains(out, "auth-service") {
		t.Errorf("expected output to contain auth-service")
	}
}

func TestFormatImpact_EmptyReports(t *testing.T) {
	out := FormatImpact(nil)
	if !strings.Contains(out, "No drift impact") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
