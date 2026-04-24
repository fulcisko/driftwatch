package drift

import (
	"fmt"
	"sort"
	"strings"
)

// ImpactLevel represents the blast radius of drift for a service.
type ImpactLevel string

const (
	ImpactCritical ImpactLevel = "critical"
	ImpactHigh     ImpactLevel = "high"
	ImpactMedium   ImpactLevel = "medium"
	ImpactLow      ImpactLevel = "low"
)

// ImpactReport summarises the drift impact for a single service.
type ImpactReport struct {
	Service    string      `json:"service"`
	Impact     ImpactLevel `json:"impact"`
	Score      int         `json:"score"`
	DriftCount int         `json:"drift_count"`
	TopKeys    []string    `json:"top_keys"`
}

// AssessImpact computes an ImpactReport for each drifted service in results.
func AssessImpact(results []CompareResult) []ImpactReport {
	var reports []ImpactReport

	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}

		score := 0
		keySet := make(map[string]struct{})
		for _, d := range r.Diffs {
			sev := ClassifyKey(d.Key)
			switch sev {
			case SeverityHigh:
				score += 10
			case SeverityMedium:
				score += 5
			case SeverityLow:
				score += 2
			default:
				score += 1
			}
			keySet[d.Key] = struct{}{}
		}

		topKeys := make([]string, 0, len(keySet))
		for k := range keySet {
			topKeys = append(topKeys, k)
		}
		sort.Strings(topKeys)
		if len(topKeys) > 5 {
			topKeys = topKeys[:5]
		}

		reports = append(reports, ImpactReport{
			Service:    r.Service,
			Impact:     scoreToImpact(score),
			Score:      score,
			DriftCount: len(r.Diffs),
			TopKeys:    topKeys,
		})
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Score > reports[j].Score
	})
	return reports
}

func scoreToImpact(score int) ImpactLevel {
	switch {
	case score >= 20:
		return ImpactCritical
	case score >= 10:
		return ImpactHigh
	case score >= 5:
		return ImpactMedium
	default:
		return ImpactLow
	}
}

// FormatImpact returns a human-readable summary of impact reports.
func FormatImpact(reports []ImpactReport) string {
	if len(reports) == 0 {
		return "No drift impact detected.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-10s %-6s %s\n", "SERVICE", "IMPACT", "SCORE", "TOP KEYS"))
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	for _, r := range reports {
		sb.WriteString(fmt.Sprintf("%-30s %-10s %-6d %s\n",
			r.Service, r.Impact, r.Score, strings.Join(r.TopKeys, ", ")))
	}
	return sb.String()
}
