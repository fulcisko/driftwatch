package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writePolicyFile(t *testing.T, p Policy) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	b, _ := json.Marshal(p)
	os.WriteFile(path, b, 0644)
	return path
}

func makePolicyResults() []CompareResult {
	return []CompareResult{
		{
			Service: "svc-a",
			Diffs: []DiffEntry{
				{Key: "env", Expected: "prod", Live: "staging"},
			},
		},
		{
			Service: "svc-b",
			Diffs: []DiffEntry{},
		},
	}
}

func TestLoadPolicy_Valid(t *testing.T) {
	p := Policy{Name: "test", Rules: []PolicyRule{{Key: "env", Severity: "high"}}}
	path := writePolicyFile(t, p)
	loaded, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Name != "test" {
		t.Errorf("expected name 'test', got %q", loaded.Name)
	}
}

func TestLoadPolicy_NotFound(t *testing.T) {
	_, err := LoadPolicy("/nonexistent/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestApplyPolicy_AllowedViolation(t *testing.T) {
	p := &Policy{
		Name: "env-policy",
		Rules: []PolicyRule{
			{Key: "env", Severity: "high", Allowed: []string{"prod", "staging"}},
		},
	}
	results := []CompareResult{
		{Service: "svc-a", Diffs: []DiffEntry{{Key: "env", Expected: "prod", Live: "dev"}}},
	}
	violations := ApplyPolicy(results, p)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Service != "svc-a" {
		t.Errorf("expected service svc-a, got %s", violations[0].Service)
	}
}

func TestApplyPolicy_NoViolation(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{Key: "env", Allowed: []string{"prod", "staging"}},
		},
	}
	results := []CompareResult{
		{Service: "svc-a", Diffs: []DiffEntry{{Key: "env", Expected: "prod", Live: "staging"}}},
	}
	violations := ApplyPolicy(results, p)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}
