package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeFingerprintResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs: []Diff{
				{Key: "timeout", Expected: "30s", Actual: "60s"},
				{Key: "replicas", Expected: "3", Actual: "1"},
			},
		},
		{
			Service: "beta",
			Diffs: []Diff{
				{Key: "log_level", Expected: "info", Actual: "debug"},
			},
		},
		{
			Service: "gamma",
			Diffs:   []Diff{},
		},
	}
}

func TestBuildFingerprint_StableHash(t *testing.T) {
	r := makeFingerprintResults()[0]
	fp1 := BuildFingerprint(r)
	fp2 := BuildFingerprint(r)
	if fp1.Hash != fp2.Hash {
		t.Errorf("expected stable hash, got %s vs %s", fp1.Hash, fp2.Hash)
	}
	if fp1.DriftCount != 2 {
		t.Errorf("expected drift count 2, got %d", fp1.DriftCount)
	}
}

func TestBuildFingerprintStore_SkipsClean(t *testing.T) {
	results := makeFingerprintResults()
	store := BuildFingerprintStore(results)
	if _, ok := store["gamma"]; ok {
		t.Error("expected clean service gamma to be excluded from store")
	}
	if len(store) != 2 {
		t.Errorf("expected 2 entries, got %d", len(store))
	}
}

func TestSaveAndLoadFingerprintStore(t *testing.T) {
	results := makeFingerprintResults()
	store := BuildFingerprintStore(results)

	dir := t.TempDir()
	path := filepath.Join(dir, "fingerprints.json")

	if err := SaveFingerprintStore(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadFingerprintStore(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["alpha"].Hash != store["alpha"].Hash {
		t.Errorf("hash mismatch for alpha")
	}
}

func TestLoadFingerprintStore_NotFound(t *testing.T) {
	store, err := LoadFingerprintStore("/nonexistent/fingerprints.json")
	if err != nil {
		t.Fatalf("expected empty store, got error: %v", err)
	}
	if len(store) != 0 {
		t.Errorf("expected empty store, got %d entries", len(store))
	}
}

func TestDiffFingerprintStore_DetectsChange(t *testing.T) {
	old := FingerprintStore{
		"alpha": {Service: "alpha", Hash: "aabbccdd11223344"},
	}
	current := FingerprintStore{
		"alpha": {Service: "alpha", Hash: "deadbeefcafebabe"},
		"beta":  {Service: "beta", Hash: "1234567890abcdef"},
	}
	changed := DiffFingerprintStore(old, current)
	if len(changed) != 2 {
		t.Errorf("expected 2 changed services, got %d: %v", len(changed), changed)
	}
}

func TestDiffFingerprintStore_NoDifference(t *testing.T) {
	store := FingerprintStore{
		"alpha": {Service: "alpha", Hash: "aabbccdd11223344"},
	}
	changed := DiffFingerprintStore(store, store)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func TestFormatFingerprintStore_Empty(t *testing.T) {
	out := FormatFingerprintStore(FingerprintStore{})
	if out == "" {
		t.Error("expected non-empty output for empty store")
	}
}

func TestFormatFingerprintStore_ContainsService(t *testing.T) {
	_ = os.Getenv // suppress unused import
	results := makeFingerprintResults()
	store := BuildFingerprintStore(results)
	out := FormatFingerprintStore(store)
	if !containsStr(out, "alpha") {
		t.Errorf("expected output to contain 'alpha', got: %s", out)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
