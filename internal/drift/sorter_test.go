package drift

import (
	"testing"
)

func makeSortResults() []CompareResult {
	return []CompareResult{
		{
			Service: "zebra-svc",
			Diffs: []DiffEntry{
				{Key: "port", Kind: DiffChanged, Expected: "8080", Actual: "9090"},
			},
		},
		{
			Service: "alpha-svc",
			Diffs: []DiffEntry{
				{Key: "host", Kind: DiffMissing, Expected: "localhost", Actual: ""},
				{Key: "timeout", Kind: DiffMissing, Expected: "30s", Actual: ""},
			},
		},
		{
			Service: "beta-svc",
			Diffs:   []DiffEntry{},
		},
	}
}

func TestSortResults_ByService(t *testing.T) {
	results := makeSortResults()
	sorted := SortResults(results, SortByService)

	if sorted[0].Service != "alpha-svc" || sorted[1].Service != "beta-svc" || sorted[2].Service != "zebra-svc" {
		t.Errorf("unexpected service order: %v, %v, %v", sorted[0].Service, sorted[1].Service, sorted[2].Service)
	}
}

func TestSortResults_ByDriftCount(t *testing.T) {
	results := makeSortResults()
	sorted := SortResults(results, SortByDriftCount)

	if sorted[0].Service != "alpha-svc" {
		t.Errorf("expected alpha-svc first (2 diffs), got %s", sorted[0].Service)
	}
	if sorted[2].Service != "beta-svc" {
		t.Errorf("expected beta-svc last (0 diffs), got %s", sorted[2].Service)
	}
}

func TestSortResults_BySeverity(t *testing.T) {
	results := makeSortResults()
	sorted := SortResults(results, SortBySeverity)

	// zebra-svc has 1 Changed (score=3), alpha-svc has 2 Missing (score=4)
	if sorted[0].Service != "alpha-svc" {
		t.Errorf("expected alpha-svc first by severity, got %s", sorted[0].Service)
	}
}

func TestSortResults_DoesNotMutateOriginal(t *testing.T) {
	results := makeSortResults()
	originalFirst := results[0].Service
	SortResults(results, SortByService)

	if results[0].Service != originalFirst {
		t.Error("original slice was mutated")
	}
}
