package drift

import (
	"fmt"
	"sort"
	"strings"
)

// HeatmapEntry represents drift frequency for a service+key pair.
type HeatmapEntry struct {
	Service string `json:"service"`
	Key     string `json:"key"`
	Count   int    `json:"count"`
	MaxSev  string `json:"max_severity"`
}

// HeatmapRow groups entries by service.
type HeatmapRow struct {
	Service string
	Entries []HeatmapEntry
	Total   int
}

// BuildHeatmap aggregates drift results into a frequency heatmap.
func BuildHeatmap(results []CompareResult) []HeatmapRow {
	type key struct{ service, field string }
	counts := map[key]int{}
	sevs := map[key]SeverityLevel{}

	for _, r := range results {
		for _, d := range r.Diffs {
			k := key{r.Service, d.Key}
			counts[k]++
			sev := ClassifyKey(d.Key)
			if sev > sevs[k] {
				sevs[k] = sev
			}
		}
	}

	rowMap := map[string]*HeatmapRow{}
	for k, count := range counts {
		row, ok := rowMap[k.service]
		if !ok {
			row = &HeatmapRow{Service: k.service}
			rowMap[k.service] = row
		}
		row.Entries = append(row.Entries, HeatmapEntry{
			Service: k.service,
			Key:     k.field,
			Count:   count,
			MaxSev:  sevs[k].String(),
		})
		row.Total += count
	}

	rows := make([]HeatmapRow, 0, len(rowMap))
	for _, row := range rowMap {
		sort.Slice(row.Entries, func(i, j int) bool {
			return row.Entries[i].Count > row.Entries[j].Count
		})
		rows = append(rows, *row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Total > rows[j].Total
	})
	return rows
}

// FormatHeatmap returns a human-readable heatmap table.
func FormatHeatmap(rows []HeatmapRow) string {
	if len(rows) == 0 {
		return "no drift data for heatmap\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-24s %-24s %6s  %s\n", "SERVICE", "KEY", "COUNT", "MAX_SEV"))
	sb.WriteString(strings.Repeat("-", 68) + "\n")
	for _, row := range rows {
		for _, e := range row.Entries {
			sb.WriteString(fmt.Sprintf("%-24s %-24s %6d  %s\n", e.Service, e.Key, e.Count, e.MaxSev))
		}
	}
	return sb.String()
}
