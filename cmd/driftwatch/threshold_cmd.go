package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

// runThresholdCheck loads manifests, fetches live config, compares, and checks thresholds.
func runThresholdCheck(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: threshold-check <manifest-dir> <source-url> <thresholds-file>")
	}
	manifestDir, sourceURL, thresholdsFile := args[0], args[1], args[2]

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("load manifests: %w", err)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()

	var results []drift.CompareResult
	for _, m := range manifests {
		live, err := fetcher.Fetch(m.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: fetch %s: %v\n", m.Name, err)
			continue
		}
		res := detector.Compare(m, live)
		results = append(results, res)
	}

	tl, err := drift.LoadThresholds(thresholdsFile)
	if err != nil {
		return fmt.Errorf("load thresholds: %w", err)
	}

	violations := drift.CheckThresholds(results, tl)
	if len(violations) == 0 {
		fmt.Println("All services within threshold limits.")
		return nil
	}

	fmt.Printf("Threshold violations (%d):\n", len(violations))
	for _, v := range violations {
		fmt.Printf("  service=%-20s drifts=%-4d severity=%-8s max_allowed=%d min_severity=%s\n",
			v.Service, v.Drifts, v.Severity, v.Rule.MaxDrifts, v.Rule.MinSeverity)
	}
	return fmt.Errorf("%d threshold violation(s) detected", len(violations))
}
