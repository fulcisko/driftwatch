package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/driftwatch/internal/drift"
)

const defaultTagPath = "tags.json"

func runTagAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: tag add <tag> <service>")
	}
	path := envOr("DRIFTWATCH_TAG_PATH", defaultTagPath)
	if err := drift.AddTag(path, args[0], args[1]); err != nil {
		return fmt.Errorf("add tag: %w", err)
	}
	fmt.Fprintf(os.Stdout, "tagged service %q with %q\n", args[1], args[0])
	return nil
}

func runTagRemove(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: tag remove <tag> <service>")
	}
	path := envOr("DRIFTWATCH_TAG_PATH", defaultTagPath)
	if err := drift.RemoveTag(path, args[0], args[1]); err != nil {
		return fmt.Errorf("remove tag: %w", err)
	}
	fmt.Fprintf(os.Stdout, "removed service %q from tag %q\n", args[1], args[0])
	return nil
}

func runTagShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: tag show <tag>")
	}
	path := envOr("DRIFTWATCH_TAG_PATH", defaultTagPath)
	store, err := drift.LoadTags(path)
	if err != nil {
		return fmt.Errorf("load tags: %w", err)
	}
	svcs := drift.FilterByTag(store, args[0])
	if len(svcs) == 0 {
		fmt.Fprintf(os.Stdout, "no services tagged with %q\n", args[0])
		return nil
	}
	fmt.Fprintf(os.Stdout, "services tagged %q:\n  %s\n", args[0], strings.Join(svcs, "\n  "))
	return nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
