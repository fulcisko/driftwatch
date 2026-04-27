package drift

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Fingerprint represents a stable hash of a service's drift state.
type Fingerprint struct {
	Service     string `json:"service"`
	Hash        string `json:"hash"`
	DriftKeys   []string `json:"drift_keys"`
	DriftCount  int    `json:"drift_count"`
}

// FingerprintStore maps service names to their fingerprints.
type FingerprintStore map[string]Fingerprint

// BuildFingerprint computes a stable fingerprint for a single CompareResult.
func BuildFingerprint(r CompareResult) Fingerprint {
	keys := make([]string, 0, len(r.Diffs))
	for _, d := range r.Diffs {
		keys = append(keys, d.Key)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s:", k)
	}
	hash := hex.EncodeToString(h.Sum(nil))[:16]

	return Fingerprint{
		Service:    r.Service,
		Hash:       hash,
		DriftKeys:  keys,
		DriftCount: len(r.Diffs),
	}
}

// BuildFingerprintStore builds a store for all drifted services.
func BuildFingerprintStore(results []CompareResult) FingerprintStore {
	store := make(FingerprintStore)
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		store[r.Service] = BuildFingerprint(r)
	}
	return store
}

// SaveFingerprintStore writes the store to a JSON file.
func SaveFingerprintStore(path string, store FingerprintStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal fingerprint store: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadFingerprintStore reads a fingerprint store from a JSON file.
func LoadFingerprintStore(path string) (FingerprintStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(FingerprintStore), nil
		}
		return nil, fmt.Errorf("read fingerprint store: %w", err)
	}
	var store FingerprintStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("unmarshal fingerprint store: %w", err)
	}
	return store, nil
}

// DiffFingerprintStore returns services whose fingerprint hash has changed.
func DiffFingerprintStore(old, current FingerprintStore) []string {
	changed := []string{}
	for svc, fp := range current {
		prev, ok := old[svc]
		if !ok || prev.Hash != fp.Hash {
			changed = append(changed, svc)
		}
	}
	sort.Strings(changed)
	return changed
}

// FormatFingerprintStore returns a human-readable summary.
func FormatFingerprintStore(store FingerprintStore) string {
	if len(store) == 0 {
		return "no drifted services fingerprinted\n"
	}
	services := make([]string, 0, len(store))
	for svc := range store {
		services = append(services, svc)
	}
	sort.Strings(services)

	var sb strings.Builder
	for _, svc := range services {
		fp := store[svc]
		fmt.Fprintf(&sb, "%-30s hash=%-16s drifts=%d keys=[%s]\n",
			fp.Service, fp.Hash, fp.DriftCount, strings.Join(fp.DriftKeys, ","))
	}
	return sb.String()
}
