package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/user/driftwatch/internal/drift"
)

const defaultWatchlistPath = "watchlist.json"

// runWatchlistAdd adds a service with a drift threshold to the watchlist.
// Usage: driftwatch watchlist add <service> <threshold>
func runWatchlistAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: driftwatch watchlist add <service> <threshold>")
	}
	service := args[0]
	threshold, err := strconv.Atoi(args[1])
	if err != nil || threshold < 1 {
		return fmt.Errorf("threshold must be a positive integer")
	}
	if err := drift.AddToWatchlist(defaultWatchlistPath, service, threshold); err != nil {
		return err
	}
	fmt.Printf("Added %q to watchlist with threshold %d\n", service, threshold)
	return nil
}

// runWatchlistRemove removes a service from the watchlist.
// Usage: driftwatch watchlist remove <service>
func runWatchlistRemove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: driftwatch watchlist remove <service>")
	}
	service := args[0]
	if err := drift.RemoveFromWatchlist(defaultWatchlistPath, service); err != nil {
		return err
	}
	fmt.Printf("Removed %q from watchlist\n", service)
	return nil
}

// runWatchlistShow prints all entries in the watchlist as a formatted table.
func runWatchlistShow(_ []string) error {
	wl, err := drift.LoadWatchlist(defaultWatchlistPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Watchlist is empty.")
			return nil
		}
		return err
	}
	if len(wl.Entries) == 0 {
		fmt.Println("Watchlist is empty.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tTHRESHOLD\tADDED AT")
	for _, e := range wl.Entries {
		fmt.Fprintf(w, "%s\t%d\t%s\n", e.Service, e.Threshold, e.AddedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
