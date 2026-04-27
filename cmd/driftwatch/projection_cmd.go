package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/manifest"
	"github.com/driftwatch/internal/source"
)

// runProjection runs a field-projected view of drift results.
// Usage: driftwatch projection <manifest-dir> <source-url> <field1>[=alias1] [field2[=alias2] ...]
func runProjection(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: projection <manifest-dir> <source-url> <field>[=alias] ...")
	}
	manifestDir := args[0]
	sourceURL := args[1]
	fieldArgs := args[2:]

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
		res := detector.Compare(m, live)
		results = append(results, res)
	}

	fields := parseProjectionFields(fieldArgs)
	serviceFilter := envOr("DRIFTWATCH_SERVICE", "")

	opts := drift.ProjectionOptions{
		Fields:  fields,
		Service: serviceFilter,
	}

	rows := drift.ApplyProjection(results, opts)
	fmt.Print(drift.FormatProjection(rows, fields))
	return nil
}

func parseProjectionFields(args []string) []drift.ProjectionField {
	fields := make([]drift.ProjectionField, 0, len(args))
	for _, a := range args {
		parts := strings.SplitN(a, "=", 2)
		f := drift.ProjectionField{Key: parts[0]}
		if len(parts) == 2 {
			f.Alias = parts[1]
		}
		fields = append(fields, f)
	}
	return fields
}
