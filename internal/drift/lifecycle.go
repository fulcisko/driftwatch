package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LifecycleStage represents the current stage of a service in its drift lifecycle.
type LifecycleStage string

const (
	StageActive     LifecycleStage = "active"
	StageWatched    LifecycleStage = "watched"
	StageDeprecated LifecycleStage = "deprecated"
	StageRetired    LifecycleStage = "retired"
)

// LifecycleEntry records the lifecycle stage for a service.
type LifecycleEntry struct {
	Service   string         `json:"service"`
	Stage     LifecycleStage `json:"stage"`
	UpdatedAt time.Time      `json:"updated_at"`
	Note      string         `json:"note,omitempty"`
}

// LifecycleStore holds all lifecycle entries.
type LifecycleStore struct {
	Entries []LifecycleEntry `json:"entries"`
}

// SetLifecycle adds or updates the lifecycle stage for a service.
func SetLifecycle(path, service string, stage LifecycleStage, note string) error {
	if service == "" {
		return fmt.Errorf("service name is required")
	}
	if stage == "" {
		return fmt.Errorf("lifecycle stage is required")
	}
	store, _ := LoadLifecycle(path)
	for i, e := range store.Entries {
		if e.Service == service {
			store.Entries[i].Stage = stage
			store.Entries[i].UpdatedAt = time.Now().UTC()
			store.Entries[i].Note = note
			return saveLifecycle(path, store)
		}
	}
	store.Entries = append(store.Entries, LifecycleEntry{
		Service:   service,
		Stage:     stage,
		UpdatedAt: time.Now().UTC(),
		Note:      note,
	})
	return saveLifecycle(path, store)
}

// LoadLifecycle reads the lifecycle store from disk.
func LoadLifecycle(path string) (LifecycleStore, error) {
	var store LifecycleStore
	data, err := os.ReadFile(path)
	if err != nil {
		return store, nil
	}
	err = json.Unmarshal(data, &store)
	return store, err
}

// FilterByStage returns entries matching the given stage.
func FilterByStage(store LifecycleStore, stage LifecycleStage) []LifecycleEntry {
	var out []LifecycleEntry
	for _, e := range store.Entries {
		if e.Stage == stage {
			out = append(out, e)
		}
	}
	return out
}

func saveLifecycle(path string, store LifecycleStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
