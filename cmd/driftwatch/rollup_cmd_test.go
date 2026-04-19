package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRollupManifest(t *testing.T, dir, name, content string) {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRunRollup_MissingArgs(t *testing.T) {
	err := runRollup([]string{}, "text")
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRunRollup_NoDrift(t *testing.T) {
	dir := t.TempDir()
	writeRollupManifest(t, dir, "svc.yaml", "name: svc\nconfig:\n  timeout: \"30\"\n")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "timeout=30")
	}))
	defer ts.Close()

	err := runRollup([]string{dir, ts.URL}, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRollup_WithDrift(t *testing.T) {
	dir := t.TempDir()
	writeRollupManifest(t, dir, "svc.yaml", "name: svc\nconfig:\n  timeout: \"30\"\n")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "timeout=99")
	}))
	defer ts.Close()

	err := runRollup([]string{dir, ts.URL}, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRollup_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	writeRollupManifest(t, dir, "svc.yaml", "name: svc\nconfig:\n  port: \"8080\"\n")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "port=9090")
	}))
	defer ts.Close()

	err := runRollup([]string{dir, ts.URL}, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
