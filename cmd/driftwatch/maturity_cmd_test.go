package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func writeMaturityManifest(t *testing.T, dir, name string) {
	t.Helper()
	content := fmt.Sprintf("name: %s\nconfig:\n  log_level: info\n", name)
	err := os.WriteFile(filepath.Join(dir, name+".yaml"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("write manifest: %v", err)
	}
}

func TestRunMaturityAssess_MissingArgs(t *testing.T) {
	err := runMaturityAssess([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunMaturityAssess_EmptyManifestDir(t *testing.T) {
	dir := t.TempDir()
	err := runMaturityAssess([]string{dir, "http://localhost:9999"})
	if err == nil {
		t.Fatal("expected error for empty manifest dir")
	}
}

func TestRunMaturityShow_MissingArgs(t *testing.T) {
	err := runMaturityShow([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunMaturityShow_EmptyReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "maturity.json")
	if err := drift.SaveMaturityReport(path, []drift.MaturityEntry{}); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := runMaturityShow([]string{path}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunMaturityShow_WithEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "maturity.json")
	entries := []drift.MaturityEntry{
		{Service: "svc-a", Level: drift.MaturityMature, DriftScore: 0},
		{Service: "svc-b", Level: drift.MaturityUnstable, DriftScore: 25},
	}
	if err := drift.SaveMaturityReport(path, entries); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := runMaturityShow([]string{path}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunMaturityShow_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "maturity.json")
	entries := []drift.MaturityEntry{
		{Service: "svc-json", Level: drift.MaturityStable, DriftScore: 3},
	}
	if err := drift.SaveMaturityReport(path, entries); err != nil {
		t.Fatalf("save: %v", err)
	}

	// Verify saved file is valid JSON
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var report drift.MaturityReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(report.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(report.Entries))
	}
}
