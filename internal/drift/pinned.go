package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PinnedKey represents a config key whose value is pinned to an expected value.
type PinnedKey struct {
	Service   string    `json:"service"`
	Key       string    `json:"key"`
	Expected  string    `json:"expected"`
	PinnedAt  time.Time `json:"pinned_at"`
	Comment   string    `json:"comment,omitempty"`
}

// PinList holds all pinned keys.
type PinList struct {
	Pins []PinnedKey `json:"pins"`
}

// AddPin adds or updates a pin for a service/key pair.
func AddPin(path, service, key, expected, comment string) error {
	list, _ := LoadPins(path)
	for i, p := range list.Pins {
		if p.Service == service && p.Key == key {
			list.Pins[i].Expected = expected
			list.Pins[i].PinnedAt = time.Now().UTC()
			list.Pins[i].Comment = comment
			return savePins(path, list)
		}
	}
	list.Pins = append(list.Pins, PinnedKey{
		Service:  service,
		Key:      key,
		Expected: expected,
		PinnedAt: time.Now().UTC(),
		Comment:  comment,
	})
	return savePins(path, list)
}

// RemovePin removes a pin for a service/key pair.
func RemovePin(path, service, key string) error {
	list, err := LoadPins(path)
	if err != nil {
		return err
	}
	filtered := list.Pins[:0]
	for _, p := range list.Pins {
		if !(p.Service == service && p.Key == key) {
			filtered = append(filtered, p)
		}
	}
	if len(filtered) == len(list.Pins) {
		return fmt.Errorf("pin not found: %s/%s", service, key)
	}
	list.Pins = filtered
	return savePins(path, list)
}

// LoadPins loads the pin list from disk.
func LoadPins(path string) (PinList, error) {
	var list PinList
	data, err := os.ReadFile(path)
	if err != nil {
		return list, nil
	}
	err = json.Unmarshal(data, &list)
	if err != nil {
		return list, fmt.Errorf("failed to parse pin file %q: %w", path, err)
	}
	return list, nil
}

// ApplyPins filters out diffs where the live value matches the pinned expected value.
func ApplyPins(results []CompareResult, list PinList) []CompareResult {
	pinMap := map[string]map[string]string{}
	for _, p := range list.Pins {
		if pinMap[p.Service] == nil {
			pinMap[p.Service] = map[string]string{}
		}
		pinMap[p.Service][p.Key] = p.Expected
	}
	var out []CompareResult
	for _, r := range results {
		filtered := r.Diffs[:0]
		for _, d := range r.Diffs {
			if exp, ok := pinMap[r.Service][d.Key]; ok && fmt.Sprintf("%v", d.LiveValue) == exp {
				continue
			}
			filtered = append(filtered, d)
		}
		r.Diffs = filtered
		out = append(out, r)
	}
	return out
}

func savePins(path string, list PinList) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
