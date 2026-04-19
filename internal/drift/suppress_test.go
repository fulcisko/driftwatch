package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func suppressPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "suppress.json")
}

func TestAddAndLoadSuppressRule(t *testing.T) {
	p := suppressPath(t)
	rule := SuppressRule{
		Service:   "api",
		Key:       "replicas",
		Reason:    "planned maintenance",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := AddSuppressRule(p, rule); err != nil {
		t.Fatalf("add: %v", err)
	}
	sl, err := LoadSuppressList(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(sl.Rules) != 1 || sl.Rules[0].Key != "replicas" {
		t.Errorf("unexpected rules: %+v", sl.Rules)
	}
}

func TestLoadSuppressList_NotFound(t *testing.T) {
	sl, err := LoadSuppressList("/nonexistent/suppress.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(sl.Rules) != 0 {
		t.Errorf("expected empty list")
	}
}

func TestApplySuppress_RemovesSuppressedDiff(t *testing.T) {
	results := []CompareResult{
		{Service: "api", Diffs: []DiffEntry{{Key: "replicas"}, {Key: "image"}}},
	}
	sl := SuppressList{Rules: []SuppressRule{
		{Service: "api", Key: "replicas", ExpiresAt: time.Now().Add(time.Hour)},
	}}
	out := ApplySuppress(results, sl)
	if len(out[0].Diffs) != 1 || out[0].Diffs[0].Key != "image" {
		t.Errorf("expected only 'image' diff, got %+v", out[0].Diffs)
	}
}

func TestApplySuppress_ExpiredRuleIgnored(t *testing.T) {
	results := []CompareResult{
		{Service: "api", Diffs: []DiffEntry{{Key: "replicas"}}},
	}
	sl := SuppressList{Rules: []SuppressRule{
		{Service: "api", Key: "replicas", ExpiresAt: time.Now().Add(-time.Hour)},
	}}
	out := ApplySuppress(results, sl)
	if len(out[0].Diffs) != 1 {
		t.Errorf("expired rule should not suppress diff")
	}
}

func TestApplySuppress_WildcardKey(t *testing.T) {
	results := []CompareResult{
		{Service: "worker", Diffs: []DiffEntry{{Key: "cpu"}, {Key: "mem"}}},
	}
	sl := SuppressList{Rules: []SuppressRule{
		{Service: "worker", Key: "*", ExpiresAt: time.Now().Add(time.Hour)},
	}}
	out := ApplySuppress(results, sl)
	if len(out[0].Diffs) != 0 {
		t.Errorf("wildcard should suppress all diffs")
	}
}

func TestSaveSuppressList_Roundtrip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "s.json")
	sl := SuppressList{Rules: []SuppressRule{
		{Service: "svc", Key: "timeout", Reason: "test", ExpiresAt: time.Now().Add(time.Hour).Truncate(time.Second)},
	}}
	if err := SaveSuppressList(p, sl); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadSuppressList(p)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Rules[0].Service != "svc" {
		t.Errorf("roundtrip mismatch")
	}
	os.Remove(p)
}
