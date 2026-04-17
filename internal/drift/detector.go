package drift

import (
	"fmt"
	"reflect"
)

// DriftResult holds the comparison result between a deployed config and a manifest.
type DriftResult struct {
	ServiceName string
	HasDrift    bool
	Diffs       []FieldDiff
}

// FieldDiff describes a single field that differs.
type FieldDiff struct {
	Field    string
	Expected interface{}
	Actual   interface{}
}

// Detector compares deployed state against source manifests.
type Detector struct{}

// NewDetector creates a new Detector.
func NewDetector() *Detector {
	return &Detector{}
}

// Compare checks deployed against expected and returns a DriftResult.
func (d *Detector) Compare(name string, expected, deployed map[string]interface{}) DriftResult {
	result := DriftResult{ServiceName: name}

	for key, expVal := range expected {
		actVal, ok := deployed[key]
		if !ok {
			result.Diffs = append(result.Diffs, FieldDiff{
				Field:    key,
				Expected: expVal,
				Actual:   nil,
			})
			continue
		}
		if !reflect.DeepEqual(expVal, actVal) {
			result.Diffs = append(result.Diffs, FieldDiff{
				Field:    key,
				Expected: expVal,
				Actual:   actVal,
			})
		}
	}

	for key, actVal := range deployed {
		if _, ok := expected[key]; !ok {
			result.Diffs = append(result.Diffs, FieldDiff{
				Field:    fmt.Sprintf("%s (unexpected)", key),
				Expected: nil,
				Actual:   actVal,
			})
		}
	}

	result.HasDrift = len(result.Diffs) > 0
	return result
}
