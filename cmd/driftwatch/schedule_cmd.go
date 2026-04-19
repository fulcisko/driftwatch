package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

func runScheduleSet(args []string, schedPath string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: schedule set <service> <interval_minutes>")
	}
	service := args[0]
	minutes, err := strconv.Atoi(args[1])
	if err != nil || minutes <= 0 {
		return fmt.Errorf("interval must be a positive integer (minutes)")
	}
	s, err := drift.LoadSchedule(schedPath)
	if err != nil {
		return fmt.Errorf("load schedule: %w", err)
	}
	drift.UpsertSchedule(s, drift.ScheduleEntry{
		Service:  service,
		Interval: time.Duration(minutes) * time.Minute,
		LastRun:  time.Time{},
		Enabled:  true,
	})
	if err := drift.SaveSchedule(schedPath, s); err != nil {
		return fmt.Errorf("save schedule: %w", err)
	}
	fmt.Fprintf(os.Stdout, "scheduled %s every %d minute(s)\n", service, minutes)
	return nil
}

func runScheduleShow(schedPath string) error {
	s, err := drift.LoadSchedule(schedPath)
	if err != nil {
		return fmt.Errorf("load schedule: %w", err)
	}
	if len(s.Entries) == 0 {
		fmt.Fprintln(os.Stdout, "no scheduled services")
		return nil
	}
	for _, e := range s.Entries {
		status := "enabled"
		if !e.Enabled {
			status = "disabled"
		}
		fmt.Fprintf(os.Stdout, "service=%-20s interval=%v last_run=%s status=%s\n",
			e.Service, e.Interval, e.LastRun.Format(time.RFC3339), status)
	}
	return nil
}
