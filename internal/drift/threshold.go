package drift

import (
	"encoding/json"
	"fmt"
	"os"
)

// ThresholdRule defines acceptable drift limits for a service.
type ThresholdRule struct {
	Service      string `json:"service"`
	MaxDrifts    int    `json:"max_drifts"`
	MinSeverity  string `json:"min_severity"`
}

// ThresholdList holds all threshold rules.
type ThresholdList struct {
	Rules []ThresholdRule `json:"rules"`
}

// ThresholdViolation describes a service that exceeded its threshold.
type ThresholdViolation struct {
	Service  string
	Rule     ThresholdRule
	Drifts   int
	Severity SeverityLevel
}

// LoadThresholds reads threshold rules from a JSON file.
func LoadThresholds(path string) (ThresholdList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ThresholdList{}, nil
		}
		return ThresholdList{}, err
	}
	var tl ThresholdList
	if err := json.Unmarshal(data, &tl); err != nil {
		return ThresholdList{}, fmt.Errorf("parse thresholds: %w", err)
	}
	return tl, nil
}

// CheckThresholds evaluates results against threshold rules and returns violations.
func CheckThresholds(results []CompareResult, tl ThresholdList) []ThresholdViolation {
	ruleMap := make(map[string]ThresholdRule)
	for _, r := range tl.Rules {
		ruleMap[r.Service] = r
	}

	var violations []ThresholdViolation
	for _, res := range results {
		rule, ok := ruleMap[res.Service]
		if !ok {
			continue
		}
		driftCount := 0
		for _, d := range res.Diffs {
			if d.Status != "match" {
				driftCount++
			}
		}
		sev := MaxSeverity(res.Diffs)
		minSev := parseSeverity(rule.MinSeverity)
		if driftCount > rule.MaxDrifts || sev >= minSev {
			violations = append(violations, ThresholdViolation{
				Service:  res.Service,
				Rule:     rule,
				Drifts:   driftCount,
				Severity: sev,
			})
		}
	}
	return violations
}

func parseSeverity(s string) SeverityLevel {
	switch s {
	case "high":
		return SeverityHigh
	case "medium":
		return SeverityMedium
	case "low":
		return SeverityLow
	default:
		return SeverityNone
	}
}
