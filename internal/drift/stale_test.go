package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeStaleHistory(t *testing.T, dir string, service string, age time.Duration) string {
	t.Helper()
	path := filepath.Join(dir, "history.json")
	record := HistoryRecord{
		Timestamp: time.Now().UTC().Add(-age),
		Results: []CompareResult{
			{Service: service, Diffs: nil},
		},
	}
	if err := AppendHistory(path, record.Results); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}
	return path
}

func TestFindStaleServices_DetectsStale(t *testing.T) {
	dir := t.TempDir()
	path := makeStaleHistory(t, dir, "auth-service", 50*time.Hour)

	entries, err := FindStaleServices(path, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 stale entry, got %d", len(entries))
	}
	if entries[0].Service != "auth-service" {
		t.Errorf("expected auth-service, got %s", entries[0].Service)
	}
}

func TestFindStaleServices_RecentNotStale(t *testing.T) {
	dir := t.TempDir()
	path := makeStaleHistory(t, dir, "billing-service", 1*time.Hour)

	entries, err := FindStaleServices(path, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 stale entries, got %d", len(entries))
	}
}

func TestFindStaleServices_NotFound(t *testing.T) {
	entries, err := FindStaleServices("/nonexistent/history.json", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty result for missing file")
	}
}

func TestSaveAndLoadStaleReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stale.json")

	input := []StaleEntry{
		{Service: "svc-a", LastSeen: time.Now().UTC().Add(-48 * time.Hour), StaleSince: 48 * time.Hour},
	}
	if err := SaveStaleReport(path, input); err != nil {
		t.Fatalf("SaveStaleReport: %v", err)
	}
	loaded, err := LoadStaleReport != nil {
		t.Fatalf("LoadStaleReport: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Service != "svc-a" {
		t.Errorf("unexpected loaded entries: %+v", loaded)
	}
}

func TestLoadStaleReport_NotFound(t *testing.T) {
	entries, err := LoadStaleReport("/no/such/file.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil for missing file")
	}
}

func TestFormatStaleReport_Empty(t *testing.T) {
	out := FormatStaleReport(nil)
	if out != "No stale services detected.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatStaleReport_WithEntries(t *testing.T) {
	entries := []StaleEntry{
		{Service: "order-service", LastSeen: time.Now().Add(-72 * time.Hour), StaleSince: 72 * time.Hour},
	}
	out := FormatStaleReport(entries)
	if !contains(out, "order-service") {
		t.Errorf("expected service name in output, got: %s", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func init() {
	_ = os.Getenv // suppress unused import
}
