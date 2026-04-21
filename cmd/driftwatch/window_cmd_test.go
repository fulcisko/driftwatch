package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
)

func writeWindowHistory(t *testing.T, entries []drift.HistoryEntry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "history.json")
	data, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("marshal history: %v", err)
	}
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatalf("write history: %v", err)
	}
	return p
}

func TestRunWindowCheck_MissingArgs(t *testing.T) {
	err := runWindowCheck([]string{"only-one"}, "")
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunWindowCheck_BadHours(t *testing.T) {
	err := runWindowCheck([]string{"dir", "http://x", "notanumber"}, "")
	if err == nil {
		t.Fatal("expected error for bad hours value")
	}
}

func TestRunWindowCheck_ZeroHours(t *testing.T) {
	err := runWindowCheck([]string{"dir", "http://x", "0"}, "")
	if err == nil {
		t.Fatal("expected error for zero hours")
	}
}

func TestRunWindowCheck_EmptyManifestDir(t *testing.T) {
	dir := t.TempDir()
	histPath := writeWindowHistory(t, []drift.HistoryEntry{
		{Service: "api", RecordedAt: time.Now().UTC()},
	})
	err := runWindowCheck([]string{dir, "http://localhost:9999", "24"}, histPath)
	if err != nil {
		t.Fatalf("unexpected error with empty manifest dir: %v", err)
	}
}

func TestRunWindowCheck_MissingHistoryFile(t *testing.T) {
	dir := t.TempDir()
	err := runWindowCheck([]string{dir, "http://localhost:9999", "12"}, filepath.Join(dir, "no-history.json"))
	if err != nil {
		t.Fatalf("should not error when history file missing: %v", err)
	}
}
