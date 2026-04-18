package drift

import (
	"testing"
)

func makeCompareResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs: []Diff{
				{Key: "replicas", Expected: "3", Actual: "2", Kind: KindChanged},
				{Key: "image", Expected: "v1", Actual: "v2", Kind: KindChanged},
			},
		},
		{
			Service: "beta",
			Diffs: []Diff{},
		},
		{
			Service: "gamma",
			Diffs: []Diff{
				{Key: "timeout", Expected: "30", Actual: "60", Kind: KindChanged},
			},
		},
	}
}

func TestFilter_OnlyDrifted(t *testing.T) {
	f := NewFilterOptions()
	f.OnlyDrifted = true
	results := f.ApplyToResults(makeCompareResults())
	if len(results) != 2 {
		t.Fatalf("expected 2 drifted services, got %d", len(results))
	}
	for _, r := range results {
		if len(r.Diffs) == 0 {
			t.Errorf("service %q has no diffs but was included", r.Service)
		}
	}
}

func TestFilter_ServicePrefix(t *testing.T) {
	f := NewFilterOptions()
	f.ServicePrefix = "al"
	results := f.ApplyToResults(makeCompareResults())
	if len(results) != 1 || results[0].Service != "alpha" {
		t.Fatalf("expected only alpha, got %+v", results)
	}
}

func TestFilter_IgnoreKeys(t *testing.T) {
	f := NewFilterOptions()
	f.AddIgnoreKey("replicas")
	results := f.ApplyToResults(makeCompareResults())
	for _, r := range results {
		for _, d := range r.Diffs {
			if d.Key == "replicas" {
				t.Errorf("service %q still contains ignored key 'replicas'", r.Service)
			}
		}
	}
}

func TestFilter_ShouldIgnoreKey(t *testing.T) {
	f := NewFilterOptions()
	f.AddIgnoreKey("env")
	if !f.ShouldIgnoreKey("env") {
		t.Error("expected 'env' to be ignored")
	}
	if f.ShouldIgnoreKey("image") {
		t.Error("expected 'image' not to be ignored")
	}
}

func TestFilter_NoOptions(t *testing.T) {
	f := NewFilterOptions()
	results := f.ApplyToResults(makeCompareResults())
	if len(results) != 3 {
		t.Fatalf("expected all 3 results, got %d", len(results))
	}
}
