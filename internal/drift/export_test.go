package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func exportResults() map[string][]CompareResult {
	return map[string][]CompareResult{
		"api": {
			{Key: "replicas", Expected: "2", Actual: "1", Kind: KindChanged},
		},
		"worker": {},
	}
}

func TestExportJSON_ValidOutput(t *testing.T) {
	results := exportResults()
	stats := Summarize(results)
	var buf bytes.Buffer
	if err := ExportJSON(&buf, results, stats); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["summary"]; !ok {
		t.Error("expected 'summary' key in JSON output")
	}
	if _, ok := out["services"]; !ok {
		t.Error("expected 'services' key in JSON output")
	}
}

func TestExportText_ContainsServiceAndSummary(t *testing.T) {
	results := exportResults()
	stats := Summarize(results)
	var buf bytes.Buffer
	if err := ExportText(&buf, results, stats); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "api") {
		t.Error("expected service name 'api' in text output")
	}
	if !strings.Contains(out, "replicas") {
		t.Error("expected key 'replicas' in text output")
	}
	if !strings.Contains(out, "Services checked") {
		t.Error("expected summary line in text output")
	}
}

func TestExportText_SkipsCleanServices(t *testing.T) {
	results := exportResults()
	stats := Summarize(results)
	var buf bytes.Buffer
	_ = ExportText(&buf, results, stats)
	if strings.Contains(buf.String(), "worker") {
		t.Error("clean service 'worker' should not appear in text output")
	}
}
