package drift

// SeverityLevel represents the importance of a detected drift.
type SeverityLevel int

const (
	SeverityNone SeverityLevel = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
)

// String returns a human-readable label for the severity level.
func (s SeverityLevel) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	default:
		return "none"
	}
}

// highSeverityKeys are config keys considered critical when changed.
var highSeverityKeys = map[string]bool{
	"replicas":        true,
	"image":           true,
	"port":            true,
	"memory_limit":    true,
	"cpu_limit":       true,
}

// mediumSeverityKeys are config keys considered moderately important.
var mediumSeverityKeys = map[string]bool{
	"env":             true,
	"log_level":       true,
	"timeout":         true,
}

// ClassifyKey returns the severity level for a given config key.
func ClassifyKey(key string) SeverityLevel {
	if highSeverityKeys[key] {
		return SeverityHigh
	}
	if mediumSeverityKeys[key] {
		return SeverityMedium
	}
	if key != "" {
		return SeverityLow
	}
	return SeverityNone
}

// MaxSeverity returns the highest severity level across a slice of diff keys.
func MaxSeverity(keys []string) SeverityLevel {
	max := SeverityNone
	for _, k := range keys {
		if s := ClassifyKey(k); s > max {
			max = s
		}
	}
	return max
}
