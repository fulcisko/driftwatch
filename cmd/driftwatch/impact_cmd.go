package main

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/manifest"
	"github.com/example/driftwatch/internal/source"
)

// runImpact loads manifests and live config, runs drift detection, and prints
// an impact assessment ranked by severity score.
func runImpact(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: driftwatch impact <manifest-dir> <source-url>")
	}

	manifestDir := args[0]
	sourceURL := args[1]

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("loading manifests: %w", err)
	}
	if len(manifests) == 0 {
		return fmt.Errorf("no manifests found in %s", manifestDir)
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
		cr := detector.Compare(m, live)
		results = append(results, cr)
	}

	reports := drift.AssessImpact(results)
	fmt.Print(drift.FormatImpact(reports))
	return nil
}
