package drift

import (
	"fmt"
	"sort"
	"strings"
)

// ExposureLevel represents how externally visible a drifted service is.
type ExposureLevel string

const (
	ExposurePublic   ExposureLevel = "public"
	ExposureInternal ExposureLevel = "internal"
	ExposurePrivate  ExposureLevel = "private"
	ExposureUnknown  ExposureLevel = "unknown"
)

// ExposureEntry holds exposure metadata for a service.
type ExposureEntry struct {
	Service  string        `json:"service"`
	Level    ExposureLevel `json:"level"`
	DriftCount int         `json:"drift_count"`
	RiskScore  int         `json:"risk_score"`
	TopKeys  []string      `json:"top_keys"`
}

// exposureWeight maps exposure level to a risk multiplier.
var exposureWeight = map[ExposureLevel]int{
	ExposurePublic:   3,
	ExposureInternal: 2,
	ExposurePrivate:  1,
	ExposureUnknown:  1,
}

// AssessExposure computes exposure risk for each drifted service.
// levelMap maps service name to its ExposureLevel; missing entries default to unknown.
func AssessExposure(results []CompareResult, levelMap map[string]ExposureLevel) []ExposureEntry {
	var entries []ExposureEntry

	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}

		level, ok := levelMap[r.Service]
		if !ok {
			level = ExposureUnknown
		}

		weight := exposureWeight[level]
		severitySum := 0
		keySet := make(map[string]struct{})

		for _, d := range r.Diffs {
			sev := ClassifyKey(d.Key)
			severitySum += int(sev)
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

		entries = append(entries, ExposureEntry{
			Service:    r.Service,
			Level:      level,
			DriftCount: len(r.Diffs),
			RiskScore:  severitySum * weight,
			TopKeys:    topKeys,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].RiskScore != entries[j].RiskScore {
			return entries[i].RiskScore > entries[j].RiskScore
		}
		return entries[i].Service < entries[j].Service
	})

	return entries
}

// FormatExposure returns a human-readable table of exposure entries.
func FormatExposure(entries []ExposureEntry) string {
	if len(entries) == 0 {
		return "no exposure risk detected\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-24s %-10s %6s %6s  top keys\n", "service", "exposure", "drifts", "risk"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-24s %-10s %6d %6d  %s\n",
			e.Service, string(e.Level), e.DriftCount, e.RiskScore,
			strings.Join(e.TopKeys, ", ")))
	}
	return sb.String()
}
