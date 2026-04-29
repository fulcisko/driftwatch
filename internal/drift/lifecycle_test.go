package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func lifecyclePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "lifecycle.json")
}

func TestSetAndLoadLifecycle(t *testing.T) {
	p := lifecyclePath(t)
	if err := SetLifecycle(p, "svc-a", StageActive, "initial"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	store, err := LoadLifecycle(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(store.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(store.Entries))
	}
	if store.Entries[0].Stage != StageActive {
		t.Errorf("expected active, got %s", store.Entries[0].Stage)
	}
}

func TestSetLifecycle_UpdatesExisting(t *testing.T) {
	p := lifecyclePath(t)
	_ = SetLifecycle(p, "svc-a", StageActive, "")
	_ = SetLifecycle(p, "svc-a", StageDeprecated, "old service")
	store, _ := LoadLifecycle(p)
	if len(store.Entries) != 1 {
		t.Fatalf("expected 1 entry after update, got %d", len(store.Entries))
	}
	if store.Entries[0].Stage != StageDeprecated {
		t.Errorf("expected deprecated, got %s", store.Entries[0].Stage)
	}
	if store.Entries[0].Note != "old service" {
		t.Errorf("expected note to be updated")
	}
}

func TestSetLifecycle_MissingService(t *testing.T) {
	p := lifecyclePath(t)
	if err := SetLifecycle(p, "", StageActive, ""); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestSetLifecycle_MissingStage(t *testing.T) {
	p := lifecyclePath(t)
	if err := SetLifecycle(p, "svc-a", "", ""); err == nil {
		t.Error("expected error for missing stage")
	}
}

func TestLoadLifecycle_NotFound(t *testing.T) {
	store, err := LoadLifecycle("/nonexistent/lifecycle.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestFilterByStage(t *testing.T) {
	p := lifecyclePath(t)
	_ = SetLifecycle(p, "svc-a", StageActive, "")
	_ = SetLifecycle(p, "svc-b", StageRetired, "")
	_ = SetLifecycle(p, "svc-c", StageActive, "")
	store, _ := LoadLifecycle(p)
	active := FilterByStage(store, StageActive)
	if len(active) != 2 {
		t.Errorf("expected 2 active entries, got %d", len(active))
	}
	retired := FilterByStage(store, StageRetired)
	if len(retired) != 1 {
		t.Errorf("expected 1 retired entry, got %d", len(retired))
	}
}

func TestFilterByStage_Empty(t *testing.T) {
	_ = os.Setenv("TEST_LIFECYCLE", "1")
	store := LifecycleStore{}
	out := FilterByStage(store, StageWatched)
	if len(out) != 0 {
		t.Errorf("expected empty result for empty store")
	}
}
