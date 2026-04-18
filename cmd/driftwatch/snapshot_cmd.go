package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

// runSaveSnapshot saves a labeled snapshot of the current drift state.
// Usage: driftwatch snapshot save <label> <manifest-dir> <source-url> <out-file>
func runSaveSnapshot(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: snapshot save <label> <manifest-dir> <source-url> <out-file>")
	}
	label, dir, sourceURL, outFile := args[0], args[1], args[2], args[3]

	manifests, err := manifest.LoadDir(dir)
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

	if err := drift.SaveSnapshot(outFile, label, results); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}
	fmt.Printf("Snapshot '%s' saved to %s (%d services)\n", label, outFile, len(results))
	return nil
}

// runDiffSnapshot compares current drift state against a saved snapshot.
// Usage: driftwatch snapshot diff <snapshot-file> <manifest-dir> <source-url>
func runDiffSnapshot(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: snapshot diff <snapshot-file> <manifest-dir> <source-url>")
	}
	snapFile, dir, sourceURL := args[0], args[1], args[2]

	snap, err := drift.LoadSnapshot(snapFile)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	manifests, err := manifest.LoadDir(dir)
	if err != nil {
		return fmt.Errorf("load manifests: %w", err)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()
	var current []drift.CompareResult
	for _, m := range manifests {
		live, err := fetcher.Fetch(m.Name)
		if err != nil {
			continue
		}
		current = append(current, detector.Compare(m, live))
	}

	changed := drift.DiffSnapshot(snap, current)
	if len(changed) == 0 {
		fmt.Println("No drift changes since snapshot:", snap.Label)
		return nil
	}
	fmt.Printf("Drift changed for %d service(s) since snapshot '%s':\n", len(changed), snap.Label)
	for _, s := range changed {
		fmt.Println(" -", strings.TrimSpace(s))
	}
	return nil
}
