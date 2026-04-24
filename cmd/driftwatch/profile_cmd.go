package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/driftwatch/driftwatch/internal/drift"
)

const defaultProfilePath = "profiles.json"

func runProfileSave(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: profile save <name> <min_severity> [ignore_keys,...] [service_prefix]")
	}
	name := args[0]
	minSev := args[1]
	var ignoreKeys []string
	if len(args) >= 3 && args[2] != "" {
		for _, k := range strings.Split(args[2], ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				ignoreKeys = append(ignoreKeys, k)
			}
		}
	}
	var prefix string
	if len(args) >= 4 {
		prefix = args[3]
	}
	path := envOr("DRIFTWATCH_PROFILES", defaultProfilePath)
	p := drift.Profile{
		Name:          name,
		MinSeverity:   minSev,
		IgnoreKeys:    ignoreKeys,
		ServicePrefix: prefix,
	}
	if err := drift.SaveProfile(path, p); err != nil {
		return fmt.Errorf("save profile: %w", err)
	}
	fmt.Printf("profile %q saved\n", name)
	return nil
}

func runProfileShow(args []string) error {
	path := envOr("DRIFTWATCH_PROFILES", defaultProfilePath)
	if len(args) == 1 {
		p, ok := drift.GetProfile(path, args[0])
		if !ok {
			return fmt.Errorf("profile %q not found", args[0])
		}
		fmt.Printf("Name:          %s\n", p.Name)
		fmt.Printf("Description:   %s\n", p.Description)
		fmt.Printf("MinSeverity:   %s\n", p.MinSeverity)
		fmt.Printf("IgnoreKeys:    %s\n", strings.Join(p.IgnoreKeys, ", "))
		fmt.Printf("ServicePrefix: %s\n", p.ServicePrefix)
		return nil
	}
	profiles, err := drift.LoadProfiles(path)
	if err != nil {
		return fmt.Errorf("load profiles: %w", err)
	}
	if len(profiles) == 0 {
		fmt.Println("no profiles defined")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tMIN_SEVERITY\tIGNORE_KEYS\tSERVICE_PREFIX")
	for _, p := range profiles {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			p.Name, p.MinSeverity,
			strings.Join(p.IgnoreKeys, ","),
			p.ServicePrefix,
		)
	}
	return w.Flush()
}

func runProfileRemove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: profile remove <name>")
	}
	path := envOr("DRIFTWATCH_PROFILES", defaultProfilePath)
	if err := drift.RemoveProfile(path, args[0]); err != nil {
		return fmt.Errorf("remove profile: %w", err)
	}
	fmt.Printf("profile %q removed\n", args[0])
	return nil
}
