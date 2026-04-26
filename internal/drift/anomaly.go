package drift

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// AnomalyResult represents a service whose drift pattern deviates significantly
// from the norm across all observed services.
type AnomalyResult struct {
	Service    string   `json:"service"`
	Score      float64  `json:"score"`
	Reason     string   `json:"reason"`
	Anomalous  bool     `json:"anomalous"`
	TopKeys    []string `json:"top_keys,omitempty"`
}

// DetectAnomalies identifies services whose drift score deviates beyond the
// given threshold (in standard deviations) from the mean across all services.
// A threshold of 1.5 is a reasonable default for most deployments.
func DetectAnomalies(results []CompareResult, stddevThreshold float64) []AnomalyResult {
	if stddevThreshold <= 0 {
		stddevThreshold = 1.5
	}

	scores := make(map[string]float64, len(results))
	topKeys := make(map[string][]string, len(results))

	for _, r := range results {
		if r.Clean {
			scores[r.Service] = 0
			continue
		}
		s := 0.0
		keyCount := make(map[string]int)
		for _, d := range r.Diffs {
			sev := ClassifyKey(d.Key)
			switch sev {
			case SeverityHigh:
				s += 3.0
			case SeverityMedium:
				s += 2.0
			case SeverityLow:
				s += 1.0
			default:
				s += 0.5
			}
			keyCount[d.Key]++
		}
		scores[r.Service] = s
		// collect top contributing keys
		type kv struct {
			key   string
			count int
		}
		var kvs []kv
		for k, c := range keyCount {
			kvs = append(kvs, kv{k, c})
		}
		sort.Slice(kvs, func(i, j int) bool { return kvs[i].count > kvs[j].count })
		var keys []string
		for i, item := range kvs {
			if i >= 3 {
				break
			}
			keys = append(keys, item.key)
		}
		topKeys[r.Service] = keys
	}

	mean, stddev := meanStddev(scores)

	var out []AnomalyResult
	for _, r := range results {
		score := scores[r.Service]
		anomalous := stddev > 0 && (score-mean)/stddev > stddevThreshold
		reason := buildReason(score, mean, stddev, stddevThreshold)
		out = append(out, AnomalyResult{
			Service:   r.Service,
			Score:     score,
			Reason:    reason,
			Anomalous: anomalous,
			TopKeys:   topKeys[r.Service],
		})
	}

	// sort anomalous services first, then by descending score
	sort.Slice(out, func(i, j int) bool {
		if out[i].Anomalous != out[j].Anomalous {
			return out[i].Anomalous
		}
		return out[i].Score > out[j].Score
	})
	return out
}

// FormatAnomalies returns a human-readable summary of anomaly detection results.
func FormatAnomalies(results []AnomalyResult) string {
	var sb strings.Builder
	sb.WriteString("Anomaly Detection Report\n")
	sb.WriteString(strings.Repeat("─", 40) + "\n")
	for _, r := range results {
		marker := " "
		if r.Anomalous {
			marker = "!"
		}
		sb.WriteString(fmt.Sprintf("[%s] %-30s score=%.2f\n", marker, r.Service, r.Score))
		if r.Anomalous {
			sb.WriteString(fmt.Sprintf("    reason: %s\n", r.Reason))
			if len(r.TopKeys) > 0 {
				sb.WriteString(fmt.Sprintf("    top keys: %s\n", strings.Join(r.TopKeys, ", ")))
			}
		}
	}
	return sb.String()
}

// meanStddev computes the mean and population standard deviation of a score map.
func meanStddev(scores map[string]float64) (mean, stddev float64) {
	if len(scores) == 0 {
		return 0, 0
	}
	for _, v := range scores {
		mean += v
	}
	mean /= float64(len(scores))
	for _, v := range scores {
		diff := v - mean
		stddev += diff * diff
	}
	stddev = math.Sqrt(stddev / float64(len(scores)))
	return mean, stddev
}

func buildReason(score, mean, stddev, threshold float64) string {
	if stddev == 0 {
		return "all services have identical drift scores"
	}
	z := (score - mean) / stddev
	return fmt.Sprintf("z-score=%.2f exceeds threshold=%.2f (mean=%.2f, stddev=%.2f)", z, threshold, mean, stddev)
}
