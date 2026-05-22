package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// QuotaRule defines a maximum number of allowed drifted keys for a service.
type QuotaRule struct {
	Service  string    `json:"service"`
	MaxDrift int       `json:"max_drift"`
	UpdatedAt time.Time `json:"updated_at"`
}

// QuotaViolation describes a service that has exceeded its drift quota.
type QuotaViolation struct {
	Service   string
	Allowed   int
	Actual    int
}

// QuotaList holds all quota rules.
type QuotaList struct {
	Rules []QuotaRule `json:"rules"`
}

// AddQuotaRule adds or updates a quota rule for a service.
func AddQuotaRule(path, service string, maxDrift int) error {
	if service == "" {
		return fmt.Errorf("service name is required")
	}
	if maxDrift < 0 {
		return fmt.Errorf("max_drift must be non-negative")
	}
	ql, _ := LoadQuotas(path)
	for i, r := range ql.Rules {
		if r.Service == service {
			ql.Rules[i].MaxDrift = maxDrift
			ql.Rules[i].UpdatedAt = time.Now().UTC()
			return saveQuotas(path, ql)
		}
	}
	ql.Rules = append(ql.Rules, QuotaRule{
		Service:   service,
		MaxDrift:  maxDrift,
		UpdatedAt: time.Now().UTC(),
	})
	return saveQuotas(path, ql)
}

// LoadQuotas reads quota rules from disk.
func LoadQuotas(path string) (QuotaList, error) {
	var ql QuotaList
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ql, nil
		}
		return ql, err
	}
	err = json.Unmarshal(data, &ql)
	return ql, err
}

// CheckQuotas returns violations for any service that exceeds its quota.
func CheckQuotas(path string, results []CompareResult) ([]QuotaViolation, error) {
	ql, err := LoadQuotas(path)
	if err != nil {
		return nil, err
	}
	ruleMap := make(map[string]int)
	for _, r := range ql.Rules {
		ruleMap[r.Service] = r.MaxDrift
	}
	var violations []QuotaViolation
	for _, res := range results {
		max, ok := ruleMap[res.Service]
		if !ok {
			continue
		}
		count := len(res.Diffs)
		if count > max {
			violations = append(violations, QuotaViolation{
				Service: res.Service,
				Allowed: max,
				Actual:  count,
			})
		}
	}
	return violations, nil
}

func saveQuotas(path string, ql QuotaList) error {
	data, err := json.MarshalIndent(ql, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
