package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func attrPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "attributions.json")
}

func TestAddAndLoadAttribution(t *testing.T) {
	p := attrPath(t)
	if err := AddAttribution(p, "svc-a", "replicas", "alice", "platform", "scale-up"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	store, err := LoadAttributions(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(store.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(store.Entries))
	}
	e := store.Entries[0]
	if e.Service != "svc-a" || e.Key != "replicas" || e.Owner != "alice" || e.Team != "platform" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestAddAttribution_MissingFields(t *testing.T) {
	p := attrPath(t)
	if err := AddAttribution(p, "", "key", "owner", "", ""); err == nil {
		t.Error("expected error for missing service")
	}
	if err := AddAttribution(p, "svc", "", "owner", "", ""); err == nil {
		t.Error("expected error for missing key")
	}
	if err := AddAttribution(p, "svc", "key", "", "", ""); err == nil {
		t.Error("expected error for missing owner")
	}
}

func TestLoadAttributions_NotFound(t *testing.T) {
	store, err := LoadAttributions("/nonexistent/attributions.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestFilterAttributions_ByService(t *testing.T) {
	p := attrPath(t)
	_ = AddAttribution(p, "svc-a", "replicas", "alice", "platform", "")
	_ = AddAttribution(p, "svc-b", "image", "bob", "infra", "")
	_ = AddAttribution(p, "svc-a", "memory", "alice", "platform", "")
	store, _ := LoadAttributions(p)
	results := FilterAttributions(store, "svc-a")
	if len(results) != 2 {
		t.Errorf("expected 2 attributions for svc-a, got %d", len(results))
	}
}

func TestFilterAttributions_EmptyServiceReturnsAll(t *testing.T) {
	p := attrPath(t)
	_ = AddAttribution(p, "svc-a", "replicas", "alice", "platform", "")
	_ = AddAttribution(p, "svc-b", "image", "bob", "infra", "")
	store, _ := LoadAttributions(p)
	results := FilterAttributions(store, "")
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestAddAttribution_Timestamp(t *testing.T) {
	p := attrPath(t)
	_ = AddAttribution(p, "svc-a", "replicas", "alice", "", "")
	store, _ := LoadAttributions(p)
	if store.Entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAttributions_Persistence(t *testing.T) {
	p := attrPath(t)
	_ = AddAttribution(p, "svc-a", "replicas", "alice", "platform", "scale-up")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty file")
	}
}
