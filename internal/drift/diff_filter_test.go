package drift

import (
	"testing"
)

func makeDiffResults() []CompareResult {
	return []CompareResult{
		{Service: "api", Diffs: []DiffEntry{
			{Key: "SECRET_KEY", Expected: "x", Actual: "y"},
			{Key: "LOG_LEVEL", Expected: "info", Actual: "debug"},
		}},
		{Service: "worker", Diffs: []DiffEntry{
			{Key: "TIMEOUT", Expected: "30", Actual: "60"},
		}},
		{Service: "clean-svc", Diffs: []DiffEntry{}},
	}
}

func TestApplyDiffFilter_OnlyDrifted(t *testing.T) {
	results := makeDiffResults()
	out := ApplyDiffFilter(results, DiffFilterOptions{OnlyDrifted: true})
	for _, r := range out {
		if len(r.Diffs) == 0 {
			t.Errorf("expected only drifted services, got clean: %s", r.Service)
		}
	}
	if len(out) != 2 {
		t.Errorf("expected 2 drifted services, got %d", len(out))
	}
}

func TestApplyDiffFilter_ServicePrefix(t *testing.T) {
	results := makeDiffResults()
	out := ApplyDiffFilter(results, DiffFilterOptions{ServicePrefix: "api"})
	if len(out) != 1 || out[0].Service != "api" {
		t.Errorf("expected only api service, got %+v", out)
	}
}

func TestApplyDiffFilter_ExcludeKeys(t *testing.T) {
	results := makeDiffResults()
	out := ApplyDiffFilter(results, DiffFilterOptions{ExcludeKeys: []string{"SECRET_KEY"}})
	for _, r := range out {
		if r.Service == "api" {
			for _, d := range r.Diffs {
				if d.Key == "SECRET_KEY" {
					t.Error("SECRET_KEY should be excluded")
				}
			}
		}
	}
}

func TestApplyDiffFilter_MinSeverity(t *testing.T) {
	results := makeDiffResults()
	out := ApplyDiffFilter(results, DiffFilterOptions{MinSeverity: SeverityHigh, OnlyDrifted: true})
	for _, r := range out {
		for _, d := range r.Diffs {
			if ClassifyKey(d.Key) < SeverityHigh {
				t.Errorf("key %s below min severity", d.Key)
			}
		}
	}
}
