package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

func runRemediationAdd(args []string, logPath string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: remediation add <service> <key> <action> [note]")
	}
	service := args[0]
	key := args[1]
	actionStr := args[2]
	note := ""
	if len(args) >= 4 {
		note = args[3]
	}
	var action drift.RemediationAction
	switch actionStr {
	case "apply":
		action = drift.ActionApply
	case "revert":
		action = drift.ActionRevert
	case "ignore":
		action = drift.ActionIgnore
	default:
		return fmt.Errorf("unknown action %q: must be apply, revert, or ignore", actionStr)
	}
	if err := drift.AddRemediation(logPath, service, key, action, note); err != nil {
		return fmt.Errorf("add remediation: %w", err)
	}
	fmt.Fprintf(os.Stdout, "remediation recorded: %s/%s -> %s\n", service, key, action)
	return nil
}

func runRemediationShow(args []string, logPath string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: remediation show <service>")
	}
	service := args[0]
	log, err := drift.LoadRemediations(logPath)
	if err != nil {
		return fmt.Errorf("load remediations: %w", err)
	}
	entries := drift.FilterRemediations(log, service)
	if len(entries) == 0 {
		fmt.Fprintf(os.Stdout, "no remediations found for %s\n", service)
		return nil
	}
	for _, e := range entries {
		note := ""
		if e.Note != "" {
			note = " (" + e.Note + ")"
		}
		fmt.Fprintf(os.Stdout, "[%s] %s/%s -> %s%s\n",
			e.CreatedAt.Format("2006-01-02"), e.Service, e.Key, e.Action, note)
	}
	return nil
}
