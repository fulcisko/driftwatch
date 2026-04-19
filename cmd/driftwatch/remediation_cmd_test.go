package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func withRemPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "rem.json")
}

func TestRunRemediationAdd_Success(t *testing.T) {
	p := withRemPath(t)
	if err := runRemediationAdd([]string{"svc-a", "replicas", "apply", "scaled up"}, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRemediationAdd_MissingArgs(t *testing.T) {
	p := withRemPath(t)
	if err := runRemediationAdd([]string{"svc-a", "replicas"}, p); err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunRemediationAdd_BadAction(t *testing.T) {
	p := withRemPath(t)
	err := runRemediationAdd([]string{"svc-a", "replicas", "destroy"}, p)
	if err == nil || !strings.Contains(err.Error(), "unknown action") {
		t.Fatalf("expected unknown action error, got %v", err)
	}
}

func TestRunRemediationShow_Empty(t *testing.T) {
	p := withRemPath(t)
	if err := runRemediationShow([]string{"svc-x"}, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRemediationShow_WithEntries(t *testing.T) {
	p := withRemPath(t)
	_ = runRemediationAdd([]string{"svc-a", "cpu", "revert"}, p)
	_ = runRemediationAdd([]string{"svc-a", "mem", "ignore", "low priority"}, p)
	if err := runRemediationShow([]string{"svc-a"}, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRemediationShow_MissingArgs(t *testing.T) {
	p := withRemPath(t)
	if err := runRemediationShow([]string{}, p); err == nil {
		t.Fatal("expected error")
	}
}
