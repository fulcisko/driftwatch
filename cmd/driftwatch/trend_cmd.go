package main

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/manifest"
	"github.com/example/driftwatch/internal/source"
)

func runTrendAppend(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: trend append <manifest-dir> <source-url> <trend-file>")
	}
	manifestDir, sourceURL, trendFile := args[0], args[1], args[2]

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
		results = append(results, detector.Compare(m, live))
	}

	if err := drift.AppendTrend(trendFile, results); err != nil {
		return fmt.Errorf("append trend: %w", err)
	}
	fmt.Println("trend updated")
	return nil
}

func runTrendShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: trend show <trend-file> [service]")
	}
	trendFile := args[0]
	service := ""
	if len(args) >= 2 {
		service = args[1]
	}

	report, err := drift.LoadTrend(trendFile)
	if err != nil {
		return fmt.Errorf("load trend: %w", err)
	}

	points := drift.FilterTrend(report, service)
	fmt.Print(drift.FormatTrend(points))
	return nil
}
