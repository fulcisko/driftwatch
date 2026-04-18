package main

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/manifest"
	"github.com/driftwatch/internal/source"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]
	if len(args) < 2 {
		return fmt.Errorf("usage: driftwatch <manifest-dir> <service-base-url>")
	}

	manifestDir := args[0]
	baseURL := args[1]

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("loading manifests: %w", err)
	}

	if len(manifests) == 0 {
		return fmt.Errorf("no manifests found in %s", manifestDir)
	}

	fetcher := source.NewFetcher(baseURL)
	detector := drift.NewDetector()
	reporter := drift.NewReporter(os.Stdout)

	var results []drift.Result

	for _, m := range manifests {
		live, err := fetcher.Fetch(m.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: could not fetch config for %s: %v\n", m.Name, err)
			continue
		}

		res := detector.Compare(m, live)
		results = append(results, res)
	}

	reporter.Print(results)

	if reporter.Summary(results) > 0 {
		os.Exit(2)
	}

	return nil
}
