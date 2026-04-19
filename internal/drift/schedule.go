package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ScheduleEntry struct {
	Service  string        `json:"service"`
	Interval time.Duration `json:"interval_ns"`
	LastRun  time.Time     `json:"last_run"`
	Enabled  bool          `json:"enabled"`
}

type Schedule struct {
	Entries []ScheduleEntry `json:"entries"`
}

func LoadSchedule(path string) (*Schedule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Schedule{}, nil
		}
		return nil, fmt.Errorf("read schedule: %w", err)
	}
	var s Schedule
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse schedule: %w", err)
	}
	return &s, nil
}

func SaveSchedule(path string, s *Schedule) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schedule: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func UpsertSchedule(s *Schedule, entry ScheduleEntry) {
	for i, e := range s.Entries {
		if e.Service == entry.Service {
			s.Entries[i] = entry
			return
		}
	}
	s.Entries = append(s.Entries, entry)
}

func DueServices(s *Schedule, now time.Time) []string {
	var due []string
	for _, e := range s.Entries {
		if !e.Enabled {
			continue
		}
		if now.Sub(e.LastRun) >= e.Interval {
			due = append(due, e.Service)
		}
	}
	return due
}
