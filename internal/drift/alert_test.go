package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeAlertResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []Diff{
				{Key: "db_password", Expected: "old", Actual: "new"},
				{Key: "log_level", Expected: "info", Actual: "debug"},
			},
		},
		{
			Service: "worker",
			Diffs:   []Diff{},
		},
	}
}

func TestGenerateAlerts_HighOnly(t *testing.T) {
	results := makeAlertResults()
	alerts := GenerateAlerts(results, AlertConfig{MinSeverity: AlertHigh})
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Key != "db_password" {
		t.Errorf("expected db_password, got %s", alerts[0].Key)
	}
}

func TestGenerateAlerts_LowThreshold(t *testing.T) {
	results := makeAlertResults()
	alerts := GenerateAlerts(results, AlertConfig{MinSeverity: AlertLow})
	if len(alerts) < 1 {
		t.Fatal("expected at least 1 alert")
	}
}

func TestGenerateAlerts_NoCleanServices(t *testing.T) {
	results := makeAlertResults()
	alerts := GenerateAlerts(results, AlertConfig{MinSeverity: AlertLow})
	for _, a := range alerts {
		if a.Service == "worker" {
			t.Error("worker has no diffs, should not appear in alerts")
		}
	}
}

func TestSaveAndLoadAlerts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")
	alerts := GenerateAlerts(makeAlertResults(), AlertConfig{MinSeverity: AlertLow})
	if err := SaveAlerts(path, alerts); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadAlerts(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != len(alerts) {
		t.Errorf("expected %d alerts, got %d", len(alerts), len(loaded))
	}
}

func TestLoadAlerts_NotFound(t *testing.T) {
	alerts, err := LoadAlerts("/nonexistent/alerts.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alerts != nil {
		t.Error("expected nil for missing file")
	}
	_ = os.Remove("/nonexistent/alerts.json")
}
