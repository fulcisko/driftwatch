package drift

// DiffKind describes the type of a drift difference.
type DiffKind string

const (
	KindChanged    DiffKind = "changed"
	KindMissing    DiffKind = "missing"
	KindUnexpected DiffKind = "unexpected"
)

// Diff represents a single field-level difference between manifest and live config.
type Diff struct {
	Key      string   `json:"key"`
	Expected string   `json:"expected"`
	Actual   string   `json:"actual"`
	Kind     DiffKind `json:"kind"`
}

// CompareResult holds the drift results for a single service.
type CompareResult struct {
	Service string `json:"service"`
	Diffs   []Diff `json:"diffs"`
}

// HasDrift returns true when the result contains at least one diff.
func (r CompareResult) HasDrift() bool {
	return len(r.Diffs) > 0
}
