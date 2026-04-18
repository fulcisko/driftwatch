package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddAndLoadTag(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tags.json")

	if err := AddTag(p, "team-a", "svc1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := AddTag(p, "team-a", "svc2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store, err := LoadTags(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(store.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(store.Tags))
	}
	if len(store.Tags[0].Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(store.Tags[0].Services))
	}
}

func TestAddTag_Duplicate(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tags.json")

	_ = AddTag(p, "env-prod", "svc1")
	_ = AddTag(p, "env-prod", "svc1")

	store, _ := LoadTags(p)
	if len(store.Tags[0].Services) != 1 {
		t.Errorf("expected 1 service, got %d", len(store.Tags[0].Services))
	}
}

func TestRemoveTag_Success(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tags.json")

	_ = AddTag(p, "team-b", "svc1")
	_ = AddTag(p, "team-b", "svc2")
	if err := RemoveTag(p, "team-b", "svc1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	store, _ := LoadTags(p)
	svcs := FilterByTag(store, "team-b")
	if len(svcs) != 1 || svcs[0] != "svc2" {
		t.Errorf("unexpected services: %v", svcs)
	}
}

func TestRemoveTag_NotFound(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tags.json")
	_ = AddTag(p, "x", "svc1")
	if err := RemoveTag(p, "missing", "svc1"); err == nil {
		t.Error("expected error for missing tag")
	}
}

func TestLoadTags_NotFound(t *testing.T) {
	store, err := LoadTags("/nonexistent/tags.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store.Tags) != 0 {
		t.Error("expected empty store")
	}
}

func TestFilterByTag_Missing(t *testing.T) {
	store := TagStore{}
	result := FilterByTag(store, "nope")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func init() { _ = os.Getenv("") }
