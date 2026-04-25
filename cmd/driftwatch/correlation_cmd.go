package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

// runCorrelation detects services sharing the same drifted config keys.
// Usage: driftwatch correlation <manifest-dir> <source-url> [--format=text|json]
func runCorrelation(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: correlation <manifest-dir> <source-url> [--format=text|json]")
	}

	manifestDir := args[0]
	sourceURL := args[1]
	format := "text"
	for _, a := range args[2:] {
		if a == "--format=json" {
			format = "json"
		} else if a == "--format=text" {
			format = "text"
		}
	}

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("loading manifests: %w", err)
	}
	if len(manifests) == 0 {
		fmt.Fprintln(os.Stdout, "No manifests found.")
		return nil
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()

	var results []drift.CompareResult
	for _, m := range manifests {
		live, err := fetcher.Fetch(m.Name)
		if err != nil {
			continue
		}
		res := detector.Compare(m, live)
		results = append(results, res)
	}

	report := drift.BuildCorrelation(results)

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	default:
		fmt.Print(drift.FormatCorrelation(report))
	}
	return nil
}
