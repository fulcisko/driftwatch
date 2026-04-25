package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCorrelationManifest(t *testing.T, dir, name, content string) {
	t.Helper()
	p := filepath.Join(dir, name+".yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeCorrelationManifest: %v", err)
	}
}

func TestRunCorrelation_MissingArgs(t *testing.T) {
	err := runCorrelation([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunCorrelation_EmptyManifestDir(t *testing.T) {
	dir := t.TempDir()
	err := runCorrelation([]string{dir, "http://localhost:9999"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunCorrelation_BadManifestDir(t *testing.T) {
	err := runCorrelation([]string{"/nonexistent/path", "http://localhost:9999"})
	if err == nil {
		t.Fatal("expected error for bad manifest dir")
	}
}

func TestRunCorrelation_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	writeCorrelationManifest(t, dir, "svc-a", "name: svc-a\nconfig:\n  timeout: 30s\n")

	// Redirect stdout to capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCorrelation([]string{dir, "http://localhost:19999", "--format=json"})

	w.Close()
	os.Stdout = old
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	output := string(buf[:n])
	_ = output

	// Source will fail to fetch but correlation should still return cleanly
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
