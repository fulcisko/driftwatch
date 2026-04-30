package drift

import (
	"strings"
	"testing"
)

func makeExposureResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api-gateway",
			Diffs: []DiffEntry{
				{Key: "replicas", Expected: "3", Actual: "1"},
				{Key: "tls_enabled", Expected: "true", Actual: "false"},
			},
		},
		{
			Service: "auth-service",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
		{
			Service: "cache",
			Diffs: []DiffEntry{},
		},
	}
}

func TestAssessExposure_SkipsClean(t *testing.T) {
	results := makeExposureResults()
	levelMap := map[string]ExposureLevel{
		"api-gateway":  ExposurePublic,
		"auth-service": ExposureInternal,
		"cache":        ExposurePrivate,
	}

	entries := AssessExposure(results, levelMap)
	for _, e := range entries {
		if e.Service == "cache" {
			t.Errorf("expected clean service 'cache' to be skipped")
		}
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestAssessExposure_PublicHigherRisk(t *testing.T) {
	results := makeExposureResults()
	levelMap := map[string]ExposureLevel{
		"api-gateway":  ExposurePublic,
		"auth-service": ExposurePrivate,
	}

	entries := AssessExposure(results, levelMap)
	if len(entries) < 2 {
		t.Fatalf("expected at least 2 entries")
	}
	if entries[0].Service != "api-gateway" {
		t.Errorf("expected api-gateway first (highest risk), got %s", entries[0].Service)
	}
	if entries[0].RiskScore <= entries[1].RiskScore {
		t.Errorf("public service should have higher risk score than private")
	}
}

func TestAssessExposure_UnknownDefault(t *testing.T) {
	results := []CompareResult{
		{Service: "mystery", Diffs: []DiffEntry{{Key: "port", Expected: "80", Actual: "8080"}}},
	}
	entries := AssessExposure(results, map[string]ExposureLevel{})
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if entries[0].Level != ExposureUnknown {
		t.Errorf("expected unknown level, got %s", entries[0].Level)
	}
}

func TestAssessExposure_TopKeysLimit(t *testing.T) {
	diffs := make([]DiffEntry, 8)
	for i := range diffs {
		diffs[i] = DiffEntry{Key: fmt.Sprintf("key%d", i), Expected: "a", Actual: "b"}
	}
	results := []CompareResult{{Service: "svc", Diffs: diffs}}
	entries := AssessExposure(results, map[string]ExposureLevel{"svc": ExposurePublic})
	if len(entries[0].TopKeys) > 5 {
		t.Errorf("expected at most 5 top keys, got %d", len(entries[0].TopKeys))
	}
}

func TestFormatExposure_ContainsHeaders(t *testing.T) {
	results := makeExposureResults()
	levelMap := map[string]ExposureLevel{
		"api-gateway":  ExposurePublic,
		"auth-service": ExposureInternal,
	}
	entries := AssessExposure(results, levelMap)
	out := FormatExposure(entries)
	for _, want := range []string{"service", "exposure", "risk", "api-gateway", "public"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestFormatExposure_Empty(t *testing.T) {
	out := FormatExposure(nil)
	if !strings.Contains(out, "no exposure") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
