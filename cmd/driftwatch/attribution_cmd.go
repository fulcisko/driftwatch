package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

// runAttributionAdd adds an attribution record for a service/key pair.
// Usage: driftwatch attribution add <file> <service> <key> <owner> [team] [reason]
func runAttributionAdd(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: attribution add <file> <service> <key> <owner> [team] [reason]")
	}
	path := args[0]
	service := args[1]
	key := args[2]
	owner := args[3]
	team := ""
	reason := ""
	if len(args) >= 5 {
		team = args[4]
	}
	if len(args) >= 6 {
		reason = args[5]
	}
	if err := drift.AddAttribution(path, service, key, owner, team, reason); err != nil {
		return fmt.Errorf("add attribution: %w", err)
	}
	fmt.Fprintf(os.Stdout, "attribution recorded: service=%s key=%s owner=%s\n", service, key, owner)
	return nil
}

// runAttributionShow prints attribution records, optionally filtered by service.
// Usage: driftwatch attribution show <file> [service]
func runAttributionShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: attribution show <file> [service]")
	}
	path := args[0]
	service := ""
	if len(args) >= 2 {
		service = args[1]
	}
	store, err := drift.LoadAttributions(path)
	if err != nil {
		return fmt.Errorf("load attributions: %w", err)
	}
	entries := drift.FilterAttributions(store, service)
	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "no attribution records found")
		return nil
	}
	for _, e := range entries {
		team := e.Team
		if team == "" {
			team = "(none)"
		}
		reason := e.Reason
		if reason == "" {
			reason = "(none)"
		}
		fmt.Fprintf(os.Stdout, "[%s] service=%s key=%s owner=%s team=%s reason=%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Service, e.Key, e.Owner, team, reason)
	}
	return nil
}
