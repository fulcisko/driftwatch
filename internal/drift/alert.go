package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type AlertSeverity string

const (
	AlertHigh   AlertSeverity = "high"
	AlertMedium AlertSeverity = "medium"
	AlertLow    AlertSeverity = "low"
)

type Alert struct {
	Service   string        `json:"service"`
	Key       string        `json:"key"`
	Severity  AlertSeverity `json:"severity"`
	Message   string        `json:"message"`
	Timestamp time.Time     `json:"timestamp"`
}

type AlertConfig struct {
	MinSeverity AlertSeverity `json:"min_severity"`
}

func GenerateAlerts(results []CompareResult, cfg AlertConfig) []Alert {
	var alerts []Alert
	for _, r := range results {
		for _, d := range r.Diffs {
			level := ClassifyKey(d.Key)
			if !meetsSeverityThreshold(level, cfg.MinSeverity) {
				continue
			}
			alerts = append(alerts, Alert{
				Service:   r.Service,
				Key:       d.Key,
				Severity:  AlertSeverity(level.String()),
				Message:   fmt.Sprintf("drift detected: %s changed from %q to %q", d.Key, d.Expected, d.Actual),
				Timestamp: time.Now().UTC(),
			})
		}
	}
	return alerts
}

func SaveAlerts(path string, alerts []Alert) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("save alerts: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(alerts)
}

func LoadAlerts(path string) ([]Alert, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("load alerts: %w", err)
	}
	defer f.Close()
	var alerts []Alert
	if err := json.NewDecoder(f).Decode(&alerts); err != nil {
		return nil, fmt.Errorf("decode alerts: %w", err)
	}
	return alerts, nil
}

func meetsSeverityThreshold(level SeverityLevel, min AlertSeverity) bool {
	order := map[AlertSeverity]int{AlertLow: 1, AlertMedium: 2, AlertHigh: 3}
	return order[AlertSeverity(level.String())] >= order[min]
}
