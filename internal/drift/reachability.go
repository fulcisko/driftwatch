package drift

import (
	"fmt"
	"sort"
	"strings"
)

// ReachabilityResult holds transitive drift impact for a service.
type ReachabilityResult struct {
	Service      string   `json:"service"`
	AffectedBy   []string `json:"affected_by"`
	ReachableFrom []string `json:"reachable_from"`
	TotalExposure int      `json:"total_exposure"`
}

// BuildReachability walks the dependency graph and computes which drifted
// services can transitively affect each service.
func BuildReachability(results []CompareResult, deps []Dependency) []ReachabilityResult {
	// Build adjacency: dependents[svc] = list of services that depend on svc
	dependents := map[string][]string{}
	for _, d := range deps {
		dependents[d.DependsOn] = append(dependents[d.DependsOn], d.Service)
	}

	// Collect drifted service names
	driftedSet := map[string]bool{}
	for _, r := range results {
		if len(r.Diffs) > 0 {
			driftedSet[r.Service] = true
		}
	}

	// For each service, BFS upstream to find all drifted ancestors
	var output []ReachabilityResult
	for _, r := range results {
		affected := reachableAncestors(r.Service, driftedSet, deps)
		reachable := dependents[r.Service]

		sort.Strings(affected)
		sort.Strings(reachable)

		output = append(output, ReachabilityResult{
			Service:       r.Service,
			AffectedBy:    affected,
			ReachableFrom: reachable,
			TotalExposure: len(affected),
		})
	}

	sort.Slice(output, func(i, j int) bool {
		if output[i].TotalExposure != output[j].TotalExposure {
			return output[i].TotalExposure > output[j].TotalExposure
		}
		return output[i].Service < output[j].Service
	})
	return output
}

func reachableAncestors(service string, driftedSet map[string]bool, deps []Dependency) []string {
	// Build: dependencies[svc] = services svc depends on
	dependencies := map[string][]string{}
	for _, d := range deps {
		dependencies[d.Service] = append(dependencies[d.Service], d.DependsOn)
	}

	visited := map[string]bool{}
	queue := []string{service}
	var affected []string

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, dep := range dependencies[curr] {
			if visited[dep] {
				continue
			}
			visited[dep] = true
			if driftedSet[dep] {
				affected = append(affected, dep)
			}
			queue = append(queue, dep)
		}
	}
	return affected
}

// FormatReachability returns a human-readable summary of reachability results.
func FormatReachability(results []ReachabilityResult) string {
	if len(results) == 0 {
		return "no reachability data\n"
	}
	var sb strings.Builder
	sb.WriteString("=== Drift Reachability ===\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("service: %s (exposure: %d)\n", r.Service, r.TotalExposure))
		if len(r.AffectedBy) > 0 {
			sb.WriteString(fmt.Sprintf("  affected by drifted: %s\n", strings.Join(r.AffectedBy, ", ")))
		}
		if len(r.ReachableFrom) > 0 {
			sb.WriteString(fmt.Sprintf("  downstream dependents: %s\n", strings.Join(r.ReachableFrom, ", ")))
		}
	}
	return sb.String()
}
