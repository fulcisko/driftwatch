package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

func runAlertGenerate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: alert-generate <manifest-dir> <source-url> <output-file> [min-severity]")
	}
	manifestDir, sourceURL, outputFile := args[0], args[1], args[2]
	minSeverity := drift.AlertHigh
	if len(args) >= 4 {
		minSeverity = drift.AlertSeverity(args[3])
	}

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("load manifests: %w", err)
	}

	fetcher := source.NewFetcher(nil)
	detector := drift.NewDetector()
	var results []drift.CompareResult
	for _, m := range manifests {
		live, err := fetcher.Fetch(sourceURL + "/" + m.Name)
		if err != nil {
			return fmt.Errorf("fetch %s: %w", m.Name, err)
		}
		res := detector.Compare(m, live)
		results = append(results, res)
	}

	alerts := drift.GenerateAlerts(results, drift.AlertConfig{MinSeverity: minSeverity})
	if err := drift.SaveAlerts(outputFile, alerts); err != nil {
		return fmt.Errorf("save alerts: %w", err)
	}
	fmt.Fprintf(os.Stdout, "generated %d alert(s) -> %s\n", len(alerts), outputFile)
	return nil
}

func runAlertShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: alert-show <alert-file>")
	}
	alerts, err := drift.LoadAlerts(args[0])
	if err != nil {
		return fmt.Errorf("load alerts: %w", err)
	}
	if len(alerts) == 0 {
		fmt.Println("no alerts found")
		return nil
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(alerts)
}
