package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

// runScore loads manifests, fetches live config, and prints drift risk scores.
func runScore(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: driftwatch score <manifest-dir> <source-url>")
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
		r := detector.Compare(m, live)
		results = append(results, r)
	}

	scores := drift.ScoreResults(results)

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	for _, s := range scores {
		fmt.Println(drift.FormatScore(s))
	}
	return nil
}
