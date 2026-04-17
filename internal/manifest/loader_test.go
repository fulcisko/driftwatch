package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/manifest"
)

func writeTemp(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return p
}

func TestLoadFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := writeTemp(t, dir, "svc.yaml", `
name: auth-service
version: "1.2.3"
namespace: production
image: auth:latest
replicas: 3
env:
  LOG_LEVEL: info
  PORT: "8080"
`)

	m, err := manifest.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "auth-service" {
		t.Errorf("expected name %q, got %q", "auth-service", m.Name)
	}
	if m.Replicas != 3 {
		t.Errorf("expected replicas 3, got %d", m.Replicas)
	}
	if m.Env["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", m.Env["LOG_LEVEL"])
	}
}

func TestLoadFile_MissingName(t *testing.T) {
	dir := t.TempDir()
	path := writeTemp(t, dir, "bad.yaml", `version: "1.0.0"`)

	_, err := manifest.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := manifest.LoadFile("/nonexistent/path/manifest.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadDir(t *testing.T) {
	dir := t.TempDir()
	writeTemp(t, dir, "a.yaml", "name: svc-a\nimage: a:latest")
	writeTemp(t, dir, "b.yml", "name: svc-b\nimage: b:latest")
	writeTemp(t, dir, "notes.txt", "ignored")

	manifests, err := manifest.LoadDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(manifests) != 2 {
		t.Errorf("expected 2 manifests, got %d", len(manifests))
	}
}
