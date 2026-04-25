package drift

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveCorrelation writes a CorrelationReport to a JSON file.
func SaveCorrelation(path string, report CorrelationReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal correlation: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write correlation: %w", err)
	}
	return nil
}

// LoadCorrelation reads a CorrelationReport from a JSON file.
func LoadCorrelation(path string) (CorrelationReport, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return CorrelationReport{}, nil
		}
		return CorrelationReport{}, fmt.Errorf("read correlation: %w", err)
	}
	var report CorrelationReport
	if err := json.Unmarshal(data, &report); err != nil {
		return CorrelationReport{}, fmt.Errorf("unmarshal correlation: %w", err)
	}
	return report, nil
}

// TopCorrelated returns the top N pairs by shared key count.
func TopCorrelated(report CorrelationReport, n int) []CorrelationEntry {
	pairs := make([]CorrelationEntry, len(report.Pairs))
	copy(pairs, report.Pairs)
	// simple insertion sort descending by SharedCount
	for i := 1; i < len(pairs); i++ {
		for j := i; j > 0 && pairs[j].SharedCount > pairs[j-1].SharedCount; j-- {
			pairs[j], pairs[j-1] = pairs[j-1], pairs[j]
		}
	}
	if n > len(pairs) {
		n = len(pairs)
	}
	return pairs[:n]
}
