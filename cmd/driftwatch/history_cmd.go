package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

func runShowHistory(historyPath string, out io.Writer) error {
	entries, err := drift.LoadHistory(historyPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(out, "No history found.")
			return nil
		}
		return fmt.Errorf("load history: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(out, "History is empty.")
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tTIMESTAMP\tSERVICES\tDRIFTED")
	for i, e := range entries {
		total := len(e.Results)
		drifted := 0
		for _, r := range e.Results {
			if len(r.Diffs) > 0 {
				drifted++
			}
		}
		fmt.Fprintf(w, "%d\t%s\t%d\t%d\n",
			i+1,
			e.Timestamp.Format(time.RFC3339),
			total,
			drifted,
		)
	}
	return w.Flush()
}

func runAppendHistory(historyPath string, results []drift.CompareResult) error {
	if err := drift.AppendHistory(historyPath, results); err != nil {
		return fmt.Errorf("append history: %w", err)
	}
	fmt.Printf("History saved to %s (%d services)\n", historyPath, len(results))
	return nil
}
