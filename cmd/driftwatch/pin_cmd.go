package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

func runPinAdd(args []string, pinFile string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: pin add <service> <key> <expected> [comment]")
	}
	service, key, expected := args[0], args[1], args[2]
	comment := ""
	if len(args) >= 4 {
		comment = args[3]
	}
	if err := drift.AddPin(pinFile, service, key, expected, comment); err != nil {
		return fmt.Errorf("add pin: %w", err)
	}
	fmt.Fprintf(os.Stdout, "pinned %s/%s = %q\n", service, key, expected)
	return nil
}

func runPinRemove(args []string, pinFile string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: pin remove <service> <key>")
	}
	if err := drift.RemovePin(pinFile, args[0], args[1]); err != nil {
		return fmt.Errorf("remove pin: %w", err)
	}
	fmt.Fprintf(os.Stdout, "removed pin %s/%s\n", args[0], args[1])
	return nil
}

func runPinShow(pinFile string) error {
	list, err := drift.LoadPins(pinFile)
	if err != nil {
		return fmt.Errorf("load pins: %w", err)
	}
	if len(list.Pins) == 0 {
		fmt.Println("no pinned keys")
		return nil
	}
	for _, p := range list.Pins {
		comment := ""
		if p.Comment != "" {
			comment = fmt.Sprintf(" # %s", p.Comment)
		}
		fmt.Printf("  %s/%s = %q (pinned %s)%s\n",
			p.Service, p.Key, p.Expected,
			p.PinnedAt.Format("2006-01-02"), comment)
	}
	return nil
}
