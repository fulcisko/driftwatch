package main

import (
	"os"
	"path/filepath"
	"testing"
)

func withAnnotationPath(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "annotations.json")
	t.Setenv("DRIFTWATCH_ANNOTATIONS", p)
	return p
}

func TestRunAnnotationAdd_Success(t *testing.T) {
	withAnnotationPath(t)
	t.Setenv("DRIFTWATCH_AUTHOR", "tester")
	err := runAnnotationAdd([]string{"svc-a", "replicas", "approved scale-down"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAnnotationAdd_MissingArgs(t *testing.T) {
	withAnnotationPath(t)
	err := runAnnotationAdd([]string{"svc-a", "key"})
	if err == nil {
		t.Error("expected error for missing note")
	}
}

func TestRunAnnotationShow_Empty(t *testing.T) {
	withAnnotationPath(t)
	err := runAnnotationShow(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAnnotationShow_WithEntries(t *testing.T) {
	p := withAnnotationPath(t)
	_ = os.WriteFile(p, []byte(`{"annotations":[{"service":"svc-a","key":"replicas","note":"ok","author":"alice","created_at":"2024-01-01T00:00:00Z"}]}`), 0644)
	err := runAnnotationShow([]string{"svc-a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAnnotationShow_FilterByKey(t *testing.T) {
	p := withAnnotationPath(t)
	_ = os.WriteFile(p, []byte(`{"annotations":[
		{"service":"svc-a","key":"replicas","note":"n1","author":"a","created_at":"2024-01-01T00:00:00Z"},
		{"service":"svc-a","key":"image","note":"n2","author":"b","created_at":"2024-01-01T00:00:00Z"}
	]}`), 0644)
	err := runAnnotationShow([]string{"svc-a", "image"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
