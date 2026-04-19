package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

func withPinFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestRunPinAdd_Success(t *testing.T) {
	p := withPinFile(t)
	if err := runPinAdd([]string{"svc-a", "LOG_LEVEL", "info", "stable"}, p); err != nil {
		t.Fatal(err)
	}
	list, _ := drift.LoadPins(p)
	if len(list.Pins) != 1 {
		t.Errorf("expected 1 pin, got %d", len(list.Pins))
	}
}

func TestRunPinAdd_MissingArgs(t *testing.T) {
	err := runPinAdd([]string{"svc-a", "LOG_LEVEL"}, withPinFile(t))
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRunPinRemove_Success(t *testing.T) {
	p := withPinFile(t)
	_ = drift.AddPin(p, "svc-a", "LOG_LEVEL", "info", "")
	if err := runPinRemove([]string{"svc-a", "LOG_LEVEL"}, p); err != nil {
		t.Fatal(err)
	}
	list, _ := drift.LoadPins(p)
	if len(list.Pins) != 0 {
		t.Error("expected 0 pins after remove")
	}
}

func TestRunPinRemove_MissingArgs(t *testing.T) {
	err := runPinRemove([]string{"svc-a"}, withPinFile(t))
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRunPinShow_Empty(t *testing.T) {
	if err := runPinShow(withPinFile(t)); err != nil {
		t.Fatal(err)
	}
}

func TestRunPinShow_WithEntries(t *testing.T) {
	p := withPinFile(t)
	_ = drift.AddPin(p, "svc-b", "TIMEOUT", "30", "do not change")
	if err := runPinShow(p); err != nil {
		t.Fatal(err)
	}
}
