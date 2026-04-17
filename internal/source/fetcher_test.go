package source_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/source"
)

func TestFetch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("IMAGE=nginx:1.25\nREPLICAS=3\n"))
	}))
	defer server.Close()

	f := source.NewFetcher(server.URL, 5*time.Second)
	cfg, err := f.Fetch(context.Background(), "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ServiceName != "web" {
		t.Errorf("expected service name 'web', got %q", cfg.ServiceName)
	}
	if cfg.Fields["IMAGE"] != "nginx:1.25" {
		t.Errorf("expected IMAGE=nginx:1.25, got %q", cfg.Fields["IMAGE"])
	}
}

func TestFetch_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	f := source.NewFetcher(server.URL, 5*time.Second)
	_, err := f.Fetch(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestFetch_InvalidBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("BADLINE\n"))
	}))
	defer server.Close()

	f := source.NewFetcher(server.URL, 5*time.Second)
	_, err := f.Fetch(context.Background(), "svc")
	if err == nil {
		t.Fatal("expected parse error")
	}
}
