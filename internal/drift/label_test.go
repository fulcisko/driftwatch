package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func labelPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "labels.json")
}

func TestAddAndLoadLabel(t *testing.T) {
	p := labelPath(t)
	if err := AddLabel(p, "svc-a", "env", "prod"); err != nil {
		t.Fatal(err)
	}
	labels, err := LoadLabels(p)
	if err != nil {
		t.Fatal(err)
	}
	if labels["svc-a"]["env"] != "prod" {
		t.Errorf("expected prod, got %s", labels["svc-a"]["env"])
	}
}

func TestAddLabel_MissingFields(t *testing.T) {
	p := labelPath(t)
	if err := AddLabel(p, "", "env", "prod"); err == nil {
		t.Error("expected error for empty service")
	}
	if err := AddLabel(p, "svc", "", "prod"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestRemoveLabel_Success(t *testing.T) {
	p := labelPath(t)
	_ = AddLabel(p, "svc-a", "env", "prod")
	if err := RemoveLabel(p, "svc-a", "env"); err != nil {
		t.Fatal(err)
	}
	labels, _ := LoadLabels(p)
	if _, ok := labels["svc-a"]; ok {
		t.Error("expected service to be removed")
	}
}

func TestRemoveLabel_NotFound(t *testing.T) {
	p := labelPath(t)
	if err := RemoveLabel(p, "ghost", "env"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestLoadLabels_NotFound(t *testing.T) {
	labels, err := LoadLabels("/nonexistent/labels.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 0 {
		t.Error("expected empty labels")
	}
}

func TestFilterByLabel(t *testing.T) {
	p := labelPath(t)
	_ = AddLabel(p, "svc-a", "team", "platform")
	_ = AddLabel(p, "svc-b", "team", "data")
	labels, _ := LoadLabels(p)
	results := []CompareResult{
		{Service: "svc-a"},
		{Service: "svc-b"},
		{Service: "svc-c"},
	}
	filtered := FilterByLabel(results, labels, "team", "platform")
	if len(filtered) != 1 || filtered[0].Service != "svc-a" {
		t.Errorf("unexpected filter result: %+v", filtered)
	}
}

func TestAddLabel_Overwrite(t *testing.T) {
	p := labelPath(t)
	_ = AddLabel(p, "svc-a", "env", "staging")
	_ = AddLabel(p, "svc-a", "env", "prod")
	labels, _ := LoadLabels(p)
	if labels["svc-a"]["env"] != "prod" {
		t.Errorf("expected prod, got %s", labels["svc-a"]["env"])
	}
	os.Remove(p)
}
