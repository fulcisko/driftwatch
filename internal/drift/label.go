package drift

import (
	"encoding/json"
	"errors"
	"os"
)

type LabelMap map[string]map[string]string // service -> key -> value

func AddLabel(path, service, key, value string) error {
	if service == "" || key == "" {
		return errors.New("service and key are required")
	}
	labels, _ := LoadLabels(path)
	if labels == nil {
		labels = make(LabelMap)
	}
	if labels[service] == nil {
		labels[service] = make(map[string]string)
	}
	labels[service][key] = value
	return saveLabels(path, labels)
}

func RemoveLabel(path, service, key string) error {
	labels, err := LoadLabels(path)
	if err != nil {
		return err
	}
	if _, ok := labels[service]; !ok {
		return errors.New("service not found")
	}
	delete(labels[service], key)
	if len(labels[service]) == 0 {
		delete(labels, service)
	}
	return saveLabels(path, labels)
}

func LoadLabels(path string) (LabelMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(LabelMap), nil
		}
		return nil, err
	}
	var labels LabelMap
	if err := json.Unmarshal(data, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

func FilterByLabel(results []CompareResult, labels LabelMap, key, value string) []CompareResult {
	var out []CompareResult
	for _, r := range results {
		if sv, ok := labels[r.Service]; ok {
			if sv[key] == value {
				out = append(out, r)
			}
		}
	}
	return out
}

func saveLabels(path string, labels LabelMap) error {
	data, err := json.MarshalIndent(labels, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
