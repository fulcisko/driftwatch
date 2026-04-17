package source

import (
	"fmt"
	"strings"
)

// parseConfigBody parses a simple KEY=VALUE newline-delimited config body.
// Lines beginning with '#' are treated as comments and ignored.
func parseConfigBody(data []byte) (map[string]string, error) {
	fields := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", i+1, line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", i+1)
		}
		fields[key] = val
	}

	return fields, nil
}
