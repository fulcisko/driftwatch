package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

// runPromotionAssess assesses whether services are safe to promote between stages.
// Usage: driftwatch promotion assess <manifest-dir> <source-url> <from-stage> <to-stage> [report-path]
func runPromotionAssess(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: promotion assess <manifest-dir> <source-url> <from-stage> <to-stage> [report-path]")
	}
	manifestDir := args[0]
	sourceURL := args[1]
	fromStage := drift.PromotionStage(args[2])
	toStage := drift.PromotionStage(args[3])

	reportPath := ""
	if len(args) >= 5 {
		reportPath = args[4]
	}

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("loading manifests: %w", err)
	}
	if len(manifests) == 0 {
		return fmt.Errorf("no manifests found in %s", manifestDir)
	}

	fetcher := source.NewFetcher(nil)
	detector := drift.NewDetector()

	var results []drift.CompareResult
	for _, m := range manifests {
		deployed, err := fetcher.Fetch(sourceURL + "/" + m.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: could not fetch %s: %v\n", m.Name, err)
			continue
		}
		res := detector.Compare(m, deployed)
		results = append(results, res)
	}

	entries := drift.AssessPromotion(results, fromStage, toStage)
	fmt.Print(drift.FormatPromotion(entries))

	if reportPath != "" {
		if err := drift.SavePromotionReport(reportPath, entries); err != nil {
			return fmt.Errorf("saving report: %w", err)
		}
		fmt.Fprintf(os.Stderr, "promotion report saved to %s\n", reportPath)
	}
	return nil
}

// runPromotionShow displays a previously saved promotion report.
// Usage: driftwatch promotion show <report-path>
func runPromotionShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: promotion show <report-path>")
	}
	entries, err := drift.LoadPromotionReport(args[0])
	if err != nil {
		return fmt.Errorf("loading report: %w", err)
	}
	fmt.Print(drift.FormatPromotion(entries))
	return nil
}
