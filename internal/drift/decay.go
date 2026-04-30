package drift

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// DecayEntry holds the computed drift decay score for a service.
type DecayEntry struct {
	Service   string    `json:"service"`
	Score     float64   `json:"score"`
	LastSeen  time.Time `json:"last_seen"`
	AgeDays   float64   `json:"age_days"`
	Decayed   bool      `json:"decayed"`
}

// DecayOptions controls decay behaviour.
type DecayOptions struct {
	// HalfLifeDays is the number of days after which a drift score is halved.
	HalfLifeDays float64
	// ThresholdScore is the minimum decayed score to still be considered active.
	ThresholdScore float64
}

// DefaultDecayOptions returns sensible defaults.
func DefaultDecayOptions() DecayOptions {
	return DecayOptions{
		HalfLifeDays:   7.0,
		ThresholdScore: 0.1,
	}
}

// ApplyDecay computes exponential decay on drift results relative to the
// timestamps recorded in history. Services with no history entry are scored
// at full weight using the current time as their reference.
func ApplyDecay(results []CompareResult, history []HistoryEntry, opts DecayOptions) []DecayEntry {
	if opts.HalfLifeDays <= 0 {
		opts.HalfLifeDays = DefaultDecayOptions().HalfLifeDays
	}

	// Build a map from service -> most recent history timestamp.
	latest := make(map[string]time.Time)
	for _, h := range history {
		if t, ok := latest[h.Service]; !ok || h.Timestamp.After(t) {
			latest[h.Service] = h.Timestamp
		}
	}

	now := time.Now().UTC()
	k := math.Log(2) / opts.HalfLifeDays

	entries := make([]DecayEntry, 0, len(results))
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		rawScore := float64(len(r.Diffs))
		ref := now
		if t, ok := latest[r.Service]; ok {
			ref = t
		}
		ageDays := now.Sub(ref).Hours() / 24.0
		decayed := rawScore * math.Exp(-k*ageDays)
		entries = append(entries, DecayEntry{
			Service:  r.Service,
			Score:    math.Round(decayed*1000) / 1000,
			LastSeen: ref,
			AgeDays:  math.Round(ageDays*10) / 10,
			Decayed:  decayed < opts.ThresholdScore,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	return entries
}

// FormatDecay returns a human-readable table of decay entries.
func FormatDecay(entries []DecayEntry) string {
	if len(entries) == 0 {
		return "no active drift decay entries\n"
	}
	out := fmt.Sprintf("%-24s %10s %10s %8s\n", "SERVICE", "SCORE", "AGE(days)", "DECAYED")
	out += fmt.Sprintf("%-24s %10s %10s %8s\n", "-------", "-----", "---------", "-------")
	for _, e := range entries {
		decayedStr := "no"
		if e.Decayed {
			decayedStr = "yes"
		}
		out += fmt.Sprintf("%-24s %10.3f %10.1f %8s\n", e.Service, e.Score, e.AgeDays, decayedStr)
	}
	return out
}
