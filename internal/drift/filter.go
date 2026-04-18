package drift

import "strings"

// FilterOptions controls which drift results are included in output.
type FilterOptions struct {
	// OnlyDrifted filters out services with no drift when true.
	OnlyDrifted bool
	// ServicePrefix limits results to services whose name starts with the given prefix.
	ServicePrefix string
	// IgnoreKeys is a set of config keys to exclude from drift comparison.
	IgnoreKeys map[string]struct{}
}

// NewFilterOptions returns a FilterOptions with an initialised IgnoreKeys map.
func NewFilterOptions() FilterOptions {
	return FilterOptions{
		IgnoreKeys: make(map[string]struct{}),
	}
}

// AddIgnoreKey registers a key to be excluded from drift results.
func (f *FilterOptions) AddIgnoreKey(key string) {
	f.IgnoreKeys[key] = struct{}{}
}

// ShouldIgnoreKey reports whether the given key should be skipped.
func (f *FilterOptions) ShouldIgnoreKey(key string) bool {
	_, ok := f.IgnoreKeys[key]
	return ok
}

// ApplyToResults filters a slice of CompareResult according to the options.
func (f *FilterOptions) ApplyToResults(results []CompareResult) []CompareResult {
	var out []CompareResult
	for _, r := range results {
		if f.ServicePrefix != "" && !strings.HasPrefix(r.Service, f.ServicePrefix) {
			continue
		}
		filtered := filterDiffs(r.Diffs, f.IgnoreKeys)
		if f.OnlyDrifted && len(filtered) == 0 {
			continue
		}
		out = append(out, CompareResult{Service: r.Service, Diffs: filtered})
	}
	return out
}

func filterDiffs(diffs []Diff, ignoreKeys map[string]struct{}) []Diff {
	if len(ignoreKeys) == 0 {
		return diffs
	}
	var out []Diff
	for _, d := range diffs {
		if _, skip := ignoreKeys[d.Key]; skip {
			continue
		}
		out = append(out, d)
	}
	return out
}
