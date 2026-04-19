package drift

import (
	"strings"
	"testing"
)

func makeHeatmapResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []DiffEntry{
				{Key: "secret_key", Expected: "a", Actual: "b"},
				{Key: "secret_key", Expected: "a", Actual: "c"},
				{Key: "replicas", Expected: "2", Actual: "3"},
			},
		},
		{
			Service: "worker",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
		{
			Service: "clean",
			Diffs:   []DiffEntry{},
		},
	}
}

func TestBuildHeatmap_SkipsClean(t *testing.T) {
	rows := BuildHeatmap(makeHeatmapResults())
	for _, row := range rows {
		if row.Service == "clean" {
			t.Error("clean service should not appear in heatmap")
		}
	}
}

func TestBuildHeatmap_CountsCorrect(t *testing.T) {
	rows := BuildHeatmap(makeHeatmapResults())
	var apiRow *HeatmapRow
	for i := range rows {
		if rows[i].Service == "api" {
			apiRow = &rows[i]
		}
	}
	if apiRow == nil {
		t.Fatal("expected api row")
	}
	if apiRow.Total != 3 {
		t.Errorf("expected total 3, got %d", apiRow.Total)
	}
	if apiRow.Entries[0].Key != "secret_key" {
		t.Errorf("expected secret_key as top entry, got %s", apiRow.Entries[0].Key)
	}
	if apiRow.Entries[0].Count != 2 {
		t.Errorf("expected count 2, got %d", apiRow.Entries[0].Count)
	}
}

func TestBuildHeatmap_SeveritySet(t *testing.T) {
	rows := BuildHeatmap(makeHeatmapResults())
	for _, row := range rows {
		if row.Service == "api" {
			for _, e := range row.Entries {
				if e.Key == "secret_key" && e.MaxSev != "high" {
					t.Errorf("expected high severity for secret_key, got %s", e.MaxSev)
				}
			}
		}
	}
}

func TestFormatHeatmap_ContainsHeaders(t *testing.T) {
	rows := BuildHeatmap(makeHeatmapResults())
	out := FormatHeatmap(rows)
	for _, want := range []string{"SERVICE", "KEY", "COUNT", "MAX_SEV", "api", "secret_key"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestFormatHeatmap_Empty(t *testing.T) {
	out := FormatHeatmap(nil)
	if !strings.Contains(out, "no drift") {
		t.Error("expected empty message")
	}
}
