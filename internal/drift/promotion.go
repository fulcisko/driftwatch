package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PromotionStage represents an environment stage in a promotion pipeline.
type PromotionStage string

const (
	StageDev     PromotionStage = "dev"
	StageStaging PromotionStage = "staging"
	StageProd    PromotionStage = "prod"
)

// PromotionEntry records whether a service is safe to promote based on drift.
type PromotionEntry struct {
	Service   string         `json:"service"`
	FromStage PromotionStage `json:"from_stage"`
	ToStage   PromotionStage `json:"to_stage"`
	Safe      bool           `json:"safe"`
	Reason    string         `json:"reason"`
	CheckedAt time.Time      `json:"checked_at"`
}

// AssessPromotion evaluates whether each service in results is safe to promote.
// A service is considered safe if it has no high-severity drift.
func AssessPromotion(results []CompareResult, from, to PromotionStage) []PromotionEntry {
	entries := make([]PromotionEntry, 0, len(results))
	for _, r := range results {
		if len(r.Diffs) == 0 {
			entries = append(entries, PromotionEntry{
				Service:   r.Service,
				FromStage: from,
				ToStage:   to,
				Safe:      true,
				Reason:    "no drift detected",
				CheckedAt: time.Now().UTC(),
			})
			continue
		}
		max := MaxSeverity(r.Diffs)
		safe := max < SeverityHigh
		reason := fmt.Sprintf("max severity: %s", max)
		entries = append(entries, PromotionEntry{
			Service:   r.Service,
			FromStage: from,
			ToStage:   to,
			Safe:      safe,
			Reason:    reason,
			CheckedAt: time.Now().UTC(),
		})
	}
	return entries
}

// SavePromotionReport writes promotion entries to a JSON file.
func SavePromotionReport(path string, entries []PromotionEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPromotionReport reads promotion entries from a JSON file.
func LoadPromotionReport(path string) ([]PromotionEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []PromotionEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// FormatPromotion returns a human-readable summary of promotion entries.
func FormatPromotion(entries []PromotionEntry) string {
	if len(entries) == 0 {
		return "no promotion assessments available\n"
	}
	out := ""
	for _, e := range entries {
		status := "SAFE"
		if !e.Safe {
			status = "BLOCKED"
		}
		out += fmt.Sprintf("[%s] %s (%s -> %s): %s\n", status, e.Service, e.FromStage, e.ToStage, e.Reason)
	}
	return out
}
