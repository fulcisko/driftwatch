package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

func runPolicyCheck(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: driftwatch policy <manifest-dir> <source-url> <policy-file>")
	}
	manifestDir, sourceURL, policyPath := args[0], args[1], args[2]

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

	policy, err := drift.LoadPolicy(policyPath)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}

	violations := drift.ApplyPolicy(results, policy)
	if len(violations) == 0 {
		fmt.Println("policy check passed: no violations")
		return nil
	}

	fmt.Printf("policy %q: %d violation(s)\n", policy.Name, len(violations))
	for _, v := range violations {
		fmt.Printf("  [%s] %s — %s\n", v.Rule.Severity, v.Service, v.Message)
	}
	return fmt.Errorf("policy violations found")
}
