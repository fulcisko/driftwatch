package drift_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func TestCompare_NoDrift(t *testing.T) {
	d := drift.NewDetector()
	expected := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}
	deployed := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}

	result := d.Compare("web", expected, deployed)
	if result.HasDrift {
		t.Errorf("expected no drift, got %+v", result.Diffs)
	}
}

func TestCompare_MissingField(t *testing.T) {
	d := drift.NewDetector()
	expected := map[string]interface{}{"replicas": 3, "image": "nginx:1.25"}
	deployed := map[string]interface{}{"replicas": 3}

	result := d.Compare("web", expected, deployed)
	if !result.HasDrift {
		t.Fatal("expected drift but got none")
	}
	if len(result.Diffs) != 1 || result.Diffs[0].Field != "image" {
		t.Errorf("unexpected diffs: %+v", result.Diffs)
	}
}

func TestCompare_ChangedValue(t *testing.T) {
	d := drift.NewDetector()
	expected := map[string]interface{}{"replicas": 3}
	deployed := map[string]interface{}{"replicas": 5}

	result := d.Compare("web", expected, deployed)
	if !result.HasDrift {
		t.Fatal("expected drift but got none")
	}
	diff := result.Diffs[0]
	if diff.Expected != 3 || diff.Actual != 5 {
		t.Errorf("unexpected diff values: %+v", diff)
	}
}

func TestCompare_UnexpectedField(t *testing.T) {
	d := drift.NewDetector()
	expected := map[string]interface{}{"replicas": 3}
	deployed := map[string]interface{}{"replicas": 3, "debug": true}

	result := d.Compare("web", expected, deployed)
	if !result.HasDrift {
		t.Fatal("expected drift for unexpected field")
	}
}
