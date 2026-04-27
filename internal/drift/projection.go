package drift

import (
	"fmt"
	"sort"
	"strings"
)

// ProjectionField defines a single field to include in a projection.
type ProjectionField struct {
	Key   string `json:"key"`
	Alias string `json:"alias,omitempty"`
}

// ProjectionOptions controls which fields are included in projected output.
type ProjectionOptions struct {
	Fields  []ProjectionField
	Service string // optional filter
}

// ProjectedRow is a single projected result row.
type ProjectedRow struct {
	Service string
	Values  map[string]string
}

// ApplyProjection filters CompareResults down to the requested fields.
func ApplyProjection(results []CompareResult, opts ProjectionOptions) []ProjectedRow {
	var rows []ProjectedRow
	for _, r := range results {
		if opts.Service != "" && r.Service != opts.Service {
			continue
		}
		row := ProjectedRow{
			Service: r.Service,
			Values:  make(map[string]string),
		}
		for _, f := range opts.Fields {
			key := f.Key
			alias := f.Alias
			if alias == "" {
				alias = key
			}
			for _, d := range r.Diffs {
				if d.Key == key {
					row.Values[alias] = fmt.Sprintf("%v", d.LiveValue)
					break
				}
			}
			if _, ok := row.Values[alias]; !ok {
				row.Values[alias] = ""
			}
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Service < rows[j].Service
	})
	return rows
}

// FormatProjection renders projected rows as a plain-text table.
func FormatProjection(rows []ProjectedRow, fields []ProjectionField) string {
	if len(rows) == 0 {
		return "no results\n"
	}
	var sb strings.Builder
	headers := []string{"SERVICE"}
	for _, f := range fields {
		alias := f.Alias
		if alias == "" {
			alias = f.Key
		}
		headers = append(headers, strings.ToUpper(alias))
	}
	sb.WriteString(strings.Join(headers, "\t") + "\n")
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, row := range rows {
		cols := []string{row.Service}
		for _, f := range fields {
			alias := f.Alias
			if alias == "" {
				alias = f.Key
			}
			cols = append(cols, row.Values[alias])
		}
		sb.WriteString(strings.Join(cols, "\t") + "\n")
	}
	return sb.String()
}
