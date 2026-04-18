package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

// runSaveBaseline saves current drift results as a baseline file.
func runSaveBaseline(resultsPath, baselinePath string) error {
	b, err := drift.LoadBaseline(resultsPath)
	if err != nil {
		return fmt.Errorf("load results: %w", err)
	}
	if err := drift.SaveBaseline(baselinePath, b.Results); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Baseline saved to %s (%d services)\n", baselinePath, len(b.Results))
	return nil
}

// runDiffBaseline compares current results against a saved baseline.
func runDiffBaseline(baselinePath, currentPath string) error {
	baseline, err := drift.LoadBaseline(baselinePath)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}
	current, err := drift.LoadBaseline(currentPath)
	if err != nil {
		return fmt.Errorf("load current: %w", err)
	}
	changed := drift.DiffBaseline(baseline, current.Results)
	if len(changed) == 0 {
		fmt.Fprintln(os.Stdout, "No drift changes since baseline.")
		return nil
	}
	fmt.Fprintf(os.Stdout, "%d service(s) changed since baseline:\n", len(changed))
	for _, r := range changed {
		status := "clean"
		if r.HasDrift() {
			status = fmt.Sprintf("%d diff(s)", len(r.Diffs))
		}
		fmt.Fprintf(os.Stdout, "  %s: %s\n", r.Service, status)
	}
	return nil
}
