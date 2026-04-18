package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

func writePolicyJSON(t *testing.T, p drift.Policy) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	b, _ := json.Marshal(p)
	os.WriteFile(path, b, 0644)
	return path
}

func TestRunPolicyCheck_MissingArgs(t *testing.T) {
	err := runPolicyCheck([]string{"only-one"})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunPolicyCheck_BadPolicyFile(t *testing.T) {
	dir := t.TempDir()
	err := runPolicyCheck([]string{dir, "http://localhost:9", "/no/such/policy.json"})
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestRunPolicyCheck_NoPolicyViolations(t *testing.T) {
	policyPath := writePolicyJSON(t, drift.Policy{
		Name:  "test",
		Rules: []drift.PolicyRule{{Key: "env", Allowed: []string{"prod"}}},
	})
	manifestDir := t.TempDir()
	yaml := "name: svc-a\nconfig:\n  env: prod\n"
	os.WriteFile(filepath.Join(manifestDir, "svc-a.yaml"), []byte(yaml), 0644)

	// source URL unreachable — fetcher will warn and skip, no results => no violations
	err := runPolicyCheck([]string{manifestDir, "http://127.0.0.1:19999", policyPath})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
