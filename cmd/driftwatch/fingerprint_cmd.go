package main

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/manifest"
	"github.com/example/driftwatch/internal/source"
)

// runFingerprintBuild loads manifests, fetches live config, computes drift,
// builds a fingerprint store, and saves it to the given path.
func runFingerprintBuild(manifestDir, sourceURL, outPath string) error {
	if manifestDir == "" || sourceURL == "" || outPath == "" {
		return fmt.Errorf("usage: fingerprint build <manifest-dir> <source-url> <output-path>")
	}

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("load manifests: %w", err)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()

	var results []drift.CompareResult
	for _, m := range manifests {
		body, err := fetcher.Fetch(m.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: fetch %s: %v\n", m.Name, err)
			continue
		}
		r := detector.Compare(m, body)
		results = append(results, r)
	}

	store := drift.BuildFingerprintStore(results)
	if err := drift.SaveFingerprintStore(outPath, store); err != nil {
		return fmt.Errorf("save fingerprint store: %w", err)
	}

	fmt.Printf("fingerprinted %d drifted service(s) → %s\n", len(store), outPath)
	return nil
}

// runFingerprintShow loads and prints a fingerprint store from disk.
func runFingerprintShow(storePath string) error {
	if storePath == "" {
		return fmt.Errorf("usage: fingerprint show <store-path>")
	}
	store, err := drift.LoadFingerprintStore(storePath)
	if err != nil {
		return fmt.Errorf("load fingerprint store: %w", err)
	}
	fmt.Print(drift.FormatFingerprintStore(store))
	return nil
}

// runFingerprintDiff compares an old fingerprint store against a newly built one.
func runFingerprintDiff(oldPath, manifestDir, sourceURL string) error {
	if oldPath == "" || manifestDir == "" || sourceURL == "" {
		return fmt.Errorf("usage: fingerprint diff <old-store> <manifest-dir> <source-url>")
	}

	old, err := drift.LoadFingerprintStore(oldPath)
	if err != nil {
		return fmt.Errorf("load old fingerprint store: %w", err)
	}

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("load manifests: %w", err)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()

	var results []drift.CompareResult
	for _, m := range manifests {
		body, err := fetcher.Fetch(m.Name)
		if err != nil {
			continue
		}
		results = append(results, detector.Compare(m, body))
	}

	current := drift.BuildFingerprintStore(results)
	changed := drift.DiffFingerprintStore(old, current)

	if len(changed) == 0 {
		fmt.Println("no fingerprint changes detected")
		return nil
	}
	fmt.Printf("%d service(s) with changed drift fingerprint:\n", len(changed))
	for _, svc := range changed {
		fmt.Printf("  %s\n", svc)
	}
	return nil
}
