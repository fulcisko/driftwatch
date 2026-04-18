package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of drift results.
type Snapshot struct {
	Timestamp time.Time       `json:"timestamp"`
	Label     string          `json:"label"`
	Results   []CompareResult `json:"results"`
}

// SaveSnapshot writes a labeled snapshot to the given file path.
func SaveSnapshot(path, label string, results []CompareResult) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot write: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from the given file path.
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, fmt.Errorf("snapshot not found: %s", path)
		}
		return Snapshot{}, fmt.Errorf("snapshot read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot unmarshal: %w", err)
	}
	return snap, nil
}

// DiffSnapshot compares current results against a saved snapshot and returns
// services whose drift status changed.
func DiffSnapshot(snap Snapshot, current []CompareResult) []string {
	prev := make(map[string]bool)
	for _, r := range snap.Results {
		prev[r.Service] = len(r.Diffs) > 0
	}
	var changed []string
	for _, r := range current {
		hasDrift := len(r.Diffs) > 0
		if was, ok := prev[r.Service]; !ok || was != hasDrift {
			changed = append(changed, r.Service)
		}
	}
	return changed
}
