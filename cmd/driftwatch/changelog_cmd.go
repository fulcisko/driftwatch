package main

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/drift"
)

// runChangelogShow prints changelog entries, optionally filtered by service.
// Usage: driftwatch changelog show <path> [service]
func runChangelogShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: changelog show <path> [service]")
	}
	path := args[0]
	service := ""
	if len(args) >= 2 {
		service = args[1]
	}

	cl, err := drift.LoadChangelog(path)
	if err != nil {
		return fmt.Errorf("load changelog: %w", err)
	}

	filtered := drift.FilterChangelog(cl, service)
	if len(filtered) == 0 {
		fmt.Println("no changelog entries found")
		return nil
	}

	for _, e := range filtered {
		fmt.Fprintf(os.Stdout, "[%s] %s — %d diff(s)\n",
			e.Timestamp.Format("2006-01-02 15:04:05"), e.Service, len(e.Diffs))
		for _, d := range e.Diffs {
			fmt.Fprintf(os.Stdout, "  key=%s expected=%v actual=%v\n", d.Key, d.Expected, d.Actual)
		}
		if e.Note != "" {
			fmt.Fprintf(os.Stdout, "  note: %s\n", e.Note)
		}
	}
	return nil
}

// runChangelogAppend appends a note entry to the changelog.
// Usage: driftwatch changelog append <path> <service> <note>
func runChangelogAppend(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: changelog append <path> <service> <note>")
	}
	entry := drift.ChangelogEntry{
		Service: args[1],
		Note:    args[2],
	}
	if err := drift.AppendChangelog(args[0], entry); err != nil {
		return fmt.Errorf("append changelog: %w", err)
	}
	fmt.Printf("changelog entry added for service %q\n", args[1])
	return nil
}
