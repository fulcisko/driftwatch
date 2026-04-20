package drift

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

type Tag struct {
	Name     string   `json:"name"`
	Services []string `json:"services"`
}

type TagStore struct {
	Tags []Tag `json:"tags"`
}

// AddTag adds a service to the named tag, creating the tag if it does not exist.
// If the service is already associated with the tag, it is a no-op.
func AddTag(path, tagName, service string) error {
	store, _ := LoadTags(path)
	for i, t := range store.Tags {
		if t.Name == tagName {
			for _, s := range t.Services {
				if s == service {
					return nil
				}
			}
			store.Tags[i].Services = append(store.Tags[i].Services, service)
			sort.Strings(store.Tags[i].Services)
			return saveTags(path, store)
		}
	}
	store.Tags = append(store.Tags, Tag{Name: tagName, Services: []string{service}})
	return saveTags(path, store)
}

// RemoveTag removes a service from the named tag.
// Returns an error if the tag does not exist.
func RemoveTag(path, tagName, service string) error {
	store, err := LoadTags(path)
	if err != nil {
		return err
	}
	for i, t := range store.Tags {
		if t.Name == tagName {
			filtered := []string{}
			for _, s := range t.Services {
				if s != service {
					filtered = append(filtered, s)
				}
			}
			store.Tags[i].Services = filtered
			return saveTags(path, store)
		}
	}
	return errors.New("tag not found")
}

// LoadTags reads and parses the tag store from the given file path.
// If the file does not exist, an empty TagStore is returned without error.
func LoadTags(path string) (TagStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return TagStore{}, nil
		}
		return TagStore{}, err
	}
	var store TagStore
	if err := json.Unmarshal(data, &store); err != nil {
		return TagStore{}, err
	}
	return store, nil
}

// FilterByTag returns the list of services associated with the given tag name.
// Returns nil if the tag does not exist.
func FilterByTag(store TagStore, tagName string) []string {
	for _, t := range store.Tags {
		if t.Name == tagName {
			return t.Services
		}
	}
	return nil
}

func saveTags(path string, store TagStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
