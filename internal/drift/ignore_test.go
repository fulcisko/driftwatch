package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeIgnoreResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []DiffEntry{
				{Key: "replicas", Expected: "3", Actual: "2"},
				{Key: "image", Expected: "v1", Actual: "v2"},
			},
		},
		{
			Service: "worker",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
	}
}

func TestLoadIgnoreList_NotFound(t *testing.T) {
	il, err := LoadIgnoreList("/nonexistent/ignore.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(il.Rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestAddAndLoadIgnoreRule(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ignore.json")

	if err := AddIgnoreRule(path, "api", "replicas", "managed by HPA"); err != nil {
		t.Fatalf("AddIgnoreRule: %v", err)
	}

	il, err := LoadIgnoreList(path)
	if err != nil {
		t.Fatalf("LoadIgnoreList: %v", err)
	}
	if len(il.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(il.Rules))
	}
	if il.Rules[0].Key != "replicas" {
		t.Errorf("unexpected key: %s", il.Rules[0].Key)
	}
}

func TestApplyIgnoreList_ExactMatch(t *testing.T) {
	il := &IgnoreList{Rules: []IgnoreRule{{Service: "api", Key: "replicas"}}}
	results := ApplyIgnoreList(makeIgnoreResults(), il)
	for _, r := range results {
		if r.Service == "api" {
			for _, d := range r.Diffs {
				if d.Key == "replicas" {
					t.Errorf("replicas should have been ignored")
				}
			}
		}
	}
}

func TestApplyIgnoreList_WildcardMatch(t *testing.T) {
	il := &IgnoreList{Rules: []IgnoreRule{{Key: "image*"}}}
	results := ApplyIgnoreList(makeIgnoreResults(), il)
	for _, r := range results {
		for _, d := range r.Diffs {
			if d.Key == "image" {
				t.Errorf("image should have been ignored by wildcard")
			}
		}
	}
}

func TestApplyIgnoreList_ServiceMismatch(t *testing.T) {
	il := &IgnoreList{Rules: []IgnoreRule{{Service: "other", Key: "timeout"}}}
	results := ApplyIgnoreList(makeIgnoreResults(), il)
	for _, r := range results {
		if r.Service == "worker" && len(r.Diffs) == 0 {
			t.Errorf("worker timeout should NOT be ignored (service mismatch)")
		}
	}
}

func TestSaveIgnoreList_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ignore.json")
	il := &IgnoreList{Rules: []IgnoreRule{
		{Service: "", Key: "debug_mode", Reason: "always off"},
	}}
	if err := SaveIgnoreList(path, il); err != nil {
		t.Fatalf("SaveIgnoreList: %v", err)
	}
	loaded, err := LoadIgnoreList(path)
	if err != nil {
		t.Fatalf("LoadIgnoreList: %v", err)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].Key != "debug_mode" {
		t.Errorf("round-trip mismatch")
	}
	_ = os.Remove(path)
}
