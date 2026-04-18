package drift

import (
	"encoding/json"
	"fmt"
	"io"
)

// ExportJSON writes the drift results and summary as JSON to w.
func ExportJSON(w io.Writer, results map[string][]CompareResult, stats SummaryStats) error {
	type payload struct {
		Summary SummaryStats                `json:"summary"`
		Services map[string][]CompareResult `json:"services"`
	}
	p := payload{
		Summary:  stats,
		Services: results,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(p); err != nil {
		return fmt.Errorf("export json: %w", err)
	}
	return nil
}

// ExportText writes a plain-text drift report to w.
func ExportText(w io.Writer, results map[string][]CompareResult, stats SummaryStats) error {
	for svc, diffs := range results {
		if len(diffs) == 0 {
			continue
		}
		fmt.Fprintf(w, "[%s]\n", svc)
		for _, d := range diffs {
			fmt.Fprintf(w, "  key=%s kind=%s expected=%q actual=%q\n",
				d.Key, d.Kind, d.Expected, d.Actual)
		}
	}
	fmt.Fprintln(w, FormatSummary(stats))
	return nil
}
