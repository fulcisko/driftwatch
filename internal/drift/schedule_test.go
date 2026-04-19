package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func schedPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "schedule.json")
}

func TestLoadSchedule_NotFound(t *testing.T) {
	s, err := LoadSchedule("/nonexistent/schedule.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(s.Entries) != 0 {
		t.Errorf("expected empty schedule")
	}
}

func TestSaveAndLoadSchedule(t *testing.T) {
	path := schedPath(t)
	s := &Schedule{
		Entries: []ScheduleEntry{
			{Service: "api", Interval: time.Hour, LastRun: time.Time{}, Enabled: true},
		},
	}
	if err := SaveSchedule(path, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadSchedule(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != 1 || loaded.Entries[0].Service != "api" {
		t.Errorf("unexpected entries: %+v", loaded.Entries)
	}
}

func TestUpsertSchedule_AddsAndUpdates(t *testing.T) {
	s := &Schedule{}
	UpsertSchedule(s, ScheduleEntry{Service: "svc", Interval: time.Minute, Enabled: true})
	if len(s.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	UpsertSchedule(s, ScheduleEntry{Service: "svc", Interval: time.Hour, Enabled: false})
	if len(s.Entries) != 1 {
		t.Errorf("expected upsert, got %d entries", len(s.Entries))
	}
	if s.Entries[0].Interval != time.Hour {
		t.Errorf("interval not updated")
	}
}

func TestDueServices(t *testing.T) {
	now := time.Now()
	s := &Schedule{
		Entries: []ScheduleEntry{
			{Service: "due-svc", Interval: time.Minute, LastRun: now.Add(-2 * time.Minute), Enabled: true},
			{Service: "not-due", Interval: time.Hour, LastRun: now.Add(-time.Minute), Enabled: true},
			{Service: "disabled", Interval: time.Minute, LastRun: now.Add(-2 * time.Minute), Enabled: false},
		},
	}
	due := DueServices(s, now)
	if len(due) != 1 || due[0] != "due-svc" {
		t.Errorf("unexpected due services: %v", due)
	}
}

func TestLoadSchedule_InvalidJSON(t *testing.T) {
	path := schedPath(t)
	os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadSchedule(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
