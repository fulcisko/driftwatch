package drift

import (
	"sort"
)

// SortOrder defines the ordering strategy for drift results.
type SortOrder string

const (
	SortByService  SortOrder = "service"
	SortByDriftCount SortOrder = "drift_count"
	SortBySeverity SortOrder = "severity"
)

// SortResults returns a new slice of CompareResult sorted by the given order.
func SortResults(results []CompareResult, order SortOrder) []CompareResult {
	copied := make([]CompareResult, len(results))
	copy(copied, results)

	switch order {
	case SortByDriftCount:
		sort.SliceStable(copied, func(i, j int) bool {
			return driftCount(copied[i]) > driftCount(copied[j])
		})
	case SortBySeverity:
		sort.SliceStable(copied, func(i, j int) bool {
			return severityScore(copied[i]) > severityScore(copied[j])
		})
	default: // SortByService
		sort.SliceStable(copied, func(i, j int) bool {
			return copied[i].Service < copied[j].Service
		})
	}

	return copied
}

func driftCount(r CompareResult) int {
	return len(r.Diffs)
}

// severityScore ranks results: changed values > missing fields > unexpected fields.
func severityScore(r CompareResult) int {
	score := 0
	for _, d := range r.Diffs {
		switch d.Kind {
		case DiffChanged:
			score += 3
		case DiffMissing:
			score += 2
		case DiffUnexpected:
			score += 1
		}
	}
	return score
}
