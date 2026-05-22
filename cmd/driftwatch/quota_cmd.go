package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

// runQuotaSet adds or updates a drift quota rule for a service.
// Usage: driftwatch quota set <quota-file> <service> <max-drift>
func runQuotaSet(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: quota set <quota-file> <service> <max-drift>")
	}
	path, service, maxStr := args[0], args[1], args[2]
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return fmt.Errorf("invalid max-drift value %q: %w", maxStr, err)
	}
	if err := drift.AddQuotaRule(path, service, max); err != nil {
		return fmt.Errorf("failed to set quota: %w", err)
	}
	fmt.Fprintf(os.Stdout, "quota set: service=%s max_drift=%d\n", service, max)
	return nil
}

// runQuotaCheck checks drift results against quota rules.
// Usage: driftwatch quota check <quota-file> <manifest-dir> <source-url>
func runQuotaCheck(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: quota check <quota-file> <manifest-dir> <source-url>")
	}
	quotaFile, manifestDir, sourceURL := args[0], args[1], args[2]

	manifests, err := manifest.LoadDir(manifestDir)
	if err != nil {
		return fmt.Errorf("loading manifests: %w", err)
	}

	fetcher := source.NewFetcher(sourceURL)
	detector := drift.NewDetector()
	var results []drift.CompareResult
	for _, m := range manifests {
		deployed, err := fetcher.Fetch(m.Name)
		if err != nil {
			continue
		}
		res := detector.Compare(m, deployed)
		results = append(results, res)
	}

	violations, err := drift.CheckQuotas(quotaFile, results)
	if err != nil {
		return fmt.Errorf("checking quotas: %w", err)
	}

	if len(violations) == 0 {
		fmt.Fprintln(os.Stdout, "all services within drift quota")
		return nil
	}

	fmt.Fprintln(os.Stdout, "quota violations:")
	for _, v := range violations {
		fmt.Fprintf(os.Stdout, "  service=%-20s allowed=%d actual=%d\n", v.Service, v.Allowed, v.Actual)
	}
	return fmt.Errorf("%d quota violation(s) detected", len(violations))
}
