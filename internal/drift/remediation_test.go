package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func remPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "remediation.json")
}

func TestAddAndLoadRemediation(t *testing.T) {
	p := remPath(t)
	if err := AddRemediation(p, "svc-a", "replicas", ActionApply, "auto-fix"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	log, err := LoadRemediations(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Service != "svc-a" || e.Key != "replicas" || e.Action != ActionApply {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestLoadRemediations_NotFound(t *testing.T) {
	log, err := LoadRemediations("/tmp/nonexistent_rem.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestFilterRemediations_ByService(t *testing.T) {
	p := remPath(t)
	_ = AddRemediation(p, "svc-a", "cpu", ActionIgnore, "")
	_ = AddRemediation(p, "svc-b", "mem", ActionRevert, "")
	_ = AddRemediation(p, "svc-a", "mem", ActionApply, "")
	log, _ := LoadRemediations(p)
	res := FilterRemediations(log, "svc-a")
	if len(res) != 2 {
		t.Errorf("expected 2, got %d", len(res))
	}
}

func TestFilterRemediations_AllServices(t *testing.T) {
	p := remPath(t)
	_ = AddRemediation(p, "svc-a", "cpu", ActionIgnore, "")
	_ = AddRemediation(p, "svc-b", "mem", ActionRevert, "")
	log, _ := LoadRemediations(p)
	res := FilterRemediations(log, "")
	if len(res) != 2 {
		t.Errorf("expected 2, got %d", len(res))
	}
}

func TestAddRemediation_MultipleEntries(t *testing.T) {
	p := remPath(t)
	for i := 0; i < 3; i++ {
		if err := AddRemediation(p, "svc", "key", ActionApply, ""); err != nil {
			t.Fatal(err)
		}
	}
	log, _ := LoadRemediations(p)
	if len(log.Entries) != 3 {
		t.Errorf("expected 3, got %d", len(log.Entries))
	}
	_ = os.Remove(p)
}
