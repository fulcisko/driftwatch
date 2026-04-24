package drift

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// DigestEntry records a hash of drift results at a point in time.
type DigestEntry struct {
	Service   string    `json:"service"`
	Hash      string    `json:"hash"`
	DiffCount int       `json:"diff_count"`
	CreatedAt time.Time `json:"created_at"`
}

// ComputeDigest produces a stable SHA-256 hash over the diffs for a single service.
func ComputeDigest(result CompareResult) string {
	keys := make([]string, 0, len(result.Diffs))
	for k := range result.Diffs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%v;", k, result.Diffs[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// SaveDigests writes a slice of DigestEntry values to the given path as JSON.
func SaveDigests(path string, entries []DigestEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal digests: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadDigests reads DigestEntry values from the given path.
// Returns an empty slice if the file does not exist.
func LoadDigests(path string) ([]DigestEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []DigestEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read digests: %w", err)
	}
	var entries []DigestEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse digests: %w", err)
	}
	return entries, nil
}

// BuildDigests creates a DigestEntry for each result and returns the slice.
func BuildDigests(results []CompareResult) []DigestEntry {
	entries := make([]DigestEntry, 0, len(results))
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		entries = append(entries, DigestEntry{
			Service:   r.Service,
			Hash:      ComputeDigest(r),
			DiffCount: len(r.Diffs),
			CreatedAt: time.Now().UTC(),
		})
	}
	return entries
}

// DigestsChanged returns the services whose digest has changed compared to
// a previously saved set of entries.
func DigestsChanged(previous, current []DigestEntry) []string {
	prev := make(map[string]string, len(previous))
	for _, e := range previous {
		prev[e.Service] = e.Hash
	}
	var changed []string
	for _, e := range current {
		if prev[e.Service] != e.Hash {
			changed = append(changed, e.Service)
		}
	}
	return changed
}
