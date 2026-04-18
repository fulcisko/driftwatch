package drift

import (
	"testing"
)

func TestClassifyKey_High(t *testing.T) {
	for _, key := range []string{"replicas", "image", "port", "memory_limit", "cpu_limit"} {
		if got := ClassifyKey(key); got != SeverityHigh {
			t.Errorf("ClassifyKey(%q) = %v, want high", key, got)
		}
	}
}

func TestClassifyKey_Medium(t *testing.T) {
	for _, key := range []string{"env", "log_level", "timeout"} {
		if got := ClassifyKey(key); got != SeverityMedium {
			t.Errorf("ClassifyKey(%q) = %v, want medium", key, got)
		}
	}
}

func TestClassifyKey_Low(t *testing.T) {
	if got := ClassifyKey("some_unknown_key"); got != SeverityLow {
		t.Errorf("ClassifyKey(unknown) = %v, want low", got)
	}
}

func TestClassifyKey_None(t *testing.T) {
	if got := ClassifyKey(""); got != SeverityNone {
		t.Errorf("ClassifyKey(\"\") = %v, want none", got)
	}
}

func TestSeverityLevel_String(t *testing.T) {
	cases := []struct {
		level SeverityLevel
		want  string
	}{
		{SeverityNone, "none"},
		{SeverityLow, "low"},
		{SeverityMedium, "medium"},
		{SeverityHigh, "high"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.want {
			t.Errorf("SeverityLevel(%d).String() = %q, want %q", c.level, got, c.want)
		}
	}
}

func TestMaxSeverity_ReturnsHighest(t *testing.T) {
	keys := []string{"log_level", "image", "some_key"}
	if got := MaxSeverity(keys); got != SeverityHigh {
		t.Errorf("MaxSeverity = %v, want high", got)
	}
}

func TestMaxSeverity_EmptyKeys(t *testing.T) {
	if got := MaxSeverity([]string{}); got != SeverityNone {
		t.Errorf("MaxSeverity([]) = %v, want none", got)
	}
}
