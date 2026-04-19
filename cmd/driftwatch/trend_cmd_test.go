package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

func writeTrendFile(t *testing.T, dir string, points []drift.TrendPoint) string {
	t.Helper()
	path := filepath.Join(dir, "trend.json")
	report := drift.TrendReport{Points: points}
	data, _ := json.MarshalIndent(report, "", "  ")
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestRunTrendShow_Empty(t *testing.T) {
	dir := t.TempDir()
	path := writeTrendFile(t, dir, nil)
	if err := runTrendShow([]string{path}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTrendShow_WithEntries(t *testing.T) {
	dir := t.TempDir()
	points := []drift.TrendPoint{
		{Timestamp: time.Now().UTC(), Service: "api", DriftCount: 2, MaxSeverity: "high"},
	}
	path := writeTrendFile(t, dir, points)
	if err := runTrendShow([]string{path, "api"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunTrendShow_MissingArgs(t *testing.T) {
	if err := runTrendShow([]string{}); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunTrendAppend_MissingArgs(t *testing.T) {
	if err := runTrendAppend([]string{"only-one"}); err == nil {
		t.Error("expected error for missing args")
	}
}
