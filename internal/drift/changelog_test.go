package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeChangelogEntry(service, key string) ChangelogEntry {
	return ChangelogEntry{
		Timestamp: time.Now().UTC(),
		Service:   service,
		Diffs: []CompareResult{
			{Key: key, Expected: "a", Actual: "b", Drifted: true},
		},
	}
}

func TestAppendAndLoadChangelog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "changelog.json")

	if err := AppendChangelog(path, makeChangelogEntry("svc-a", "port")); err != nil {
		t.Fatalf("first append: %v", err)
	}
	if err := AppendChangelog(path, makeChangelogEntry("svc-b", "timeout")); err != nil {
		t.Fatalf("second append: %v", err)
	}

	cl, err := LoadChangelog(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(cl) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cl))
	}
}

func TestLoadChangelog_NotFound(t *testing.T) {
	cl, err := LoadChangelog("/nonexistent/changelog.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cl) != 0 {
		t.Fatalf("expected empty changelog")
	}
}

func TestFilterChangelog_ByService(t *testing.T) {
	path := filepath.Join(t.TempDir(), "changelog.json")
	_ = AppendChangelog(path, makeChangelogEntry("alpha", "k1"))
	_ = AppendChangelog(path, makeChangelogEntry("beta", "k2"))
	_ = AppendChangelog(path, makeChangelogEntry("alpha", "k3"))

	cl, _ := LoadChangelog(path)
	filtered := FilterChangelog(cl, "alpha")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 alpha entries, got %d", len(filtered))
	}
}

func TestFilterChangelog_EmptyService(t *testing.T) {
	path := filepath.Join(t.TempDir(), "changelog.json")
	_ = AppendChangelog(path, makeChangelogEntry("x", "k"))
	_ = AppendChangelog(path, makeChangelogEntry("y", "k"))

	cl, _ := LoadChangelog(path)
	all := FilterChangelog(cl, "")
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestAppendChangelog_SetsTimestamp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "changelog.json")
	before := time.Now().UTC()
	entry := ChangelogEntry{Service: "svc", Diffs: []CompareResult{}}
	_ = AppendChangelog(path, entry)
	cl, _ := LoadChangelog(path)
	if cl[0].Timestamp.Before(before) {
		t.Error("timestamp should be set automatically")
	}
	_ = os.Remove(path)
}
