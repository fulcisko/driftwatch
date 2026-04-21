package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/manifest"
	"github.com/driftwatch/internal/source"
)

// runWindowCheck loads manifests, fetches live config, compares, then filters
// results to services seen in history within the given hour window.
func runWindowCheck(args []string, historyPath string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: window <manifest-dir> <source-url> <hours>")
	}

	manifestDir := args[0]
	sourceURL := args[1]
	hours, err := strconv.Atoi(args[2])
	if err != nil || hours <= 0 {
		return fmt.Errorf("hours must be a positive integer, got: %s", args[2])
	}

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
		r := detector.Compare(m, live)
		results = append(results, r)
	}

	history, err := drift.LoadHistory(historyPath)
	if err != nil {
		history = nil
	}

	opts := drift.NewWindowOptions(time.Duration(hours) * time.Hour)
	windowed := drift.ApplyWindow(results, history, opts)

	fmt.Printf("Window: %s\n", drift.FormatWindow(opts))
	fmt.Printf("Services in window: %d\n", len(windowed))
	for _, r := range windowed {
		if len(r.Diffs) == 0 {
			fmt.Printf("  [clean]  %s\n", r.Service)
		} else {
			fmt.Printf("  [drift]  %s (%d diffs)\n", r.Service, len(r.Diffs))
		}
	}
	return nil
}
