package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// WatchEntry represents a service being actively watched for drift.
type WatchEntry struct {
	Service   string    `json:"service"`
	AddedAt   time.Time `json:"added_at"`
	Threshold int       `json:"threshold"` // min drift count to alert
}

// Watchlist holds all watched services.
type Watchlist struct {
	Entries []WatchEntry `json:"entries"`
}

// AddToWatchlist appends a service entry to the watchlist file.
func AddToWatchlist(path, service string, threshold int) error {
	wl, err := LoadWatchlist(path)
	if err != nil {
		wl = &Watchlist{}
	}
	for _, e := range wl.Entries {
		if e.Service == service {
			return fmt.Errorf("service %q already in watchlist", service)
		}
	}
	wl.Entries = append(wl.Entries, WatchEntry{
		Service:   service,
		AddedAt:   time.Now().UTC(),
		Threshold: threshold,
	})
	return saveWatchlist(path, wl)
}

// RemoveFromWatchlist removes a service entry by name.
func RemoveFromWatchlist(path, service string) error {
	wl, err := LoadWatchlist(path)
	if err != nil {
		return err
	}
	filtered := wl.Entries[:0]
	for _, e := range wl.Entries {
		if e.Service != service {
			filtered = append(filtered, e)
		}
	}
	if len(filtered) == len(wl.Entries) {
		return fmt.Errorf("service %q not found in watchlist", service)
	}
	wl.Entries = filtered
	return saveWatchlist(path, wl)
}

// LoadWatchlist reads the watchlist from disk.
func LoadWatchlist(path string) (*Watchlist, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var wl Watchlist
	if err := json.Unmarshal(data, &wl); err != nil {
		return nil, err
	}
	return &wl, nil
}

// MatchWatchlist filters results to only those in the watchlist that exceed threshold.
func MatchWatchlist(wl *Watchlist, results []CompareResult) []CompareResult {
	index := make(map[string]WatchEntry, len(wl.Entries))
	for _, e := range wl.Entries {
		index[e.Service] = e
	}
	var matched []CompareResult
	for _, r := range results {
		if e, ok := index[r.Service]; ok {
			if len(r.Diffs) >= e.Threshold {
				matched = append(matched, r)
			}
		}
	}
	return matched
}

func saveWatchlist(path string, wl *Watchlist) error {
	data, err := json.MarshalIndent(wl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
