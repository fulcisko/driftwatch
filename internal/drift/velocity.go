package drift

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// VelocityEntry represents the drift velocity for a single service.
type VelocityEntry struct {
	Service       string  `json:"service"`
	DriftsPerDay  float64 `json:"drifts_per_day"`
	TotalDrifts   int     `json:"total_drifts"`
	SpanDays      float64 `json:"span_days"`
	Accelerating  bool    `json:"accelerating"`
}

// VelocityReport holds velocity entries for all services.
type VelocityReport struct {
	GeneratedAt time.Time       `json:"generated_at"`
	Entries     []VelocityEntry `json:"entries"`
}

// ComputeVelocity calculates drift velocity from trend history.
// It measures how many drifts per day each service accumulates over the
// observed window, and flags services whose rate is increasing.
func ComputeVelocity(history []TrendEntry, minSpanHours float64) VelocityReport {
	type bucket struct {
		times  []time.Time
		counts []int
	}
	byService := map[string]*bucket{}

	for _, e := range history {
		if _, ok := byService[e.Service]; !ok {
			byService[e.Service] = &bucket{}
		}
		b := byService[e.Service]
		b.times = append(b.times, e.RecordedAt)
		b.counts = append(b.counts, e.DriftCount)
	}

	var entries []VelocityEntry
	for svc, b := range byService {
		if len(b.times) < 2 {
			continue
		}
		sort.Slice(b.times, func(i, j int) bool { return b.times[i].Before(b.times[j]) })
		span := b.times[len(b.times)-1].Sub(b.times[0]).Hours()
		if span < minSpanHours {
			continue
		}
		total := 0
		for _, c := range b.counts {
			total += c
		}
		spanDays := span / 24.0
		rate := float64(total) / spanDays

		// Detect acceleration: compare first-half rate vs second-half rate.
		mid := len(b.counts) / 2
		firstHalf, secondHalf := 0, 0
		for i, c := range b.counts {
			if i < mid {
				firstHalf += c
			} else {
				secondHalf += c
			}
		}
		accel := secondHalf > firstHalf

		entries = append(entries, VelocityEntry{
			Service:      svc,
			DriftsPerDay: math.Round(rate*100) / 100,
			TotalDrifts:  total,
			SpanDays:     math.Round(spanDays*100) / 100,
			Accelerating: accel,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].DriftsPerDay > entries[j].DriftsPerDay
	})

	return VelocityReport{
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
	}
}

// FormatVelocity returns a human-readable summary of the velocity report.
func FormatVelocity(r VelocityReport) string {
	if len(r.Entries) == 0 {
		return "no velocity data available\n"
	}
	out := fmt.Sprintf("drift velocity report (%s)\n", r.GeneratedAt.Format(time.RFC3339))
	out += fmt.Sprintf("%-30s %12s %10s %10s %12s\n", "service", "drifts/day", "total", "span(days)", "accelerating")
	for _, e := range r.Entries {
		acc := ""
		if e.Accelerating {
			acc = "yes"
		}
		out += fmt.Sprintf("%-30s %12.2f %10d %10.2f %12s\n",
			e.Service, e.DriftsPerDay, e.TotalDrifts, e.SpanDays, acc)
	}
	return out
}
