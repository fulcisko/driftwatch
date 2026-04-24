package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

// runAuditAppend records a new audit event to the log file at the given path.
// Arguments: <path> <action> <service> <user> [detail]
func runAuditAppend(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: audit append <path> <action> <service> <user> [detail]")
	}
	path := args[0]
	action := args[1]
	service := args[2]
	user := args[3]
	detail := ""
	if len(args) >= 5 {
		detail = args[4]
	}
	if err := drift.AppendAuditEvent(path, action, service, user, detail); err != nil {
		return fmt.Errorf("audit append: %w", err)
	}
	fmt.Fprintf(os.Stdout, "audit event recorded: action=%s service=%s user=%s\n", action, service, user)
	return nil
}

// runAuditShow loads and displays audit log entries from the given path,
// optionally filtering by service and action.
// Arguments: <path> [service] [action]
func runAuditShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: audit show <path> [service] [action]")
	}
	path := args[0]
	service := ""
	action := ""
	if len(args) >= 2 {
		service = args[1]
	}
	if len(args) >= 3 {
		action = args[2]
	}
	events, err := drift.LoadAuditLog(path)
	if err != nil {
		return fmt.Errorf("load audit log: %w", err)
	}
	filtered := drift.FilterAuditLog(events, service, action)
	if len(filtered) == 0 {
		fmt.Fprintln(os.Stdout, "no audit events found")
		return nil
	}
	fmt.Print(drift.FormatAuditLog(filtered))
	return nil
}
