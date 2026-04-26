package drift

import (
	"strings"
	"testing"
)

func makeClusterResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30s", Actual: "60s"},
				{Key: "replicas", Expected: "2", Actual: "3"},
			},
		},
		{
			Service: "beta",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30s", Actual: "90s"},
				{Key: "replicas", Expected: "2", Actual: "4"},
			},
		},
		{
			Service: "gamma",
			Diffs: []DiffEntry{
				{Key: "log_level", Expected: "info", Actual: "debug"},
			},
		},
		{
			Service: "clean",
			Diffs:   []DiffEntry{},
		},
	}
}

func TestClusterByDriftPattern_FindsCluster(t *testing.T) {
	results := makeClusterResults()
	cr := ClusterByDriftPattern(results, 2)

	if len(cr.Clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(cr.Clusters))
	}
	g := cr.Clusters[0]
	if len(g.Services) != 2 {
		t.Errorf("expected 2 services in cluster, got %d", len(g.Services))
	}
	if len(g.SharedKeys) != 2 {
		t.Errorf("expected 2 shared keys, got %d: %v", len(g.SharedKeys), g.SharedKeys)
	}
}

func TestClusterByDriftPattern_Unclustered(t *testing.T) {
	results := makeClusterResults()
	cr := ClusterByDriftPattern(results, 2)

	// gamma has unique key; clean has no diffs — neither should be in a cluster
	unclustered := cr.Unclustered
	found := false
	for _, svc := range unclustered {
		if svc == "gamma" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected gamma in unclustered, got %v", unclustered)
	}
}

func TestClusterByDriftPattern_SkipsCleanServices(t *testing.T) {
	results := makeClusterResults()
	cr := ClusterByDriftPattern(results, 1)

	for _, g := range cr.Clusters {
		for _, svc := range g.Services {
			if svc == "clean" {
				t.Errorf("clean service should not appear in any cluster")
			}
		}
	}
	for _, svc := range cr.Unclustered {
		if svc == "clean" {
			t.Errorf("clean service should not appear in unclustered")
		}
	}
}

func TestClusterByDriftPattern_NoClusters(t *testing.T) {
	results := []CompareResult{
		{Service: "svc-a", Diffs: []DiffEntry{{Key: "x", Expected: "1", Actual: "2"}}},
		{Service: "svc-b", Diffs: []DiffEntry{{Key: "y", Expected: "1", Actual: "2"}}},
	}
	cr := ClusterByDriftPattern(results, 1)
	if len(cr.Clusters) != 0 {
		t.Errorf("expected no clusters, got %d", len(cr.Clusters))
	}
}

func TestFormatCluster_ContainsClusterLabel(t *testing.T) {
	results := makeClusterResults()
	cr := ClusterByDriftPattern(results, 2)
	out := FormatCluster(cr)

	if !strings.Contains(out, "cluster-1") {
		t.Errorf("expected cluster-1 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Clusters found:") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}
