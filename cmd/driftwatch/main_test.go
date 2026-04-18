package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_MissingArgs(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"driftwatch"}
	if err := run(); err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRun_EmptyManifestDir(t *testing.T) {
	dir := t.TempDir()

	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"driftwatch", dir, "http://localhost"}
	if err := run(); err == nil {
		t.Fatal("expected error for empty manifest dir")
	}
}

func TestRun_Success(t *testing.T) {
	dir := t.TempDir()
	yaml := `name: svc1
config:
  PORT: "8080"
  LOG_LEVEL: info
`
	if err := os.WriteFile(filepath.Join(dir, "svc1.yaml"), []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("PORT=8080\nLOG_LEVEL=info\n"))
	}))
	defer ts.Close()

	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"driftwatch", dir, ts.URL}
	if err := run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
