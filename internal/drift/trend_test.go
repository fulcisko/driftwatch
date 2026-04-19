package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTrendResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []DiffEntry{
				{Key: "replicas", Expected: "3", Actual: "2"},
			},
		},
		{
			Service: "worker",
			Diffs:   []DiffEntry{},
		},
	}
}

func TestAppendAndLoadTrend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "trend.json")
	results := makeTrendResults()
	if err := AppendTrend(path, results); err != nil {
		t.Fatalf("AppendTrend: %v", err)
	}
	report, err := LoadTrend(path)
	if err != nil {
		t.Fatalf("LoadTrend: %v", err)
	}
	if len(report.Points) != 1 {
		t.Errorf("expected 1 point (clean services skipped), got %d", len(report.Points))
	}
	if report.Points[0].Service != "api" {
		t.Errorf("expected service api, got %s", report.Points[0].Service)
	}
}

func TestLoadTrend_NotFound(t *testing.T) {
	report, err := LoadTrend(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Points) != 0 {
		t.Errorf("expected empty report")
	}
}

func TestFilterTrend_ByService(t *testing.T) {
	path := filepath.Join(t.TempDir(), "trend.json")
	results := []CompareResult{
		{Service: "api", Diffs: []DiffEntry{{Key: "x", Expected: "1", Actual: "2"}}},
		{Service: "db", Diffs: []DiffEntry{{Key: "y", Expected: "a", Actual: "b"}}},
	}
	_ = AppendTrend(path, results)
	report, _ := LoadTrend(path)
	points := FilterTrend(report, "api")
	if len(points) != 1 || points[0].Service != "api" {
		t.Errorf("expected 1 api point, got %+v", points)
	}
}

func TestFormatTrend_Empty(t *testing.T) {
	out := FormatTrend(nil)
	if out != "no trend data available\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatTrend_ContainsService(t *testing.T) {
	path := filepath.Join(t.TempDir(), "trend.json")
	_ = AppendTrend(path, makeTrendResults())
	report, _ := LoadTrend(path)
	out := FormatTrend(report.Points)
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
	_ = os.Remove(path)
}
