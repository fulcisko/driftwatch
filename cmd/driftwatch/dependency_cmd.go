package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// runDependencyAdd adds a directed dependency edge between two services.
// Usage: driftwatch dependency add <from> <to> [label] --file <path>
func runDependencyAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: dependency add <from> <to> [label]")
	}
	from, to := args[0], args[1]
	label := ""
	if len(args) >= 3 {
		label = args[2]
	}
	path := envOr("DRIFTWATCH_DEP_FILE", "deps.json")
	if err := drift.AddDependency(path, from, to, label); err != nil {
		return fmt.Errorf("add dependency: %w", err)
	}
	fmt.Fprintf(os.Stdout, "dependency added: %s -> %s\n", from, to)
	return nil
}

// runDependencyShow prints all dependencies for a given service.
// Usage: driftwatch dependency show <service> [--direction dependents|dependencies]
func runDependencyShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: dependency show <service>")
	}
	service := args[0]
	direction := "dependencies"
	if len(args) >= 2 {
		direction = strings.ToLower(args[1])
	}
	path := envOr("DRIFTWATCH_DEP_FILE", "deps.json")
	graph, err := drift.LoadDependencies(path)
	if err != nil {
		return fmt.Errorf("load dependencies: %w", err)
	}
	var results []string
	switch direction {
	case "dependents":
		results = drift.DependentsOf(graph, service)
	default:
		results = drift.DependenciesOf(graph, service)
	}
	if len(results) == 0 {
		fmt.Fprintf(os.Stdout, "no %s found for %s\n", direction, service)
		return nil
	}
	fmt.Fprintf(os.Stdout, "%s of %s:\n", direction, service)
	for _, r := range results {
		fmt.Fprintf(os.Stdout, "  - %s\n", r)
	}
	return nil
}
