package drift

import (
	"strings"
	"testing"
	"time"
)

func makeTrendEntries() []TrendEntry {
	now := time.Now().UTC()
	return []TrendEntry{
		{Service: "api", DriftCount: 3, RecordedAt: now.Add(-72 * time.Hour)},
		{Service: "api", DriftCount: 5, RecordedAt: now.Add(-48 * time.Hour)},
		{Service: "api", DriftCount: 7, RecordedAt: now.Add(-24 * time.Hour)},
		{Service: "worker", DriftCount: 2, RecordedAt: now.Add(-96 * time.Hour)},
		{Service: "worker", DriftCount: 2, RecordedAt: now.Add(-48 * time.Hour)},
		{Service: "worker", DriftCount: 2, RecordedAt: now.Add(-24 * time.Hour)},
	}
}

func TestComputeVelocity_BasicRate(t *testing.T) {
	entries := makeTrendEntries()
	report := ComputeVelocity(entries, 0)
	if len(report.Entries) == 0 {
		t.Fatal("expected velocity entries, got none")
	}
	var found bool
	for _, e := range report.Entries {
		if e.Service == "api" {
			found = true
			if e.DriftsPerDay <= 0 {
				t.Errorf("expected positive drifts/day for api, got %f", e.DriftsPerDay)
			}
		}
	}
	if !found {
		t.Error("api service not found in velocity report")
	}
}

func TestComputeVelocity_DetectsAcceleration(t *testing.T) {
	entries := makeTrendEntries()
	report := ComputeVelocity(entries, 0)
	for _, e := range report.Entries {
		if e.Service == "api" && !e.Accelerating {
			t.Error("expected api to be marked as accelerating")
		}
		if e.Service == "worker" && e.Accelerating {
			t.Error("expected worker to not be marked as accelerating")
		}
	}
}

func TestComputeVelocity_MinSpanFilters(t *testing.T) {
	entries := makeTrendEntries()
	// require at least 200 hours span — no entries qualify
	report := ComputeVelocity(entries, 200)
	if len(report.Entries) != 0 {
		t.Errorf("expected 0 entries with high minSpan, got %d", len(report.Entries))
	}
}

func TestComputeVelocity_SortedByRate(t *testing.T) {
	entries := makeTrendEntries()
	report := ComputeVelocity(entries, 0)
	for i := 1; i < len(report.Entries); i++ {
		if report.Entries[i].DriftsPerDay > report.Entries[i-1].DriftsPerDay {
			t.Error("velocity entries not sorted descending by drifts/day")
		}
	}
}

func TestFormatVelocity_ContainsServiceName(t *testing.T) {
	entries := makeTrendEntries()
	report := ComputeVelocity(entries, 0)
	out := FormatVelocity(report)
	if !strings.Contains(out, "api") {
		t.Error("formatted velocity output missing service name 'api'")
	}
	if !strings.Contains(out, "drifts/day") {
		t.Error("formatted velocity output missing header 'drifts/day'")
	}
}

func TestFormatVelocity_EmptyReport(t *testing.T) {
	out := FormatVelocity(VelocityReport{})
	if !strings.Contains(out, "no velocity data") {
		t.Error("expected 'no velocity data' for empty report")
	}
}
