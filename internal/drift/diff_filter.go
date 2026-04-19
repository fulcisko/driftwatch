package drift

import "strings"

// DiffFilterOptions controls which diffs are included in output.
type DiffFilterOptions struct {
	MinSeverity  SeverityLevel
	ServicePrefix string
	OnlyDrifted  bool
	ExcludeKeys  []string
}

// ApplyDiffFilter returns a filtered copy of results based on DiffFilterOptions.
func ApplyDiffFilter(results []CompareResult, opts DiffFilterOptions) []CompareResult {
	var out []CompareResult
	for _, r := range results {
		if opts.OnlyDrifted && len(r.Diffs) == 0 {
			continue
		}
		if opts.ServicePrefix != "" && !strings.HasPrefix(r.Service, opts.ServicePrefix) {
			continue
		}
		filtered := filterDiffsByOptions(r.Diffs, opts)
		if opts.OnlyDrifted && len(filtered) == 0 {
			continue
		}
		out = append(out, CompareResult{Service: r.Service, Diffs: filtered})
	}
	return out
}

func filterDiffsByOptions(diffs []DiffEntry, opts DiffFilterOptions) []DiffEntry {
	var out []DiffEntry
	for _, d := range diffs {
		if isExcludedKey(d.Key, opts.ExcludeKeys) {
			continue
		}
		if ClassifyKey(d.Key) < opts.MinSeverity {
			continue
		}
		out = append(out, d)
	}
	return out
}

func isExcludedKey(key string, excludes []string) bool {
	for _, e := range excludes {
		if strings.EqualFold(key, e) {
			return true
		}
	}
	return false
}
