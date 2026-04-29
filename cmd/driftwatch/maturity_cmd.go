package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/driftwatch/driftwatch/internal/drift"
	"github.com/driftwatch/driftwatch/internal/manifest"
	"github.com/driftwatch/driftwatch/internal/source"
)

// runMaturityAssess assesses config maturity for all services in a manifest dir.
// Usage: maturity assess <manifest-dir> <source-url> [--format=text|json] [--save=<path>]
func runMaturityAssess(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: maturity assess <manifest-dir> <source-url> [--format=text|json] [--save=<path>]")
	}
	manifestDir := args[0]
	sourceURL := args[1]
	format := envOr("DRIFTWATCH_FORMAT", "text")
	savePath := ""
	for _, a := range args[2:] {
		if len(a) > 9 && a[:9] == "--format=" {
			format = a[9:]
		}
		if len(a) > 7 && a[:7] == "--save=" {
			savePath = a[7:]
		}
	}

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("load manifests: %w", err)
	}
	if len(manifests) == 0 {
		return fmt.Errorf("no manifests found in %s", manifestDir)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()
	var results []drift.CompareResult
	for _, m := range manifests {
		body, err := fetcher.Fetch(m.Name)
		if err != nil {
			continue
		}
		res := detector.Compare(m, body)
		results = append(results, res)
	}

	entries := drift.AssessMaturity(results)

	if savePath != "" {
		if err := drift.SaveMaturityReport(savePath, entries); err != nil {
			return fmt.Errorf("save maturity report: %w", err)
		}
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(entries)
	default:
		fmt.Print(drift.FormatMaturity(entries))
	}
	return nil
}

// runMaturityShow loads and displays a saved maturity report.
// Usage: maturity show <path> [--format=text|json]
func runMaturityShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: maturity show <path> [--format=text|json]")
	}
	path := args[0]
	format := "text"
	for _, a := range args[1:] {
		if len(a) > 9 && a[:9] == "--format=" {
			format = a[9:]
		}
	}

	entries, err := drift.LoadMaturityReport(path)
	if err != nil {
		return fmt.Errorf("load maturity report: %w", err)
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(entries)
	default:
		fmt.Print(drift.FormatMaturity(entries))
	}
	return nil
}
