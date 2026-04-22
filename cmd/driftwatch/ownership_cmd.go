package main

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
)

const defaultOwnershipPath = "ownership.json"

// runOwnershipAdd adds or updates an ownership entry.
// Usage: driftwatch ownership add <service> <team> [contact]
func runOwnershipAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ownership add <service> <team> [contact]")
	}
	service := args[0]
	team := args[1]
	contact := ""
	if len(args) >= 3 {
		contact = args[2]
	}
	path := envOr("DRIFTWATCH_OWNERSHIP_PATH", defaultOwnershipPath)
	if err := drift.AddOwner(path, service, team, contact); err != nil {
		return fmt.Errorf("add owner: %w", err)
	}
	fmt.Fprintf(os.Stdout, "ownership added: service=%s team=%s\n", service, team)
	return nil
}

// runOwnershipRemove removes an ownership entry for a service.
// Usage: driftwatch ownership remove <service>
func runOwnershipRemove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ownership remove <service>")
	}
	path := envOr("DRIFTWATCH_OWNERSHIP_PATH", defaultOwnershipPath)
	if err := drift.RemoveOwner(path, args[0]); err != nil {
		return fmt.Errorf("remove owner: %w", err)
	}
	fmt.Fprintf(os.Stdout, "ownership removed: service=%s\n", args[0])
	return nil
}

// runOwnershipShow prints all ownership entries, optionally looking up a single service.
// Usage: driftwatch ownership show [service]
func runOwnershipShow(args []string) error {
	path := envOr("DRIFTWATCH_OWNERSHIP_PATH", defaultOwnershipPath)
	om, err := drift.LoadOwnership(path)
	if err != nil {
		return fmt.Errorf("load ownership: %w", err)
	}
	if len(args) >= 1 {
		o, ok := drift.LookupOwner(om, args[0])
		if !ok {
			fmt.Fprintf(os.Stdout, "no owner found for service %q\n", args[0])
			return nil
		}
		fmt.Fprintf(os.Stdout, "service=%-20s team=%-20s contact=%s\n", o.Service, o.Team, o.Contact)
		return nil
	}
	if len(om.Owners) == 0 {
		fmt.Fprintln(os.Stdout, "no ownership entries found")
		return nil
	}
	fmt.Fprintf(os.Stdout, "%-20s %-20s %s\n", "SERVICE", "TEAM", "CONTACT")
	for _, o := range om.Owners {
		fmt.Fprintf(os.Stdout, "%-20s %-20s %s\n", o.Service, o.Team, o.Contact)
	}
	return nil
}
