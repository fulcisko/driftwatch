package source

import (
	"testing"
)

func TestParseConfigBody_Valid(t *testing.T) {
	input := []byte("# comment\nKEY=value\nFOO=bar baz\n")
	fields, err := parseConfigBody(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", fields["KEY"])
	}
	if fields["FOO"] != "bar baz" {
		t.Errorf("expected FOO='bar baz', got %q", fields["FOO"])
	}
}

func TestParseConfigBody_Empty(t *testing.T) {
	fields, err := parseConfigBody([]byte("\n# only comments\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("expected empty map, got %v", fields)
	}
}

func TestParseConfigBody_InvalidLine(t *testing.T) {
	_, err := parseConfigBody([]byte("NOEQUALSSIGN\n"))
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestParseConfigBody_EmptyKey(t *testing.T) {
	_, err := parseConfigBody([]byte("=value\n"))
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestParseConfigBody_ValueWithEquals(t *testing.T) {
	fields, err := parseConfigBody([]byte("URL=http://example.com?a=1\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["URL"] != "http://example.com?a=1" {
		t.Errorf("unexpected URL value: %q", fields["URL"])
	}
}
