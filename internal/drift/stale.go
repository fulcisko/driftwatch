package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// StaleEntry represents a service whose drift has not been re-checked within a threshold.
type StaleEntry struct {
	Service    string    `json:"service"`
	LastSeen   time.Time `json:"last_seen"`
	StaleSince time.Duration `json:"stale_since_seconds"`
}

// FindStaleServices returns services whose last drift check is older than maxAge.
func FindStaleServices(historyPath string, maxAge time.Duration) ([]StaleEntry, error) {
	records, err := LoadHistory(historyPath)
	if err != nil {
		return nil, fmt.Errorf("load history: %w", err)
	}

	// Build a map of service -> latest timestamp
	latest := map[string]time.Time{}
	for _, r := range records {
		for _, res := range r.Results {
			if t, ok := latest[res.Service]; !ok || r.Timestamp.After(t) {
				latest[res.Service] = r.Timestamp
			}
		}
	}

	now := time.Now().UTC()
	var stale []StaleEntry
	for svc, ts := range latest {
		age := now.Sub(ts)
		if age > maxAge {
			stale = append(stale, StaleEntry{
				Service:    svc,
				LastSeen:   ts,
				StaleSince: age,
			})
		}
	}
	return stale, nil
}

// SaveStaleReport writes stale entries to a JSON file.
func SaveStaleReport(path string, entries []StaleEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal stale report: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadStaleReport reads a previously saved stale report from disk.
func LoadStaleReport(path string) ([]StaleEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read stale report: %w", err)
	}
	var entries []StaleEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse stale report: %w", err)
	}
	return entries, nil
}

// FormatStaleReport returns a human-readable summary of stale services.
func FormatStaleReport(entries []StaleEntry) string {
	if len(entries) == 0 {
		return "No stale services detected.\n"
	}
	out := fmt.Sprintf("Stale services (%d):\n", len(entries))
	for _, e := range entries {
		out += fmt.Sprintf("  %-30s last seen: %s (%.0fh ago)\n",
			e.Service,
			e.LastSeen.Format(time.RFC3339),
			e.StaleSince.Hours(),
		)
	}
	return out
}
