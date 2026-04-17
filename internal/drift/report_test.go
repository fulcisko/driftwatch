package drift_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeResults() []drift.DriftResult {
	return []drift.DriftResult{
		{ServiceName: "api", HasDrift: false},
		{
			ServiceName: "worker",
			HasDrift:    true,
			Diffs: []drift.FieldDiff{
				{Field: "replicas", Expected: 3, Actual: 1},
			},
		},
	}
}

func TestPrint_ContainsServiceNames(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)
	r.Print(makeResults())
	out := buf.String()
	if !strings.Contains(out, "api") {
		t.Error("expected 'api' in output")
	}
	if !strings.Contains(out, "worker") {
		t.Error("expected 'worker' in output")
	}
	if !strings.Contains(out, "[DRIFT]") {
		t.Error("expected [DRIFT] marker in output")
	}
	if !strings.Contains(out, "[OK]") {
		t.Error("expected [OK] marker in output")
	}
}

func TestSummary_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)
	summary := r.Summary(makeResults())
	if !strings.Contains(summary, "1/2") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestSummary_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)
	results := []drift.DriftResult{
		{ServiceName: "api", HasDrift: false},
	}
	summary := r.Summary(results)
	if !strings.Contains(summary, "in sync") {
		t.Errorf("unexpected summary: %s", summary)
	}
}
