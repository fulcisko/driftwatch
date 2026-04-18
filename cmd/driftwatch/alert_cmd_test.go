package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

func withAlertFile(t *testing.T, alerts []drift.Alert) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(alerts); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunAlertShow_Empty(t *testing.T) {
	path := withAlertFile(t, []drift.Alert{})
	if err := runAlertShow([]string{path}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunAlertShow_WithEntries(t *testing.T) {
	alerts := []drift.Alert{
		{Service: "api", Key: "db_password", Severity: drift.AlertHigh, Message: "drift detected", Timestamp: time.Now().UTC()},
	}
	path := withAlertFile(t, alerts)
	if err := runAlertShow([]string{path}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunAlertShow_MissingArgs(t *testing.T) {
	if err := runAlertShow([]string{}); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunAlertGenerate_MissingArgs(t *testing.T) {
	if err := runAlertGenerate([]string{}); err == nil {
		t.Error("expected error for missing args")
	}
	if err := runAlertGenerate([]string{"dir", "url"}); err == nil {
		t.Error("expected error for missing output file")
	}
}
