package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeThresholdFile(t *testing.T, tl ThresholdList) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "thresholds.json")
	data, _ := json.Marshal(tl)
	os.WriteFile(path, data, 0644)
	return path
}

func makeThresholdResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []DiffEntry{
				{Key: "replicas", Expected: "3", Actual: "1", Status: "changed"},
				{Key: "image", Expected: "v2", Actual: "v1", Status: "changed"},
			},
		},
		{
			Service: "worker",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30", Actual: "30", Status: "match"},
			},
		},
	}
}

func TestLoadThresholds_NotFound(t *testing.T) {
	tl, err := LoadThresholds("/nonexistent/thresholds.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tl.Rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestLoadThresholds_Valid(t *testing.T) {
	input := ThresholdList{Rules: []ThresholdRule{
		{Service: "api", MaxDrifts: 1, MinSeverity: "high"},
	}}
	path := writeThresholdFile(t, input)
	tl, err := LoadThresholds(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tl.Rules) != 1 || tl.Rules[0].Service != "api" {
		t.Errorf("unexpected rules: %+v", tl.Rules)
	}
}

func TestCheckThresholds_Violation(t *testing.T) {
	tl := ThresholdList{Rules: []ThresholdRule{
		{Service: "api", MaxDrifts: 1, MinSeverity: "high"},
	}}
	results := makeThresholdResults()
	violations := CheckThresholds(results, tl)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Service != "api" {
		t.Errorf("expected api violation")
	}
}

func TestCheckThresholds_NoViolation(t *testing.T) {
	tl := ThresholdList{Rules: []ThresholdRule{
		{Service: "worker", MaxDrifts: 5, MinSeverity: "high"},
	}}
	results := makeThresholdResults()
	violations := CheckThresholds(results, tl)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheckThresholds_NoMatchingRule(t *testing.T) {
	tl := ThresholdList{Rules: []ThresholdRule{
		{Service: "unknown", MaxDrifts: 0, MinSeverity: "low"},
	}}
	results := makeThresholdResults()
	violations := CheckThresholds(results, tl)
	if len(violations) != 0 {
		t.Errorf("expected no violations for unmatched service")
	}
}
