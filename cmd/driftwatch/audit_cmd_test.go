package main

import (
	"os"
	"path/filepath"
	"testing"
)

func withAuditFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.json")
}

func TestRunAuditAppend_Success(t *testing.T) {
	p := withAuditFile(t)
	err := runAuditAppend([]string{p, "check", "svc-a", "alice", "drift detected"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if len(data) == 0 {
		t.Error("expected non-empty audit file")
	}
}

func TestRunAuditAppend_MissingArgs(t *testing.T) {
	err := runAuditAppend([]string{"/tmp/x.json", "check"})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunAuditShow_Empty(t *testing.T) {
	p := withAuditFile(t)
	err := runAuditShow([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAuditShow_WithEntries(t *testing.T) {
	p := withAuditFile(t)
	_ = runAuditAppend([]string{p, "diff", "svc-b", "bob", ""})
	err := runAuditShow([]string{p, "svc-b", ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAuditShow_MissingArgs(t *testing.T) {
	err := runAuditShow([]string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunAuditShow_FilterByAction(t *testing.T) {
	p := withAuditFile(t)
	_ = runAuditAppend([]string{p, "save-baseline", "svc-a", "carol", ""})
	_ = runAuditAppend([]string{p, "diff", "svc-a", "carol", ""})
	err := runAuditShow([]string{p, "svc-a", "diff"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
