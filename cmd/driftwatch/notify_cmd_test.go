package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

func writeNotifyRules(t *testing.T, dir string, rules []drift.NotifyRule) string {
	t.Helper()
	path := filepath.Join(dir, "rules.json")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(rules); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunNotifyGenerate_MissingArgs(t *testing.T) {
	if err := runNotifyGenerate([]string{"only-one"}); err == nil {
		t.Error("expected error")
	}
}

func TestRunNotifyGenerate_BadRulesFile(t *testing.T) {
	dir := t.TempDir()
	err := runNotifyGenerate([]string{dir, "http://localhost", filepath.Join(dir, "missing.json")})
	if err == nil {
		t.Error("expected error for missing rules file")
	}
}

func TestRunNotifyShow_MissingArgs(t *testing.T) {
	if err := runNotifyShow([]string{}); err == nil {
		t.Error("expected error")
	}
}

func TestRunNotifyShow_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.json")
	events := []drift.NotifyEvent{}
	f, _ := os.Create(path)
	json.NewEncoder(f).Encode(events)
	f.Close()
	if err := runNotifyShow([]string{path}); err != nil {
		t.Fatal(err)
	}
}

func TestRunNotifyShow_WithEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.json")
	rules := []drift.NotifyRule{
		{Channel: drift.ChannelSlack, Target: "#ops", MinSeverity: "low"},
	}
	results := []drift.CompareResult{
		{Service: "api", Diffs: []drift.DiffEntry{{Key: "timeout", Expected: "10", Actual: "20"}}},
	}
	events := drift.GenerateNotifyEvents(results, rules)
	drift.SaveNotifyEvents(path, events)
	if err := runNotifyShow([]string{path}); err != nil {
		t.Fatal(err)
	}
}
