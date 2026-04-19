package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/manifest"
	"github.com/example/driftwatch/internal/source"
)

// runRollup loads manifests, fetches live config, runs comparison, and prints a rollup.
func runRollup(args []string, format string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: driftwatch rollup <manifest-dir> <source-url>")
	}
	manifestDir := args[0]
	sourceURL := args[1]

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("loading manifests: %w", err)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()

	var results []drift.CompareResult
	for _, m := range manifests {
		live, err := fetcher.Fetch(m.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: could not fetch %s: %v\n", m.Name, err)
			continue
		}
		res := detector.Compare(m, live)
		results = append(results, res)
	}

	rollup := drift.BuildRollup(results)

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rollup)
	default:
		fmt.Print(drift.FormatRollup(rollup))
	}
	return nil
}
