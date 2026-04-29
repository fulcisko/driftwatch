package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeMaturityResults() []CompareResult {
	return []CompareResult{
		{
			Service: "svc-stable",
			Diffs: []DiffEntry{
				{Key: "log_level", Expected: "info", Actual: "debug"},
			},
		},
		{
			Service: "svc-clean",
			Diffs:   []DiffEntry{},
		},
		{
			Service: "svc-unstable",
			Diffs: []DiffEntry{
				{Key: "db_password", Expected: "secret", Actual: "changed"},
				{Key: "api_key", Expected: "abc", Actual: "xyz"},
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
	}
}

func TestAssessMaturity_CleanIsMature(t *testing.T) {
	results := makeMaturityResults()
	entries := AssessMaturity(results)
	for _, e := range entries {
		if e.Service == "svc-clean" {
			if e.Level != MaturityMature {
				t.Errorf("expected mature, got %s", e.Level)
			}
			return
		}
	}
	t.Error("svc-clean not found in entries")
}

func TestAssessMaturity_HighDriftIsUnstable(t *testing.T) {
	results := makeMaturityResults()
	entries := AssessMaturity(results)
	for _, e := range entries {
		if e.Service == "svc-unstable" {
			if e.Level == MaturityMature || e.Level == MaturityStable {
				t.Errorf("expected unstable/developing/unknown, got %s", e.Level)
			}
			return
		}
	}
	t.Error("svc-unstable not found in entries")
}

func TestAssessMaturity_AllServicesPresent(t *testing.T) {
	results := makeMaturityResults()
	entries := AssessMaturity(results)
	if len(entries) != len(results) {
		t.Errorf("expected %d entries, got %d", len(results), len(entries))
	}
}

func TestSaveAndLoadMaturityReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "maturity.json")
	results := makeMaturityResults()
	entries := AssessMaturity(results)

	if err := SaveMaturityReport(path, entries); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadMaturityReport(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(loaded))
	}
}

func TestLoadMaturityReport_NotFound(t *testing.T) {
	entries, err := LoadMaturityReport("/nonexistent/maturity.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d", len(entries))
	}
}

func TestFormatMaturity_ContainsServiceName(t *testing.T) {
	results := makeMaturityResults()
	entries := AssessMaturity(results)
	out := FormatMaturity(entries)
	if !containsStr(out, "svc-clean") {
		t.Error("expected svc-clean in output")
	}
}

func TestFormatMaturity_Empty(t *testing.T) {
	out := FormatMaturity([]MaturityEntry{})
	if out == "" {
		t.Error("expected non-empty output for empty entries")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstr(s, sub))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestScoreToMaturity_Zero(t *testing.T) {
	if scoreToMaturity(0) != MaturityMature {
		t.Error("score 0 should be mature")
	}
}

func TestScoreToMaturity_High(t *testing.T) {
	if scoreToMaturity(50) != MaturityUnknown {
		t.Error("score 50 should be unknown")
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
