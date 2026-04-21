package drift

import (
	"fmt"
	"sort"
	"strings"
)

// GroupByField defines the field to group results by.
type GroupByField string

const (
	GroupByService  GroupByField = "service"
	GroupBySeverity GroupByField = "severity"
	GroupByKey      GroupByField = "key"
)

// GroupedResults holds drift results organized by a grouping key.
type GroupedResults struct {
	Field  GroupByField
	Groups map[string][]CompareResult
}

// GroupResults groups a slice of CompareResult by the specified field.
// Supported fields: "service", "severity", "key".
func GroupResults(results []CompareResult, field GroupByField) (*GroupedResults, error) {
	switch field {
	case GroupByService, GroupBySeverity, GroupByKey:
		// valid
	default:
		return nil, fmt.Errorf("unsupported group-by field: %q", field)
	}

	groups := make(map[string][]CompareResult)

	for _, r := range results {
		keys := groupKeys(r, field)
		for _, k := range keys {
			groups[k] = append(groups[k], r)
		}
	}

	return &GroupedResults{
		Field:  field,
		Groups: groups,
	}, nil
}

// groupKeys returns the grouping key(s) for a single CompareResult.
// A result may map to multiple keys when grouping by key or severity,
// since a single service can have many diffs.
func groupKeys(r CompareResult, field GroupByField) []string {
	switch field {
	case GroupByService:
		return []string{r.Service}
	case GroupBySeverity:
		if len(r.Diffs) == 0 {
			return []string{"none"}
		}
		seen := make(map[string]struct{})
		var keys []string
		for _, d := range r.Diffs {
			s := strings.ToLower(ClassifyKey(d.Key).String())
			if _, ok := seen[s]; !ok {
				seen[s] = struct{}{}
				keys = append(keys, s)
			}
		}
		return keys
	case GroupByKey:
		if len(r.Diffs) == 0 {
			return nil
		}
		seen := make(map[string]struct{})
		var keys []string
		for _, d := range r.Diffs {
			if _, ok := seen[d.Key]; !ok {
				seen[d.Key] = struct{}{}
				keys = append(keys, d.Key)
			}
		}
		return keys
	}
	return nil
}

// SortedGroupKeys returns the group keys in deterministic alphabetical order.
func (g *GroupedResults) SortedGroupKeys() []string {
	keys := make([]string, 0, len(g.Groups))
	for k := range g.Groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FormatGrouped returns a human-readable representation of grouped results.
func FormatGrouped(g *GroupedResults) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Grouped by: %s\n", g.Field))
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	for _, groupKey := range g.SortedGroupKeys() {
		results := g.Groups[groupKey]
		totalDiffs := 0
		for _, r := range results {
			totalDiffs += len(r.Diffs)
		}
		sb.WriteString(fmt.Sprintf("[%s] %d service(s), %d diff(s)\n", groupKey, len(results), totalDiffs))
		for _, r := range results {
			if len(r.Diffs) > 0 {
				sb.WriteString(fmt.Sprintf("  - %s (%d diffs)\n", r.Service, len(r.Diffs)))
			}
		}
	}
	return sb.String()
}
