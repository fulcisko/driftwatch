package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func quotaPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "quotas.json")
}

func makeQuotaResults() []CompareResult {
	return []CompareResult{
		{
			Service: "svc-a",
			Diffs: []DiffEntry{
				{Key: "port", Expected: "8080", Actual: "9090"},
				{Key: "replicas", Expected: "3", Actual: "1"},
			},
		},
		{
			Service: "svc-b",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30s", Actual: "60s"},
			},
		},
		{Service: "svc-c", Diffs: []DiffEntry{}},
	}
}

func TestAddAndLoadQuota(t *testing.T) {
	p := quotaPath(t)
	if err := AddQuotaRule(p, "svc-a", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ql, err := LoadQuotas(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(ql.Rules) != 1 || ql.Rules[0].Service != "svc-a" || ql.Rules[0].MaxDrift != 5 {
		t.Errorf("unexpected rules: %+v", ql.Rules)
	}
}

func TestAddQuota_UpdatesExisting(t *testing.T) {
	p := quotaPath(t)
	_ = AddQuotaRule(p, "svc-a", 3)
	_ = AddQuotaRule(p, "svc-a", 10)
	ql, _ := LoadQuotas(p)
	if len(ql.Rules) != 1 || ql.Rules[0].MaxDrift != 10 {
		t.Errorf("expected updated max_drift=10, got %+v", ql.Rules)
	}
}

func TestAddQuota_MissingService(t *testing.T) {
	p := quotaPath(t)
	if err := AddQuotaRule(p, "", 5); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestAddQuota_NegativeMax(t *testing.T) {
	p := quotaPath(t)
	if err := AddQuotaRule(p, "svc-a", -1); err == nil {
		t.Error("expected error for negative max_drift")
	}
}

func TestLoadQuotas_NotFound(t *testing.T) {
	ql, err := LoadQuotas("/nonexistent/quotas.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(ql.Rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestCheckQuotas_Violation(t *testing.T) {
	p := quotaPath(t)
	_ = AddQuotaRule(p, "svc-a", 1) // svc-a has 2 diffs → violation
	_ = AddQuotaRule(p, "svc-b", 5) // svc-b has 1 diff → ok
	results := makeQuotaResults()
	violations, err := CheckQuotas(p, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 || violations[0].Service != "svc-a" {
		t.Errorf("expected violation for svc-a, got %+v", violations)
	}
	if violations[0].Actual != 2 || violations[0].Allowed != 1 {
		t.Errorf("unexpected violation counts: %+v", violations[0])
	}
}

func TestCheckQuotas_NoViolation(t *testing.T) {
	p := quotaPath(t)
	_ = AddQuotaRule(p, "svc-a", 10)
	results := makeQuotaResults()
	violations, err := CheckQuotas(p, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %+v", violations)
	}
}

func TestCheckQuotas_UnmatchedServiceIgnored(t *testing.T) {
	p := quotaPath(t)
	_ = AddQuotaRule(p, "svc-z", 0)
	results := makeQuotaResults()
	violations, err := CheckQuotas(p, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations for unmatched service")
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
