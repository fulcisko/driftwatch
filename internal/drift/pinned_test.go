package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func pinPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestAddAndLoadPin(t *testing.T) {
	p := pinPath(t)
	if err := AddPin(p, "svc-a", "LOG_LEVEL", "info", "stable"); err != nil {
		t.Fatal(err)
	}
	list, err := LoadPins(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(list.Pins))
	}
	if list.Pins[0].Expected != "info" {
		t.Errorf("unexpected expected value: %s", list.Pins[0].Expected)
	}
}

func TestAddPin_UpdatesExisting(t *testing.T) {
	p := pinPath(t)
	_ = AddPin(p, "svc-a", "LOG_LEVEL", "info", "")
	_ = AddPin(p, "svc-a", "LOG_LEVEL", "debug", "updated")
	list, _ := LoadPins(p)
	if len(list.Pins) != 1 {
		t.Fatalf("expected 1 pin after update, got %d", len(list.Pins))
	}
	if list.Pins[0].Expected != "debug" {
		t.Errorf("expected updated value debug, got %s", list.Pins[0].Expected)
	}
}

func TestRemovePin_Success(t *testing.T) {
	p := pinPath(t)
	_ = AddPin(p, "svc-a", "LOG_LEVEL", "info", "")
	if err := RemovePin(p, "svc-a", "LOG_LEVEL"); err != nil {
		t.Fatal(err)
	}
	list, _ := LoadPins(p)
	if len(list.Pins) != 0 {
		t.Errorf("expected 0 pins after removal")
	}
}

func TestRemovePin_NotFound(t *testing.T) {
	p := pinPath(t)
	err := RemovePin(p, "svc-x", "MISSING")
	if err == nil {
		t.Error("expected error for missing pin")
	}
}

func TestLoadPins_NotFound(t *testing.T) {
	list, err := LoadPins("/tmp/nonexistent_pins_xyz.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Pins) != 0 {
		t.Error("expected empty list")
	}
}

func TestApplyPins_FiltersPinnedMatch(t *testing.T) {
	_ = os.Setenv("TEST_APPLY_PINS", "1")
	list := PinList{
		Pins: []PinnedKey{
			{Service: "svc-a", Key: "LOG_LEVEL", Expected: "info"},
		},
	}
	results := []CompareResult{
		{Service: "svc-a", Diffs: []DiffEntry{
			{Key: "LOG_LEVEL", LiveValue: "info"},
			{Key: "TIMEOUT", LiveValue: "30"},
		}},
	}
	out := ApplyPins(results, list)
	if len(out[0].Diffs) != 1 {
		t.Errorf("expected 1 diff after pin filter, got %d", len(out[0].Diffs))
	}
	if out[0].Diffs[0].Key != "TIMEOUT" {
		t.Errorf("expected TIMEOUT diff to remain")
	}
}
