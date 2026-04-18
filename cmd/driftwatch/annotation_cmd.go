package main

import (
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/drift"
)

func runAnnotationAdd(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: annotation add <service> <key> <note> [author]")
	}
	service := args[0]
	key := args[1]
	note := args[2]
	author := envOr("DRIFTWATCH_AUTHOR", "unknown")
	if len(args) >= 4 {
		author = args[3]
	}
	path := envOr("DRIFTWATCH_ANNOTATIONS", "annotations.json")
	if err := drift.AddAnnotation(path, service, key, note, author); err != nil {
		return fmt.Errorf("add annotation: %w", err)
	}
	fmt.Fprintf(os.Stdout, "annotation added for %s/%s\n", service, key)
	return nil
}

func runAnnotationShow(args []string) error {
	path := envOr("DRIFTWATCH_ANNOTATIONS", "annotations.json")
	anns, err := drift.LoadAnnotations(path)
	if err != nil {
		return fmt.Errorf("load annotations: %w", err)
	}
	service := ""
	key := ""
	if len(args) >= 1 {
		service = args[0]
	}
	if len(args) >= 2 {
		key = args[1]
	}
	filtered := drift.FilterAnnotations(anns, service, key)
	if len(filtered) == 0 {
		fmt.Fprintln(os.Stdout, "no annotations found")
		return nil
	}
	for _, a := range filtered {
		fmt.Fprintf(os.Stdout, "[%s] %s/%s — %s (by %s at %s)\n",
			a.CreatedAt.Format("2006-01-02"), a.Service, a.Key, a.Note, a.Author, a.CreatedAt.Format("15:04:05"))
	}
	return nil
}
