package drift

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// IgnoreRule defines a rule to suppress specific drift findings.
type IgnoreRule struct {
	Service string `json:"service"` // empty = match all
	Key     string `json:"key"`     // exact or prefix with "*" wildcard
	Reason  string `json:"reason,omitempty"`
}

// IgnoreList holds a set of ignore rules.
type IgnoreList struct {
	Rules []IgnoreRule `json:"rules"`
}

// LoadIgnoreList reads an ignore list from a JSON file.
func LoadIgnoreList(path string) (*IgnoreList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &IgnoreList{}, nil
		}
		return nil, err
	}
	var il IgnoreList
	if err := json.Unmarshal(data, &il); err != nil {
		return nil, err
	}
	return &il, nil
}

// SaveIgnoreList writes an ignore list to a JSON file.
func SaveIgnoreList(path string, il *IgnoreList) error {
	data, err := json.MarshalIndent(il, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// AddIgnoreRule appends a rule to the ignore list file.
func AddIgnoreRule(path, service, key, reason string) error {
	il, err := LoadIgnoreList(path)
	if err != nil {
		return err
	}
	il.Rules = append(il.Rules, IgnoreRule{Service: service, Key: key, Reason: reason})
	return SaveIgnoreList(path, il)
}

// ApplyIgnoreList filters out diffs that match any ignore rule.
func ApplyIgnoreList(results []CompareResult, il *IgnoreList) []CompareResult {
	if il == nil || len(il.Rules) == 0 {
		return results
	}
	out := make([]CompareResult, 0, len(results))
	for _, r := range results {
		r.Diffs = filterIgnoredDiffs(r.Service, r.Diffs, il)
		out = append(out, r)
	}
	return out
}

func filterIgnoredDiffs(service string, diffs []DiffEntry, il *IgnoreList) []DiffEntry {
	var kept []DiffEntry
	for _, d := range diffs {
		if !isIgnored(service, d.Key, il) {
			kept = append(kept, d)
		}
	}
	return kept
}

func isIgnored(service, key string, il *IgnoreList) bool {
	for _, rule := range il.Rules {
		if rule.Service != "" && rule.Service != service {
			continue
		}
		pattern := rule.Key
		if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(key, strings.TrimSuffix(pattern, "*")) {
				return true
			}
		} else if pattern == key {
			return true
		}
	}
	return false
}
