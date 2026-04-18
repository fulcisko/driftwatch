package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved snapshot of drift results.
type Baseline struct {
	CreatedAt time.Time        `json:"created_at"`
	Results   []CompareResult  `json:"results"`
}

// SaveBaseline writes the current results to a JSON file at path.
func SaveBaseline(path string, results []CompareResult) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal baseline: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write baseline: %w", err)
	}
	return nil
}

// LoadBaseline reads a previously saved baseline from path.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read baseline: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("unmarshal baseline: %w", err)
	}
	return &b, nil
}

// DiffBaseline compares current results against a baseline and returns
// only results whose drift state has changed since the baseline.
func DiffBaseline(baseline *Baseline, current []CompareResult) []CompareResult {
	index := make(map[string]CompareResult, len(baseline.Results))
	for _, r := range baseline.Results {
		index[r.Service] = r
	}
	var changed []CompareResult
	for _, cur := range current {
		prev, ok := index[cur.Service]
		if !ok || driftCount(cur) != driftCount(prev) {
			changed = append(changed, cur)
		}
	}
	return changed
}
