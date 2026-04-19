package drift

import (
	"encoding/json"
	"os"
	"time"
)

// SuppressRule silences drift alerts for a service/key until a given expiry.
type SuppressRule struct {
	Service   string    `json:"service"`
	Key       string    `json:"key"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SuppressList struct {
	Rules []SuppressRule `json:"rules"`
}

func LoadSuppressList(path string) (SuppressList, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return SuppressList{}, nil
	}
	if err != nil {
		return SuppressList{}, err
	}
	var sl SuppressList
	if err := json.Unmarshal(data, &sl); err != nil {
		return SuppressList{}, err
	}
	return sl, nil
}

func SaveSuppressList(path string, sl SuppressList) error {
	data, err := json.MarshalIndent(sl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddSuppressRule(path string, rule SuppressRule) error {
	sl, err := LoadSuppressList(path)
	if err != nil {
		return err
	}
	sl.Rules = append(sl.Rules, rule)
	return SaveSuppressList(path, sl)
}

// ApplySuppress filters out CompareResults whose diffs are currently suppressed.
func ApplySuppress(results []CompareResult, sl SuppressList) []CompareResult {
	now := time.Now()
	out := make([]CompareResult, 0, len(results))
	for _, r := range results {
		filtered := r.Diffs[:0]
		for _, d := range r.Diffs {
			if !isSuppressed(r.Service, d.Key, sl, now) {
				filtered = append(filtered, d)
			}
		}
		r.Diffs = filtered
		out = append(out, r)
	}
	return out
}

func isSuppressed(service, key string, sl SuppressList, now time.Time) bool {
	for _, rule := range sl.Rules {
		if rule.ExpiresAt.Before(now) {
			continue
		}
		if rule.Service == service && (rule.Key == key || rule.Key == "*") {
			return true
		}
	}
	return false
}
