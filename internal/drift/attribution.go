package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Attribution records who is responsible for a drift event.
type Attribution struct {
	Service   string    `json:"service"`
	Key       string    `json:"key"`
	Owner     string    `json:"owner"`
	Team      string    `json:"team"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// AttributionStore holds all attribution records.
type AttributionStore struct {
	Entries []Attribution `json:"entries"`
}

// AddAttribution appends an attribution record to the store at path.
func AddAttribution(path, service, key, owner, team, reason string) error {
	if service == "" || key == "" || owner == "" {
		return fmt.Errorf("service, key, and owner are required")
	}
	store, _ := LoadAttributions(path)
	store.Entries = append(store.Entries, Attribution{
		Service:   service,
		Key:       key,
		Owner:     owner,
		Team:      team,
		Reason:    reason,
		Timestamp: time.Now().UTC(),
	})
	return saveAttributions(path, store)
}

// LoadAttributions reads the attribution store from path.
func LoadAttributions(path string) (AttributionStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AttributionStore{}, nil
		}
		return AttributionStore{}, err
	}
	var store AttributionStore
	if err := json.Unmarshal(data, &store); err != nil {
		return AttributionStore{}, err
	}
	return store, nil
}

// FilterAttributions returns entries matching service (empty matches all).
func FilterAttributions(store AttributionStore, service string) []Attribution {
	if service == "" {
		return store.Entries
	}
	var out []Attribution
	for _, e := range store.Entries {
		if e.Service == service {
			out = append(out, e)
		}
	}
	return out
}

func saveAttributions(path string, store AttributionStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
