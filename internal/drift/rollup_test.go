package drift

import (
	"strings"
	"testing"
)

func makeRollupResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs: []Diff{
				{Key: "secret_key", Expected: "x", Actual: "y"},
				{Key: "replicas", Expected: "2", Actual: "3"},
			},
		},
		{
			Service: "beta",
			Diffs:   []Diff{},
		},
		{
			Service: "gamma",
			Diffs: []Diff{
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
	}
}

func TestBuildRollup_SkipsClean(t *testing.T) {
	r := BuildRollup(makeRollupResults())
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestBuildRollup_TotalDiffs(t *testing.T) {
	r := BuildRollup(makeRollupResults())
	if r.Total != 3 {
		t.Errorf("expected total 3, got %d", r.Total)
	}
}

func TestBuildRollup_BySeverity(t *testing.T) {
	r := BuildRollup(makeRollupResults())
	var alpha RollupEntry
	for _, e := range r.Entries {
		if e.Service == "alpha" {
			alpha = e
		}
	}
	if alpha.BySeverity["high"] != 1 {
		t.Errorf("expected 1 high severity diff for alpha, got %d", alpha.BySeverity["high"])
	}
}

func TestFormatRollup_ContainsServiceName(t *testing.T) {
	r := BuildRollup(makeRollupResults())
	out := FormatRollup(r)
	if !strings.Contains(out, "alpha") {
		t.Error("expected output to contain 'alpha'")
	}
}

func TestFormatRollup_NoDrift(t *testing.T) {
	r := BuildRollup([]CompareResult{{Service: "clean", Diffs: []Diff{}}})
	out := FormatRollup(r)
	if !strings.Contains(out, "No drift") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}
