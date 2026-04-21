package drift

import (
	"fmt"
	"time"
)

// WindowOptions defines a time range for filtering drift results.
type WindowOptions struct {
	From time.Time
	To   time.Time
}

// WindowedResult holds drift results scoped to a time window.
type WindowedResult struct {
	Window  WindowOptions
	Results []CompareResult
}

// ApplyWindow filters a slice of CompareResults to those whose service
// appears in the history within the given time window.
func ApplyWindow(results []CompareResult, history []HistoryEntry, opts WindowOptions) []CompareResult {
	servicesInWindow := map[string]bool{}
	for _, h := range history {
		if !h.RecordedAt.Before(opts.From) && !h.RecordedAt.After(opts.To) {
			servicesInWindow[h.Service] = true
		}
	}

	var filtered []CompareResult
	for _, r := range results {
		if servicesInWindow[r.Service] {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// FormatWindow returns a human-readable description of the window.
func FormatWindow(opts WindowOptions) string {
	return fmt.Sprintf("from %s to %s",
		opts.From.Format(time.RFC3339),
		opts.To.Format(time.RFC3339),
	)
}

// NewWindowOptions constructs a WindowOptions from duration back from now.
func NewWindowOptions(duration time.Duration) WindowOptions {
	now := time.Now().UTC()
	return WindowOptions{
		From: now.Add(-duration),
		To:   now,
	}
}
