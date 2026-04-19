package drift

import "fmt"

// RollupEntry summarizes drift for a single service.
type RollupEntry struct {
	Service    string         `json:"service"`
	TotalDiffs int            `json:"total_diffs"`
	BySeverity map[string]int `json:"by_severity"`
	TopKeys    []string       `json:"top_keys"`
}

// RollupReport holds rollup entries for all services.
type RollupReport struct {
	Entries []RollupEntry `json:"entries"`
	Total   int           `json:"total"`
}

// BuildRollup aggregates CompareResults into a RollupReport.
func BuildRollup(results []CompareResult) RollupReport {
	report := RollupReport{}
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		entry := RollupEntry{
			Service:    r.Service,
			TotalDiffs: len(r.Diffs),
			BySeverity: map[string]int{},
		}
		seen := map[string]bool{}
		for _, d := range r.Diffs {
			lvl := ClassifyKey(d.Key)
			entry.BySeverity[lvl.String()]++
			if !seen[d.Key] && len(entry.TopKeys) < 5 {
				entry.TopKeys = append(entry.TopKeys, d.Key)
				seen[d.Key] = true
			}
		}
		report.Entries = append(report.Entries, entry)
		report.Total += len(r.Diffs)
	}
	return report
}

// FormatRollup returns a human-readable rollup string.
func FormatRollup(r RollupReport) string {
	if len(r.Entries) == 0 {
		return "No drift detected across all services.\n"
	}
	out := fmt.Sprintf("Rollup: %d total diffs across %d service(s)\n", r.Total, len(r.Entries))
	for _, e := range r.Entries {
		out += fmt.Sprintf("  %s: %d diffs", e.Service, e.TotalDiffs)
		if h := e.BySeverity["high"]; h > 0 {
			out += fmt.Sprintf(" [high:%d]", h)
		}
		if m := e.BySeverity["medium"]; m > 0 {
			out += fmt.Sprintf(" [medium:%d]", m)
		}
		out += "\n"
	}
	return out
}
