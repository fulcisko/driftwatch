package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func annotationPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "annotations.json")
}

func TestAddAndLoadAnnotation(t *testing.T) {
	p := annotationPath(t)
	if err := AddAnnotation(p, "svc-a", "replicas", "intentional change", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	anns, err := LoadAnnotations(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Service != "svc-a" || anns[0].Key != "replicas" {
		t.Errorf("unexpected annotation: %+v", anns[0])
	}
}

func TestAddAnnotation_MissingFields(t *testing.T) {
	p := annotationPath(t)
	if err := AddAnnotation(p, "", "key", "note", "author"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestLoadAnnotations_NotFound(t *testing.T) {
	anns, err := LoadAnnotations(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(anns) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestFilterAnnotations(t *testing.T) {
	p := annotationPath(t)
	_ = AddAnnotation(p, "svc-a", "replicas", "note1", "alice")
	_ = AddAnnotation(p, "svc-b", "image", "note2", "bob")
	_ = AddAnnotation(p, "svc-a", "image", "note3", "carol")

	all, _ := LoadAnnotations(p)

	filtered := FilterAnnotations(all, "svc-a", "")
	if len(filtered) != 2 {
		t.Errorf("expected 2, got %d", len(filtered))
	}

	filtered2 := FilterAnnotations(all, "svc-a", "image")
	if len(filtered2) != 1 {
		t.Errorf("expected 1, got %d", len(filtered2))
	}
}

func TestAddAnnotation_MultipleEntries(t *testing.T) {
	p := annotationPath(t)
	for i := 0; i < 3; i++ {
		_ = AddAnnotation(p, "svc", "key", "note", "author")
	}
	anns, _ := LoadAnnotations(p)
	if len(anns) != 3 {
		t.Errorf("expected 3 annotations, got %d", len(anns))
	}
	os.Remove(p)
}
