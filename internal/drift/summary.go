package drift

import "fmt"

// SummaryStats holds aggregate counts from a drift detection run.
type SummaryStats struct {
	TotalServices  int
	DriftedServices int
	CleanServices  int
	TotalDiffs     int
}

// Summarize computes aggregate statistics from a map of compare results.
func Summarize(results map[string][]CompareResult) SummaryStats {
	stats := SummaryStats{
		TotalServices: len(results),
	}
	for _, diffs := range results {
		if len(diffs) > 0 {
			stats.DriftedServices++
			stats.TotalDiffs += len(diffs)
		} else {
			stats.CleanServices++
		}
	}
	return stats
}

// FormatSummary returns a human-readable summary string.
func FormatSummary(stats SummaryStats) string {
	return fmt.Sprintf(
		"Services checked: %d | Drifted: %d | Clean: %d | Total diffs: %d",
		stats.TotalServices,
		stats.DriftedServices,
		stats.CleanServices,
		stats.TotalDiffs,
	)
}
