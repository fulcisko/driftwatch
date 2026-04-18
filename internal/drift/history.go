package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a snapshot of drift results at a point in time.
type HistoryEntry struct {
	Timestamp time.Time       `json:"timestamp"`
	Results   []CompareResult `json:"results"`
}

// AppendHistory appends the current results to a history file.
func AppendHistory(path string, results []CompareResult) error {
	entries, err := LoadHistory(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}
	entries = append(entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Results:   results,
	})
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create history file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

// LoadHistory reads all history entries from the given file.
func LoadHistory(path string) ([]HistoryEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []HistoryEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode history: %w", err)
	}
	return entries, nil
}

// LatestHistory returns the most recent history entry, if any.
func LatestHistory(path string) (*HistoryEntry, error) {
	entries, err := LoadHistory(path)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	e := entries[len(entries)-1]
	return &e, nil
}
