package drift

import (
	"strings"
	"testing"
)

func makeProjectionResults() []CompareResult {
	return []CompareResult{
		{
			Service: "api",
			Diffs: []DiffEntry{
				{Key: "replicas", ExpectedValue: "3", LiveValue: "2", DriftType: "changed"},
				{Key: "image", ExpectedValue: "v1.2", LiveValue: "v1.1", DriftType: "changed"},
			},
		},
		{
			Service: "worker",
			Diffs: []DiffEntry{
				{Key: "replicas", ExpectedValue: "2", LiveValue: "2", DriftType: "none"},
			},
		},
		{
			Service: "cache",
			Diffs:   []DiffEntry{},
		},
	}
}

func TestApplyProjection_SelectedFields(t *testing.T) {
	results := makeProjectionResults()
	opts := ProjectionOptions{
		Fields: []ProjectionField{
			{Key: "replicas"},
			{Key: "image", Alias: "img"},
		},
	}
	rows := ApplyProjection(results, opts)
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[0].Service != "api" {
		t.Errorf("expected api first, got %s", rows[0].Service)
	}
	if rows[0].Values["replicas"] != "2" {
		t.Errorf("expected live replicas=2, got %s", rows[0].Values["replicas"])
	}
	if rows[0].Values["img"] != "v1.1" {
		t.Errorf("expected img=v1.1, got %s", rows[0].Values["img"])
	}
}

func TestApplyProjection_ServiceFilter(t *testing.T) {
	results := makeProjectionResults()
	opts := ProjectionOptions{
		Fields:  []ProjectionField{{Key: "replicas"}},
		Service: "worker",
	}
	rows := ApplyProjection(results, opts)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0].Service != "worker" {
		t.Errorf("expected worker, got %s", rows[0].Service)
	}
}

func TestApplyProjection_MissingKeyIsEmpty(t *testing.T) {
	results := makeProjectionResults()
	opts := ProjectionOptions{
		Fields: []ProjectionField{{Key: "nonexistent"}},
	}
	rows := ApplyProjection(results, opts)
	for _, row := range rows {
		if v := row.Values["nonexistent"]; v != "" {
			t.Errorf("expected empty for missing key on %s, got %q", row.Service, v)
		}
	}
}

func TestFormatProjection_ContainsHeaders(t *testing.T) {
	results := makeProjectionResults()
	fields := []ProjectionField{{Key: "replicas"}, {Key: "image", Alias: "img"}}
	rows := ApplyProjection(results, ProjectionOptions{Fields: fields})
	out := FormatProjection(rows, fields)
	if !strings.Contains(out, "SERVICE") {
		t.Error("expected SERVICE header")
	}
	if !strings.Contains(out, "REPLICAS") {
		t.Error("expected REPLICAS header")
	}
	if !strings.Contains(out, "IMG") {
		t.Error("expected IMG header")
	}
	if !strings.Contains(out, "api") {
		t.Error("expected api row")
	}
}

func TestFormatProjection_Empty(t *testing.T) {
	out := FormatProjection(nil, nil)
	if !strings.Contains(out, "no results") {
		t.Errorf("expected 'no results', got %q", out)
	}
}
