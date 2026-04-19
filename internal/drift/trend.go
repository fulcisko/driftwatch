package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type TrendPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Service     string    `json:"service"`
	DriftCount  int       `json:"drift_count"`
	MaxSeverity string    `json:"max_severity"`
}

type TrendReport struct {
	Points []TrendPoint `json:"points"`
}

func AppendTrend(path string, results []CompareResult) error {
	report, _ := LoadTrend(path)
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		point := TrendPoint{
			Timestamp:   time.Now().UTC(),
			Service:     r.Service,
			DriftCount:  len(r.Diffs),
			MaxSeverity: MaxSeverity(r.Diffs).String(),
		}
		report.Points = append(report.Points, point)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadTrend(path string) (TrendReport, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TrendReport{}, nil
		}
		return TrendReport{}, err
	}
	var report TrendReport
	if err := json.Unmarshal(data, &report); err != nil {
		return TrendReport{}, err
	}
	return report, nil
}

func FilterTrend(report TrendReport, service string) []TrendPoint {
	var out []TrendPoint
	for _, p := range report.Points {
		if service == "" || p.Service == service {
			out = append(out, p)
		}
	}
	return out
}

func FormatTrend(points []TrendPoint) string {
	if len(points) == 0 {
		return "no trend data available\n"
	}
	out := ""
	for _, p := range points {
		out += fmt.Sprintf("[%s] %s: %d diffs (max severity: %s)\n",
			p.Timestamp.Format(time.RFC3339), p.Service, p.DriftCount, p.MaxSeverity)
	}
	return out
}
