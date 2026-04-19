package drift

import (
	"strings"
	"testing"
)

func makeScoreResults() []CompareResult {
	return []CompareResult{
		{
			Service: "alpha",
			Diffs: []DiffEntry{
				{Key: "secret_key", Expected: "x", Actual: "y"},
				{Key: "replicas", Expected: "2", Actual: "3"},
			},
		},
		{
			Service: "beta",
			Diffs: []DiffEntry{},
		},
		{
			Service: "gamma",
			Diffs: []DiffEntry{
				{Key: "timeout", Expected: "30", Actual: "60"},
			},
		},
	}
}

func TestScoreResults_HighKey(t *testing.T) {
	results := makeScoreResults()
	scores := ScoreResults(results)
	if len(scores) != 3 {
		t.Fatalf("expected 3 scores, got %d", len(scores))
	}
	alpha := scores[0]
	if alpha.Highs == 0 {
		t.Error("expected at least one high severity diff for alpha")
	}
	if alpha.Score < weightHigh {
		t.Errorf("expected score >= %.1f, got %.1f", weightHigh, alpha.Score)
	}
}

func TestScoreResults_CleanService(t *testing.T) {
	results := makeScoreResults()
	scores := ScoreResults(results)
	beta := scores[1]
	if beta.Score != 0 {
		t.Errorf("expected score 0 for clean service, got %.1f", beta.Score)
	}
	if beta.Drifted != 0 {
		t.Errorf("expected 0 drifted for clean service, got %d", beta.Drifted)
	}
}

func TestScoreResults_Count(t *testing.T) {
	results := makeScoreResults()
	scores := ScoreResults(results)
	if len(scores) != len(results) {
		t.Errorf("expected %d scores, got %d", len(results), len(scores))
	}
}

func TestFormatScore_ContainsService(t *testing.T) {
	ds := DriftScore{Service: "mysvc", Score: 15.0, Drifted: 2, Highs: 1, Mediums: 1}
	out := FormatScore(ds)
	if !strings.Contains(out, "mysvc") {
		t.Error("expected service name in formatted score")
	}
	if !strings.Contains(out, "15.0") {
		t.Error("expected score value in formatted score")
	}
}
