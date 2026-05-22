package drift

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makePromotionResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs:   []DiffEntry{},
		},
		{
			Service: "beta",
			Diffs: []DiffEntry{
				{Key: "replicas", Expected: "3", Actual: "1"},
			},
		},
		{
			Service: "gamma",
			Diffs: []DiffEntry{
				{Key: "secret_key", Expected: "abc", Actual: "xyz"},
			},
		},
	}
}

func TestAssessPromotion_SafeWhenNoDrift(t *testing.T) {
	results := makePromotionResults()
	entries := AssessPromotion(results[:1], StageDev, StageStaging)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !entries[0].Safe {
		t.Errorf("expected safe=true for clean service")
	}
	if entries[0].Service != "alpha" {
		t.Errorf("unexpected service: %s", entries[0].Service)
	}
}

func TestAssessPromotion_BlockedOnHighSeverity(t *testing.T) {
	results := makePromotionResults()
	entries := AssessPromotion(results[2:], StageStaging, StageProd)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Safe {
		t.Errorf("expected safe=false for high-severity drift")
	}
	if !strings.Contains(entries[0].Reason, "high") {
		t.Errorf("expected reason to mention high, got: %s", entries[0].Reason)
	}
}

func TestAssessPromotion_StagesRecorded(t *testing.T) {
	results := makePromotionResults()
	entries := AssessPromotion(results[:1], StageDev, StageProd)
	if entries[0].FromStage != StageDev || entries[0].ToStage != StageProd {
		t.Errorf("unexpected stages: %v -> %v", entries[0].FromStage, entries[0].ToStage)
	}
}

func TestSaveAndLoadPromotionReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "promotion.json")
	results := makePromotionResults()
	entries := AssessPromotion(results, StageDev, StageStaging)
	if err := SavePromotionReport(path, entries); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := LoadPromotionReport(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(loaded))
	}
}

func TestLoadPromotionReport_NotFound(t *testing.T) {
	entries, err := LoadPromotionReport("/nonexistent/promotion.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file")
	}
}

func TestFormatPromotion_ContainsServiceAndStatus(t *testing.T) {
	results := makePromotionResults()
	entries := AssessPromotion(results, StageDev, StageStaging)
	out := FormatPromotion(entries)
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected alpha in output")
	}
	if !strings.Contains(out, "SAFE") {
		t.Errorf("expected SAFE in output")
	}
}

func TestFormatPromotion_Empty(t *testing.T) {
	out := FormatPromotion(nil)
	if !strings.Contains(out, "no promotion") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
