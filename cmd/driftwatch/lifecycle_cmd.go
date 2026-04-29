package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

// runLifecycleSet sets the lifecycle stage for a service.
// Usage: driftwatch lifecycle set <path> <service> <stage> [note]
func runLifecycleSet(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: lifecycle set <path> <service> <stage> [note]")
	}
	path := args[1]
	service := args[2]
	stage := drift.LifecycleStage(args[3])
	note := ""
	if len(args) >= 5 {
		note = args[4]
	}
	if err := drift.SetLifecycle(path, service, stage, note); err != nil {
		return fmt.Errorf("set lifecycle: %w", err)
	}
	fmt.Fprintf(os.Stdout, "lifecycle stage for %q set to %q\n", service, stage)
	return nil
}

// runLifecycleShow prints all lifecycle entries, optionally filtered by stage.
// Usage: driftwatch lifecycle show <path> [stage]
func runLifecycleShow(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: lifecycle show <path> [stage]")
	}
	path := args[1]
	store, err := drift.LoadLifecycle(path)
	if err != nil {
		return fmt.Errorf("load lifecycle: %w", err)
	}
	entries := store.Entries
	if len(args) >= 3 {
		stage := drift.LifecycleStage(args[2])
		entries = drift.FilterByStage(store, stage)
	}
	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "no lifecycle entries found")
		return nil
	}
	for _, e := range entries {
		line := fmt.Sprintf("%-30s %-12s %s", e.Service, e.Stage, e.UpdatedAt.Format("2006-01-02"))
		if e.Note != "" {
			line += "  (" + e.Note + ")"
		}
		fmt.Fprintln(os.Stdout, line)
	}
	return nil
}
