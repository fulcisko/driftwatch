package drift

import (
	"encoding/json"
	"os"
	"time"
)

// ChangelogEntry records a single drift event for a service.
type ChangelogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Service   string            `json:"service"`
	Diffs     []CompareResult   `json:"diffs"`
	Note      string            `json:"note,omitempty"`
}

// Changelog is an ordered list of entries.
type Changelog []ChangelogEntry

// AppendChangelog adds a new entry to the changelog file.
func AppendChangelog(path string, entry ChangelogEntry) error {
	cl, err := LoadChangelog(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	cl = append(cl, entry)
	data, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadChangelog reads all changelog entries from disk.
func LoadChangelog(path string) (Changelog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Changelog{}, nil
		}
		return nil, err
	}
	var cl Changelog
	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, err
	}
	return cl, nil
}

// FilterChangelog returns entries matching the given service name.
// Pass an empty string to return all entries.
func FilterChangelog(cl Changelog, service string) Changelog {
	if service == "" {
		return cl
	}
	var out Changelog
	for _, e := range cl {
		if e.Service == service {
			out = append(out, e)
		}
	}
	return out
}
