package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Profile represents a named drift configuration profile that groups
// filter options, severity thresholds, and ignored keys together.
type Profile struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IgnoreKeys  []string  `json:"ignore_keys"`
	MinSeverity string    `json:"min_severity"`
	ServicePrefix string  `json:"service_prefix,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type profileStore struct {
	Profiles []Profile `json:"profiles"`
}

// SaveProfile adds or updates a profile in the given file.
func SaveProfile(path string, p Profile) error {
	if p.Name == "" {
		return fmt.Errorf("profile name is required")
	}
	store, _ := loadProfiles(path)
	now := time.Now().UTC()
	for i, existing := range store.Profiles {
		if existing.Name == p.Name {
			p.CreatedAt = existing.CreatedAt
			p.UpdatedAt = now
			store.Profiles[i] = p
			return writeProfiles(path, store)
		}
	}
	p.CreatedAt = now
	p.UpdatedAt = now
	store.Profiles = append(store.Profiles, p)
	return writeProfiles(path, store)
}

// LoadProfiles returns all profiles from the given file.
func LoadProfiles(path string) ([]Profile, error) {
	store, err := loadProfiles(path)
	if err != nil {
		return nil, err
	}
	return store.Profiles, nil
}

// GetProfile returns a single profile by name.
func GetProfile(path, name string) (Profile, bool) {
	store, _ := loadProfiles(path)
	for _, p := range store.Profiles {
		if p.Name == name {
			return p, true
		}
	}
	return Profile{}, false
}

// RemoveProfile deletes a profile by name.
func RemoveProfile(path, name string) error {
	store, err := loadProfiles(path)
	if err != nil {
		return err
	}
	updated := store.Profiles[:0]
	found := false
	for _, p := range store.Profiles {
		if p.Name == name {
			found = true
			continue
		}
		updated = append(updated, p)
	}
	if !found {
		return fmt.Errorf("profile %q not found", name)
	}
	store.Profiles = updated
	return writeProfiles(path, store)
}

func loadProfiles(path string) (profileStore, error) {
	var store profileStore
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return store, err
	}
	err = json.Unmarshal(data, &store)
	return store, err
}

func writeProfiles(path string, store profileStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
