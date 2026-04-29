package drift

import (
	"fmt"
	"math"
	"sort"
)

// ConfidenceLevel represents how confident we are in a drift result.
type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "high"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceLow    ConfidenceLevel = "low"
)

// ConfidenceResult holds confidence scoring for a single service.
type ConfidenceResult struct {
	Service    string          `json:"service"`
	Score      float64         `json:"score"`
	Level      ConfidenceLevel `json:"level"`
	DriftCount int             `json:"drift_count"`
	Reason     string          `json:"reason"`
}

// ScoreConfidence computes a confidence score for each service based on
// historical consistency and current drift count. A higher score means
// we are more confident the drift is real and significant.
func ScoreConfidence(results []CompareResult, history []HistoryEntry) []ConfidenceResult {
	// Build a map of how many times each service has drifted historically.
	historyCounts := map[string]int{}
	for _, h := range history {
		for _, r := range h.Results {
			if len(r.Diffs) > 0 {
				historyCounts[r.Service]++
			}
		}
	}

	var out []ConfidenceResult
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		histFreq := historyCounts[r.Service]
		driftCount := len(r.Diffs)

		// Score: blend of current drift magnitude and historical frequency.
		// Normalize drift count with a soft cap at 10.
		driftNorm := math.Min(float64(driftCount)/10.0, 1.0)
		// Historical frequency normalized with soft cap at 5.
		histNorm := math.Min(float64(histFreq)/5.0, 1.0)
		score := math.Round((0.6*driftNorm+0.4*histNorm)*100) / 100

		level := levelFromScore(score)
		reason := buildConfidenceReason(driftCount, histFreq, level)

		out = append(out, ConfidenceResult{
			Service:    r.Service,
			Score:      score,
			Level:      level,
			DriftCount: driftCount,
			Reason:     reason,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Score > out[j].Score
	})
	return out
}

func levelFromScore(score float64) ConfidenceLevel {
	switch {
	case score >= 0.65:
		return ConfidenceHigh
	case score >= 0.35:
		return ConfidenceMedium
	default:
		return ConfidenceLow
	}
}

func buildConfidenceReason(driftCount, histFreq int, level ConfidenceLevel) string {
	return fmt.Sprintf("%d diffs detected, seen drifting %d time(s) historically — confidence %s",
		driftCount, histFreq, level)
}

// FormatConfidence returns a human-readable summary of confidence results.
func FormatConfidence(results []ConfidenceResult) string {
	if len(results) == 0 {
		return "no drifted services to score\n"
	}
	out := "drift confidence scores:\n"
	for _, r := range results {
		out += fmt.Sprintf("  %-30s score=%.2f  level=%-6s  diffs=%d\n",
			r.Service, r.Score, r.Level, r.DriftCount)
	}
	return out
}
