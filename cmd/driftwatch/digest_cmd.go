package main

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/manifest"
	"github.com/example/driftwatch/internal/source"
)

// runDigestBuild loads manifests, fetches live config, computes digests, and
// saves them to the given output path.
func runDigestBuild(manifestDir, sourceURL, outputPath string) error {
	if manifestDir == "" || sourceURL == "" || outputPath == "" {
		return fmt.Errorf("usage: digest build <manifest-dir> <source-url> <output-path>")
	}

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

	entries := drift.BuildDigests(results)
	if err := drift.SaveDigests(outputPath, entries); err != nil {
		return fmt.Errorf("save digests: %w", err)
	}
	fmt.Printf("saved %d digest(s) to %s\n", len(entries), outputPath)
	return nil
}

// runDigestDiff loads a previously saved digest file and compares it against
// freshly computed digests, printing services whose drift has changed.
func runDigestDiff(manifestDir, sourceURL, baselinePath string) error {
	if manifestDir == "" || sourceURL == "" || baselinePath == "" {
		return fmt.Errorf("usage: digest diff <manifest-dir> <source-url> <baseline-path>")
	}

	previous, err := drift.LoadDigests(baselinePath)
	if err != nil {
		return fmt.Errorf("load baseline digests: %w", err)
	}

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

	current := drift.BuildDigests(results)
	changed := drift.DigestsChanged(previous, current)

	if len(changed) == 0 {
		fmt.Println("no digest changes detected")
		return nil
	}
	fmt.Printf("%d service(s) with changed drift digest:\n", len(changed))
	for _, svc := range changed {
		fmt.Printf("  - %s\n", svc)
	}
	return nil
}
