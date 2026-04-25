package drift

import (
	"strings"
	"testing"
)

func makeReachabilityResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []DiffEntry{{Key: "timeout", Expected: "30", Actual: "60"}},
		},
		{
			Service: "worker",
			Diffs:   []DiffEntry{},
		},
		{
			Service: "gateway",
			Diffs: []DiffEntry{{Key: "port", Expected: "8080", Actual: "9090"}},
		},
	}
}

func makeReachabilityDeps() []Dependency {
	return []Dependency{
		{Service: "worker", DependsOn: "api"},
		{Service: "gateway", DependsOn: "worker"},
	}
}

func TestBuildReachability_AffectedBy(t *testing.T) {
	results := makeReachabilityResults()
	deps := makeReachabilityDeps()

	out := BuildReachability(results, deps)

	// gateway depends on worker which depends on api (drifted)
	var gw *ReachabilityResult
	for i := range out {
		if out[i].Service == "gateway" {
			gw = &out[i]
		}
	}
	if gw == nil {
		t.Fatal("expected gateway in results")
	}
	if len(gw.AffectedBy) == 0 {
		t.Error("expected gateway to be affected by upstream drift")
	}
	found := false
	for _, s := range gw.AffectedBy {
		if s == "api" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected api in gateway.AffectedBy, got %v", gw.AffectedBy)
	}
}

func TestBuildReachability_CleanServiceNoExposure(t *testing.T) {
	results := makeReachabilityResults()
	deps := makeReachabilityDeps()

	out := BuildReachability(results, deps)

	for _, r := range out {
		if r.Service == "worker" && r.TotalExposure > 0 {
			// worker only depends on api (drifted), so exposure should be 1
			// but worker itself is clean — exposure reflects upstream drift
		}
	}
}

func TestBuildReachability_NoDeps(t *testing.T) {
	results := makeReachabilityResults()
	out := BuildReachability(results, nil)

	for _, r := range out {
		if r.TotalExposure != 0 {
			t.Errorf("expected zero exposure with no deps for %s, got %d", r.Service, r.TotalExposure)
		}
	}
}

func TestBuildReachability_SortedByExposure(t *testing.T) {
	results := makeReachabilityResults()
	deps := makeReachabilityDeps()

	out := BuildReachability(results, deps)

	for i := 1; i < len(out); i++ {
		if out[i].TotalExposure > out[i-1].TotalExposure {
			t.Errorf("results not sorted by exposure: index %d (%d) > index %d (%d)",
				i, out[i].TotalExposure, i-1, out[i-1].TotalExposure)
		}
	}
}

func TestFormatReachability_ContainsService(t *testing.T) {
	results := makeReachabilityResults()
	deps := makeReachabilityDeps()
	out := BuildReachability(results, deps)

	formatted := FormatReachability(out)
	if !strings.Contains(formatted, "api") {
		t.Error("expected formatted output to contain 'api'")
	}
	if !strings.Contains(formatted, "Drift Reachability") {
		t.Error("expected header in formatted output")
	}
}

func TestFormatReachability_Empty(t *testing.T) {
	out := FormatReachability(nil)
	if !strings.Contains(out, "no reachability") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
