package drift

import (
	"fmt"
	"sort"
	"strings"
)

// CorrelationEntry represents a pair of services with shared drift keys.
type CorrelationEntry struct {
	ServiceA    string   `json:"service_a"`
	ServiceB    string   `json:"service_b"`
	SharedKeys  []string `json:"shared_keys"`
	SharedCount int      `json:"shared_count"`
}

// CorrelationReport holds all correlated service pairs.
type CorrelationReport struct {
	Pairs []CorrelationEntry `json:"pairs"`
}

// BuildCorrelation finds services that share drifted config keys.
func BuildCorrelation(results []CompareResult) CorrelationReport {
	// Map service -> set of drifted keys
	serviceKeys := make(map[string]map[string]struct{})
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		keys := make(map[string]struct{})
		for _, d := range r.Diffs {
			keys[d.Key] = struct{}{}
		}
		serviceKeys[r.Service] = keys
	}

	services := make([]string, 0, len(serviceKeys))
	for svc := range serviceKeys {
		services = append(services, svc)
	}
	sort.Strings(services)

	var pairs []CorrelationEntry
	for i := 0; i < len(services); i++ {
		for j := i + 1; j < len(services); j++ {
			shared := sharedKeys(serviceKeys[services[i]], serviceKeys[services[j]])
			if len(shared) > 0 {
				pairs = append(pairs, CorrelationEntry{
					ServiceA:    services[i],
					ServiceB:    services[j],
					SharedKeys:  shared,
					SharedCount: len(shared),
				})
			}
		}
	}
	return CorrelationReport{Pairs: pairs}
}

func sharedKeys(a, b map[string]struct{}) []string {
	var shared []string
	for k := range a {
		if _, ok := b[k]; ok {
			shared = append(shared, k)
		}
	}
	sort.Strings(shared)
	return shared
}

// FormatCorrelation returns a human-readable summary of the correlation report.
func FormatCorrelation(report CorrelationReport) string {
	if len(report.Pairs) == 0 {
		return "No correlated drift found across services.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Correlated drift pairs: %d\n\n", len(report.Pairs)))
	for _, p := range report.Pairs {
		sb.WriteString(fmt.Sprintf("  %s <-> %s (%d shared key(s)): %s\n",
			p.ServiceA, p.ServiceB, p.SharedCount, strings.Join(p.SharedKeys, ", ")))
	}
	return sb.String()
}
