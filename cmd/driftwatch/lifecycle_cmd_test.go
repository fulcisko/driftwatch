package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

func withLifecyclePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "lifecycle.json")
}

func TestRunLifecycleSet_Success(t *testing.T) {
	p := withLifecyclePath(t)
	err := runLifecycleSet([]string{"lifecycle", p, "svc-x", "active", "running fine"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	store, _ := drift.LoadLifecycle(p)
	if len(store.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(store.Entries))
	}
	if store.Entries[0].Stage != drift.StageActive {
		t.Errorf("expected active stage")
	}
}

func TestRunLifecycleSet_MissingArgs(t *testing.T) {
	err := runLifecycleSet([]string{"lifecycle", "/tmp/x.json", "svc-x"})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunLifecycleShow_Empty(t *testing.T) {
	p := withLifecyclePath(t)
	err := runLifecycleShow([]string{"lifecycle", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunLifecycleShow_WithEntries(t *testing.T) {
	p := withLifecyclePath(t)
	_ = drift.SetLifecycle(p, "svc-a", drift.StageWatched, "monitoring")
	_ = drift.SetLifecycle(p, "svc-b", drift.StageRetired, "")
	err := runLifecycleShow([]string{"lifecycle", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunLifecycleShow_FilteredByStage(t *testing.T) {
	p := withLifecyclePath(t)
	_ = drift.SetLifecycle(p, "svc-a", drift.StageActive, "")
	_ = drift.SetLifecycle(p, "svc-b", drift.StageRetired, "")
	err := runLifecycleShow([]string{"lifecycle", p, "retired"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunLifecycleShow_MissingArgs(t *testing.T) {
	err := runLifecycleShow([]string{"lifecycle"})
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestRunLifecycleShow_JSONStructure(t *testing.T) {
	p := withLifecyclePath(t)
	_ = drift.SetLifecycle(p, "svc-z", drift.StageDeprecated, "will be removed")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var store struct {
		Entries []struct {
			Service   string    `json:"service"`
			Stage     string    `json:"stage"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(data, &store); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !strings.Contains(string(data), "deprecated") {
		t.Errorf("expected 'deprecated' in output")
	}
}
