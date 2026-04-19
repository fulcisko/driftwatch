package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/source"
)

func runNotifyGenerate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: notify-generate <manifest-dir> <source-url> <rules-file>")
	}
	manifestDir, sourceURL, rulesFile := args[0], args[1], args[2]

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
			return fmt.Errorf("fetch %s: %w", m.Name, err)
		}
		results = append(results, detector.Compare(m, live))
	}

	rf, err := os.Open(rulesFile)
	if err != nil {
		return fmt.Errorf("open rules: %w", err)
	}
	defer rf.Close()
	var rules []drift.NotifyRule
	if err := json.NewDecoder(rf).Decode(&rules); err != nil {
		return fmt.Errorf("parse rules: %w", err)
	}

	events := drift.GenerateNotifyEvents(results, rules)
	if len(events) == 0 {
		fmt.Println("no notifications to send")
		return nil
	}
	for _, e := range events {
		fmt.Printf("[%s] %s -> %s (%s): %s\n", e.Severity, e.Service, e.Channel, e.Target, e.Message)
	}
	return nil
}

func runNotifyShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: notify-show <events-file>")
	}
	events, err := drift.LoadNotifyEvents(args[0])
	if err != nil {
		return fmt.Errorf("load events: %w", err)
	}
	if len(events) == 0 {
		fmt.Println("no notify events found")
		return nil
	}
	for _, e := range events {
		fmt.Printf("%s  [%s] %s  channel=%s target=%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"), e.Severity, e.Service, e.Channel, e.Target)
	}
	return nil
}
