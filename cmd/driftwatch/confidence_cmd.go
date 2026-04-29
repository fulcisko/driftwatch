package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/manifest"
	"github.com/example/driftwatch/internal/source"
)

// runConfidence loads manifests, fetches live config, compares them,
// optionally loads history, and prints confidence scores.
func runConfidence(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: driftwatch confidence <manifest-dir> <source-url> [history-file] [--json]")
	}

	manifestDir := args[0]
	sourceURL := args[1]

	var historyFile string
	jsonOut := false
	for _, a := range args[2:] {
		if a == "--json" {
			jsonOut = true
		} else {
			historyFile = a
		}
	}

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("loading manifests: %w", err)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()

	var compareResults []drift.CompareResult
	for _, m := range manifests {
		body, err := fetcher.Fetch(m.Name)
		if err != nil {
			return fmt.Errorf("fetching %s: %w", m.Name, err)
		}
		res := detector.Compare(m, body)
		compareResults = append(compareResults, res)
	}

	var history []drift.HistoryEntry
	if historyFile != "" {
		history, err = drift.LoadHistory(historyFile)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("loading history: %w", err)
		}
	}

	scores := drift.ScoreConfidence(compareResults, history)

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(scores)
	}

	fmt.Print(drift.FormatConfidence(scores))
	return nil
}
