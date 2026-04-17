package drift

import (
	"fmt"
	"io"
	"strings"
)

// Reporter formats drift results for output.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter writing to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Print writes a human-readable summary of results to the writer.
func (r *Reporter) Print(results []DriftResult) {
	for _, res := range results {
		if !res.HasDrift {
			fmt.Fprintf(r.w, "[OK]    %s — no drift detected\n", res.ServiceName)
			continue
		}
		fmt.Fprintf(r.w, "[DRIFT] %s — %d field(s) differ:\n", res.ServiceName, len(res.Diffs))
		for _, d := range res.Diffs {
			fmt.Fprintf(r.w, "  %-30s expected=%-20v actual=%v\n",
				d.Field,
				formatVal(d.Expected),
				formatVal(d.Actual),
			)
		}
	}
}

// Summary returns a one-line summary string.
func (r *Reporter) Summary(results []DriftResult) string {
	total := len(results)
	drifted := 0
	for _, res := range results {
		if res.HasDrift {
			drifted++
		}
	}
	if drifted == 0 {
		return fmt.Sprintf("All %d service(s) in sync.", total)
	}
	return fmt.Sprintf("%d/%d service(s) have drift.", drifted, total)
}

func formatVal(v interface{}) string {
	if v == nil {
		return "<absent>"
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}
