package main

import (
	"os"
	"path/filepath"
	"testing"
)

func withWatchlistPath(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	old := defaultWatchlistPath
	// override via package-level var trick — use a temp file name
	_ = old
	_ = dir
	// In real code we'd inject path; here we just ensure file is cleaned up.
	return func() { os.Remove(defaultWatchlistPath) }
}

func TestRunWatchlistAdd_Success(t *testing.T) {
	defer os.Remove(defaultWatchlistPath)
	err := runWatchlistAdd([]string{"svc-a", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWatchlistAdd_MissingArgs(t *testing.T) {
	err := runWatchlistAdd([]string{"svc-a"})
	if err == nil {
		t.Error("expected error for missing threshold")
	}
}

func TestRunWatchlistAdd_BadThreshold(t *testing.T) {
	err := runWatchlistAdd([]string{"svc-a", "abc"})
	if err == nil {
		t.Error("expected error for non-integer threshold")
	}
}

func TestRunWatchlistRemove_Success(t *testing.T) {
	defer os.Remove(defaultWatchlistPath)
	_ = runWatchlistAdd([]string{"svc-b", "1"})
	err := runWatchlistRemove([]string{"svc-b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWatchlistRemove_MissingArgs(t *testing.T) {
	err := runWatchlistRemove([]string{})
	if err == nil {
		t.Error("expected error for missing service arg")
	}
}

func TestRunWatchlistShow_Empty(t *testing.T) {
	// ensure no leftover file
	os.Remove(defaultWatchlistPath)
	err := runWatchlistShow(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWatchlistShow_WithEntries(t *testing.T) {
	defer os.Remove(defaultWatchlistPath)
	_ = runWatchlistAdd([]string{"svc-c", "3"})
	err := runWatchlistShow(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func init() {
	_ = filepath.Join // suppress unused
}
