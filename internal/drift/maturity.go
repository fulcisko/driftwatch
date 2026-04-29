package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// MaturityLevel represents how mature/stable a service's config is.
type MaturityLevel string

const (
	MaturityUnknown  MaturityLevel = "unknown"
	MaturityUnstable MaturityLevel = "unstable"
	MaturityDeveloping MaturityLevel = "developing"
	MaturityStable   MaturityLevel = "stable"
	MaturityMature   MaturityLevel = "mature"
)

// MaturityEntry holds the maturity assessment for a single service.
type MaturityEntry struct {
	Service     string        `json:"service"`
	Level       MaturityLevel `json:"level"`
	DriftScore  float64       `json:"drift_score"`
	AssessedAt  time.Time     `json:"assessed_at"`
	Note        string        `json:"note,omitempty"`
}

// MaturityReport is the full set of maturity entries.
type MaturityReport struct {
	Entries []MaturityEntry `json:"entries"`
}

// AssessMaturity derives a maturity level from scored drift results.
func AssessMaturity(results []CompareResult) []MaturityEntry {
	scores := ScoreResults(results)
	entries := make([]MaturityEntry, 0, len(scores))
	for _, s := range scores {
		entries = append(entries, MaturityEntry{
			Service:    s.Service,
			Level:      scoreToMaturity(s.Score),
			DriftScore: s.Score,
			AssessedAt: time.Now().UTC(),
		})
	}
	return entries
}

func scoreToMaturity(score float64) MaturityLevel {
	switch {
	case score == 0:
		return MaturityMature
	case score <= 5:
		return MaturityStable
	case score <= 15:
		return MaturityDeveloping
	case score <= 30:
		return MaturityUnstable
	default:
		return MaturityUnknown
	}
}

// SaveMaturityReport persists the maturity report to a JSON file.
func SaveMaturityReport(path string, entries []MaturityEntry) error {
	report := MaturityReport{Entries: entries}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal maturity report: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadMaturityReport reads a maturity report from a JSON file.
func LoadMaturityReport(path string) ([]MaturityEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []MaturityEntry{}, nil
		}
		return nil, fmt.Errorf("read maturity report: %w", err)
	}
	var report MaturityReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unmarshal maturity report: %w", err)
	}
	return report.Entries, nil
}

// FormatMaturity returns a human-readable summary of maturity entries.
func FormatMaturity(entries []MaturityEntry) string {
	if len(entries) == 0 {
		return "no maturity data available\n"
	}
	out := "Service Maturity Report\n"
	out += fmt.Sprintf("%-30s %-12s %s\n", "SERVICE", "LEVEL", "SCORE")
	out += fmt.Sprintf("%-30s %-12s %s\n", "-------", "-----", "-----")
	for _, e := range entries {
		out += fmt.Sprintf("%-30s %-12s %.1f\n", e.Service, e.Level, e.DriftScore)
	}
	return out
}
