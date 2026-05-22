package drift

import (
	"fmt"
	"sort"
	"strings"
)

// CascadeRisk represents the propagated drift risk for a service,
// accounting for its own drift and the drift of its dependencies.
type CascadeRisk struct {
	Service      string   `json:"service"`
	DirectDrift  int      `json:"direct_drift"`
	CascadeDrift int      `json:"cascade_drift"`
	TotalRisk    int      `json:"total_risk"`
	RiskLevel    string   `json:"risk_level"`
	AffectedBy   []string `json:"affected_by,omitempty"`
}

// CascadeReport holds the full cascade risk assessment across all services.
type CascadeReport struct {
	Entries []CascadeRisk `json:"entries"`
}

// AssessCascade computes the cascading drift risk for each service by
// walking the dependency graph and summing drift counts from upstream services.
func AssessCascade(results []CompareResult, deps []Dependency) CascadeReport {
	// Build a map of direct drift counts per service.
	driftCounts := make(map[string]int)
	for _, r := range results {
		if !r.Clean {
			driftCounts[r.Service] = len(r.Diffs)
		}
	}

	// Build a dependents index: for each service, who depends on it?
	// We want to propagate drift upstream → downstream.
	// deps[i].DependsOn means Service depends on DependsOn.
	dependsOnMap := make(map[string][]string) // service -> list of its dependencies
	for _, d := range deps {
		dependsOnMap[d.Service] = append(dependsOnMap[d.Service], d.DependsOn)
	}

	var entries []CascadeRisk
	for _, r := range results {
		cascadeDrift := 0
		var affectedBy []string

		// Walk direct dependencies and accumulate their drift.
		visited := make(map[string]bool)
		queue := append([]string{}, dependsOnMap[r.Service]...)
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			if visited[cur] {
				continue
			}
			visited[cur] = true
			if count, ok := driftCounts[cur]; ok && count > 0 {
				cascadeDrift += count
				affectedBy = append(affectedBy, cur)
			}
			// Recurse into transitive dependencies.
			queue = append(queue, dependsOnMap[cur]...)
		}

		direct := driftCounts[r.Service]
		total := direct + cascadeDrift

		sort.Strings(affectedBy)
		entries = append(entries, CascadeRisk{
			Service:      r.Service,
			DirectDrift:  direct,
			CascadeDrift: cascadeDrift,
			TotalRisk:    total,
			RiskLevel:    cascadeRiskLevel(total),
			AffectedBy:   affectedBy,
		})
	}

	// Sort by total risk descending, then by service name.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].TotalRisk != entries[j].TotalRisk {
			return entries[i].TotalRisk > entries[j].TotalRisk
		}
		return entries[i].Service < entries[j].Service
	})

	return CascadeReport{Entries: entries}
}

// cascadeRiskLevel maps a total risk score to a human-readable level.
func cascadeRiskLevel(total int) string {
	switch {
	case total == 0:
		return "none"
	case total <= 3:
		return "low"
	case total <= 8:
		return "medium"
	default:
		return "high"
	}
}

// FormatCascade returns a human-readable table of cascade risk entries.
func FormatCascade(report CascadeReport) string {
	if len(report.Entries) == 0 {
		return "no cascade risk data available\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-24s %8s %8s %8s %8s\n",
		"SERVICE", "DIRECT", "CASCADE", "TOTAL", "LEVEL"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")

	for _, e := range report.Entries {
		sb.WriteString(fmt.Sprintf("%-24s %8d %8d %8d %8s\n",
			e.Service, e.DirectDrift, e.CascadeDrift, e.TotalRisk, e.RiskLevel))
		if len(e.AffectedBy) > 0 {
			sb.WriteString(fmt.Sprintf("  affected by: %s\n", strings.Join(e.AffectedBy, ", ")))
		}
	}

	return sb.String()
}
