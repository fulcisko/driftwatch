package drift

import (
	"fmt"
	"sort"
	"strings"
)

// ClusterGroup represents a set of services that share similar drift patterns.
type ClusterGroup struct {
	Label    string   `json:"label"`
	Services []string `json:"services"`
	SharedKeys []string `json:"shared_keys"`
}

// ClusterResult holds all discovered clusters from a drift analysis.
type ClusterResult struct {
	Clusters []ClusterGroup `json:"clusters"`
	Unclustered []string   `json:"unclustered"`
}

// ClusterByDriftPattern groups services that share at least minShared drifted keys.
func ClusterByDriftPattern(results []CompareResult, minShared int) ClusterResult {
	// Build a map of service -> set of drifted keys
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

	assigned := make(map[string]bool)
	var groups []ClusterGroup

	for i := 0; i < len(services); i++ {
		if assigned[services[i]] {
			continue
		}
		group := []string{services[i]}
		shared := copyKeySet(serviceKeys[services[i]])

		for j := i + 1; j < len(services); j++ {
			if assigned[services[j]] {
				continue
			}
			intersect := intersectKeys(shared, serviceKeys[services[j]])
			if len(intersect) >= minShared {
				group = append(group, services[j])
				shared = intersect
				assigned[services[j]] = true
			}
		}

		if len(group) > 1 {
			assigned[services[i]] = true
			sharedList := keySetToSlice(shared)
			groups = append(groups, ClusterGroup{
				Label:      fmt.Sprintf("cluster-%d", len(groups)+1),
				Services:   group,
				SharedKeys: sharedList,
			})
		}
	}

	var unclustered []string
	for _, svc := range services {
		if !assigned[svc] {
			unclustered = append(unclustered, svc)
		}
	}

	return ClusterResult{Clusters: groups, Unclustered: unclustered}
}

// FormatCluster returns a human-readable summary of cluster results.
func FormatCluster(cr ClusterResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Clusters found: %d\n", len(cr.Clusters)))
	for _, g := range cr.Clusters {
		sb.WriteString(fmt.Sprintf("  [%s] services: %s | shared keys: %s\n",
			g.Label,
			strings.Join(g.Services, ", "),
			strings.Join(g.SharedKeys, ", "),
		))
	}
	if len(cr.Unclustered) > 0 {
		sb.WriteString(fmt.Sprintf("Unclustered: %s\n", strings.Join(cr.Unclustered, ", ")))
	}
	return sb.String()
}

func copyKeySet(m map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{}, len(m))
	for k := range m {
		out[k] = struct{}{}
	}
	return out
}

func intersectKeys(a, b map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for k := range a {
		if _, ok := b[k]; ok {
			out[k] = struct{}{}
		}
	}
	return out
}

func keySetToSlice(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
