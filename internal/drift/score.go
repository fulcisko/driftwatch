package drift

import "fmt"

// DriftScore holds a numeric risk score for a service.
type DriftScore struct {
	Service  string  `json:"service"`
	Score    float64 `json:"score"`
	Drifted  int     `json:"drifted"`
	Highs    int     `json:"highs"`
	Mediums  int     `json:"mediums"`
	Lows     int     `json:"lows"`
}

// weightings for severity levels
const (
	weightHigh   = 10.0
	weightMedium = 5.0
	weightLow    = 1.0
)

// ScoreResults computes a DriftScore for each CompareResult.
func ScoreResults(results []CompareResult) []DriftScore {
	scores := make([]DriftScore, 0, len(results))
	for _, r := range results {
		ds := DriftScore{Service: r.Service}
		for _, d := range r.Diffs {
			switch ClassifyKey(d.Key) {
			case SeverityHigh:
				ds.Highs++
				ds.Score += weightHigh
			case SeverityMedium:
				ds.Mediums++
				ds.Score += weightMedium
			case SeverityLow:
				ds.Lows++
				ds.Score += weightLow
			default:
				ds.Score += weightLow
			}
			ds.Drifted++
		}
		scores = append(scores, ds)
	}
	return scores
}

// FormatScore returns a human-readable string for a DriftScore.
func FormatScore(ds DriftScore) string {
	return fmt.Sprintf("service=%-20s score=%.1f drifted=%d (high=%d medium=%d low=%d)",
		ds.Service, ds.Score, ds.Drifted, ds.Highs, ds.Mediums, ds.Lows)
}
