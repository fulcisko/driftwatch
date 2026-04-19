package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeNotifyResults() []CompareResult {
	return []CompareResult{
		{
			Service: "auth",
			Diffs: []DiffEntry{
				{Key: "secret_key", Expected: "x", Actual: "y"},
			},
		},
		{
			Service: "cache",
			Diffs:   []DiffEntry{},
		},
		{
			Service: "worker",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
	}
}

func TestGenerateNotifyEvents_HighOnly(t *testing.T) {
	rules := []NotifyRule{
		{Channel: ChannelSlack, Target: "#alerts", MinSeverity: "high"},
	}
	events := GenerateNotifyEvents(makeNotifyResults(), rules)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Service != "auth" {
		t.Errorf("expected auth, got %s", events[0].Service)
	}
}

func TestGenerateNotifyEvents_AllDrifted(t *testing.T) {
	rules := []NotifyRule{
		{Channel: ChannelEmail, Target: "ops@example.com", MinSeverity: "low"},
	}
	events := GenerateNotifyEvents(makeNotifyResults(), rules)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestGenerateNotifyEvents_ServiceFilter(t *testing.T) {
	rules := []NotifyRule{
		{Channel: ChannelWebhook, Target: "http://hook", MinSeverity: "low", Services: []string{"worker"}},
	}
	events := GenerateNotifyEvents(makeNotifyResults(), rules)
	if len(events) != 1 || events[0].Service != "worker" {
		t.Errorf("expected worker only, got %+v", events)
	}
}

func TestSaveAndLoadNotifyEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notify.json")
	rules := []NotifyRule{{Channel: ChannelSlack, Target: "#ops", MinSeverity: "low"}}
	events := GenerateNotifyEvents(makeNotifyResults(), rules)
	if err := SaveNotifyEvents(path, events); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadNotifyEvents(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != len(events) {
		t.Errorf("expected %d, got %d", len(events), len(loaded))
	}
}

func TestLoadNotifyEvents_NotFound(t *testing.T) {
	events, err := LoadNotifyEvents(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if events != nil {
		t.Error("expected nil")
	}
}

func TestGenerateNotifyEvents_NoRules(t *testing.T) {
	events := GenerateNotifyEvents(makeNotifyResults(), nil)
	if len(events) != 0 {
		t.Errorf("expected 0 events")
	}
	_ = os.Getenv("CI")
}
