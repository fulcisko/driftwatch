package drift

import (
	"strings"
	"testing"
)

func makeCorrelationResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30s", Actual: "60s"},
				{Key: "replicas", Expected: "2", Actual: "3"},
			},
		},
		{
			Service: "beta",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30s", Actual: "90s"},
				{Key: "log_level", Expected: "info", Actual: "debug"},
			},
		},
		{
			Service: "gamma",
			Diffs: []DiffEntry{
				{Key: "log_level", Expected: "info", Actual: "warn"},
			},
		},
		{
			Service: "clean",
			Diffs:   []DiffEntry{},
		},
	}
}

func TestBuildCorrelation_FindsPairs(t *testing.T) {
	results := makeCorrelationResults()
	report := BuildCorrelation(results)
	if len(report.Pairs) == 0 {
		t.Fatal("expected at least one correlated pair")
	}
	found := false
	for _, p := range report.Pairs {
		if p.ServiceA == "alpha" && p.ServiceB == "beta" {
			found = true
			if p.SharedCount != 1 {
				t.Errorf("expected 1 shared key, got %d", p.SharedCount)
			}
			if len(p.SharedKeys) < 1 || p.SharedKeys[0] != "timeout" {
				t.Errorf("expected shared key 'timeout', got %v", p.SharedKeys)
			}
		}
	}
	if !found {
		t.Error("expected alpha <-> beta pair")
	}
}

func TestBuildCorrelation_SkipsClean(t *testing.T) {
	results := makeCorrelationResults()
	report := BuildCorrelation(results)
	for _, p := range report.Pairs {
		if p.ServiceA == "clean" || p.ServiceB == "clean" {
			t.Error("clean service should not appear in correlation pairs")
		}
	}
}

func TestBuildCorrelation_NoDrift(t *testing.T) {
	results := []CompareResult{
		{Service: "a", Diffs: []DiffEntry{}},
		{Service: "b", Diffs: []DiffEntry{}},
	}
	report := BuildCorrelation(results)
	if len(report.Pairs) != 0 {
		t.Errorf("expected no pairs, got %d", len(report.Pairs))
	}
}

func TestFormatCorrelation_ContainsPair(t *testing.T) {
	results := makeCorrelationResults()
	report := BuildCorrelation(results)
	out := FormatCorrelation(report)
	if !strings.Contains(out, "alpha") {
		t.Error("expected output to contain 'alpha'")
	}
	if !strings.Contains(out, "timeout") {
		t.Error("expected output to contain 'timeout'")
	}
}

func TestFormatCorrelation_Empty(t *testing.T) {
	report := CorrelationReport{}
	out := FormatCorrelation(report)
	if !strings.Contains(out, "No correlated") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
