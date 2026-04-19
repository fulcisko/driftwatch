package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func auditPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.json")
}

func TestAppendAndLoadAuditEvent(t *testing.T) {
	p := auditPath(t)
	if err := AppendAuditEvent(p, "save-baseline", "svc-a", "alice", "saved"); err != nil {
		t.Fatalf("append: %v", err)
	}
	events, err := LoadAuditLog(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Service != "svc-a" || events[0].User != "alice" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestLoadAuditLog_NotFound(t *testing.T) {
	events, err := LoadAuditLog("/tmp/driftwatch-no-such-audit.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected empty slice")
	}
	os.Remove("/tmp/driftwatch-no-such-audit.json")
}

func TestFilterAuditLog_ByService(t *testing.T) {
	events := []AuditEvent{
		{Action: "check", Service: "svc-a", User: "bob"},
		{Action: "check", Service: "svc-b", User: "bob"},
	}
	out := FilterAuditLog(events, "svc-a", "")
	if len(out) != 1 || out[0].Service != "svc-a" {
		t.Errorf("unexpected filter result: %+v", out)
	}
}

func TestFilterAuditLog_ByAction(t *testing.T) {
	events := []AuditEvent{
		{Action: "save-baseline", Service: "svc-a", User: "alice"},
		{Action: "diff", Service: "svc-a", User: "alice"},
	}
	out := FilterAuditLog(events, "", "diff")
	if len(out) != 1 || out[0].Action != "diff" {
		t.Errorf("unexpected filter result: %+v", out)
	}
}

func TestFormatAuditLog_Empty(t *testing.T) {
	out := FormatAuditLog(nil)
	if out != "no audit events found\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatAuditLog_WithEntries(t *testing.T) {
	events := []AuditEvent{
		{Action: "check", Service: "svc-a", User: "carol", Detail: "ok"},
	}
	out := FormatAuditLog(events)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
