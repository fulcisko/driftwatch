package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func ownerPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ownership.json")
}

func TestAddAndLoadOwner(t *testing.T) {
	p := ownerPath(t)
	if err := AddOwner(p, "api", "platform", "platform@example.com"); err != nil {
		t.Fatalf("AddOwner: %v", err)
	}
	om, err := LoadOwnership(p)
	if err != nil {
		t.Fatalf("LoadOwnership: %v", err)
	}
	if len(om.Owners) != 1 {
		t.Fatalf("expected 1 owner, got %d", len(om.Owners))
	}
	if om.Owners[0].Team != "platform" {
		t.Errorf("expected team 'platform', got %q", om.Owners[0].Team)
	}
}

func TestAddOwner_UpdatesExisting(t *testing.T) {
	p := ownerPath(t)
	_ = AddOwner(p, "api", "old-team", "")
	_ = AddOwner(p, "api", "new-team", "new@example.com")
	om, _ := LoadOwnership(p)
	if len(om.Owners) != 1 {
		t.Fatalf("expected 1 owner after update, got %d", len(om.Owners))
	}
	if om.Owners[0].Team != "new-team" {
		t.Errorf("expected updated team 'new-team', got %q", om.Owners[0].Team)
	}
}

func TestAddOwner_MissingFields(t *testing.T) {
	p := ownerPath(t)
	if err := AddOwner(p, "", "team", ""); err == nil {
		t.Error("expected error for missing service")
	}
	if err := AddOwner(p, "svc", "", ""); err == nil {
		t.Error("expected error for missing team")
	}
}

func TestRemoveOwner_Success(t *testing.T) {
	p := ownerPath(t)
	_ = AddOwner(p, "api", "platform", "")
	_ = AddOwner(p, "worker", "infra", "")
	if err := RemoveOwner(p, "api"); err != nil {
		t.Fatalf("RemoveOwner: %v", err)
	}
	om, _ := LoadOwnership(p)
	if len(om.Owners) != 1 || om.Owners[0].Service != "worker" {
		t.Errorf("unexpected owners after remove: %+v", om.Owners)
	}
}

func TestRemoveOwner_NotFound(t *testing.T) {
	p := ownerPath(t)
	_ = AddOwner(p, "api", "platform", "")
	if err := RemoveOwner(p, "missing"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestLoadOwnership_NotFound(t *testing.T) {
	om, err := LoadOwnership(filepath.Join(t.TempDir(), "none.json"))
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
	if len(om.Owners) != 0 {
		t.Errorf("expected empty owners, got %+v", om.Owners)
	}
}

func TestLookupOwner(t *testing.T) {
	p := ownerPath(t)
	_ = AddOwner(p, "api", "platform", "platform@example.com")
	om, _ := LoadOwnership(p)
	o, ok := LookupOwner(om, "api")
	if !ok {
		t.Fatal("expected to find owner for 'api'")
	}
	if o.Contact != "platform@example.com" {
		t.Errorf("unexpected contact: %q", o.Contact)
	}
	_, ok = LookupOwner(om, "unknown")
	if ok {
		t.Error("expected no owner for 'unknown'")
	}
}

func TestLoadOwnership_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(p, []byte("not-json"), 0644)
	_, err := LoadOwnership(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
