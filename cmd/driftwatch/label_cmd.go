package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

func runLabelAdd(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: label add <service> <key> <value>")
	}
	path := envOr("DRIFTWATCH_LABEL_PATH", "labels.json")
	return drift.AddLabel(path, args[0], args[1], args[2])
}

func runLabelRemove(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: label remove <service> <key>")
	}
	path := envOr("DRIFTWATCH_LABEL_PATH", "labels.json")
	return drift.RemoveLabel(path, args[0], args[1])
}

func runLabelShow(args []string) error {
	path := envOr("DRIFTWATCH_LABEL_PATH", "labels.json")
	labels, err := drift.LoadLabels(path)
	if err != nil {
		return err
	}
	if len(labels) == 0 {
		fmt.Println("no labels defined")
		return nil
	}
	for svc, kv := range labels {
		for k, v := range kv {
			fmt.Fprintf(os.Stdout, "%-30s %s=%s\n", svc, k, v)
		}
	}
	return nil
}

func runLabelFilter(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: label filter <key> <value>")
	}
	labelPath := envOr("DRIFTWATCH_LABEL_PATH", "labels.json")
	labels, err := drift.LoadLabels(labelPath)
	if err != nil {
		return err
	}
	sourceURL := envOr("DRIFTWATCH_SOURCE_URL", "")
	if sourceURL == "" {
		return fmt.Errorf("DRIFTWATCH_SOURCE_URL is required")
	}
	// Placeholder: in full impl, load results then filter
	for svc, kv := range labels {
		if kv[args[0]] == args[1] {
			fmt.Println(svc)
		}
	}
	return nil
}
