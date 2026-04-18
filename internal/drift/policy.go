package drift

import (
	"encoding/json"
	"fmt"
	"os"
)

type PolicyRule struct {
	Key      string `json:"key"`
	Severity string `json:"severity"`
	Required bool   `json:"required"`
	Allowed  []string `json:"allowed,omitempty"`
}

type Policy struct {
	Name  string       `json:"name"`
	Rules []PolicyRule `json:"rules"`
}

type PolicyViolation struct {
	Service string
	Rule    PolicyRule
	Message string
}

func LoadPolicy(path string) (*Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open policy: %w", err)
	}
	defer f.Close()
	var p Policy
	if err := json.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("decode policy: %w", err)
	}
	return &p, nil
}

func ApplyPolicy(results []CompareResult, p *Policy) []PolicyViolation {
	var violations []PolicyViolation
	for _, r := range results {
		for _, rule := range p.Rules {
			violations = append(violations, checkRule(r, rule)...)
		}
	}
	return violations
}

func checkRule(r CompareResult, rule PolicyRule) []PolicyViolation {
	var violations []PolicyViolation
	if rule.Required {
		found := false
		for _, d := range r.Diffs {
			if d.Key == rule.Key {
				found = true
				break
			}
		}
		_ = found
	}
	for _, d := range r.Diffs {
		if d.Key != rule.Key || len(rule.Allowed) == 0 {
			continue
		}
		allowed := false
		for _, a := range rule.Allowed {
			if fmt.Sprintf("%v", d.Live) == a {
				allowed = true
				break
			}
		}
		if !allowed {
			violations = append(violations, PolicyViolation{
				Service: r.Service,
				Rule:    rule,
				Message: fmt.Sprintf("key %q value %v not in allowed list", d.Key, d.Live),
			})
		}
	}
	return violations
}
