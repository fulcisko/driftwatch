package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func ttlPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ttl.json")
}

func TestAddAndLoadTTLRule(t *testing.T) {
	p := ttlPath(t)
	if err := AddTTLRule(p, "svc-a", 2*time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, err := LoadTTLList(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(list.Rules) != 1 || list.Rules[0].Service != "svc-a" {
		t.Errorf("expected svc-a, got %+v", list.Rules)
	}
}

func TestAddTTLRule_UpdatesExisting(t *testing.T) {
	p := ttlPath(t)
	_ = AddTTLRule(p, "svc-a", time.Hour)
	_ = AddTTLRule(p, "svc-a", 3*time.Hour)
	list, _ := LoadTTLList(p)
	if len(list.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(list.Rules))
	}
	if list.Rules[0].TTL != 3*time.Hour {
		t.Errorf("expected updated TTL 3h, got %v", list.Rules[0].TTL)
	}
}

func TestAddTTLRule_MissingService(t *testing.T) {
	p := ttlPath(t)
	if err := AddTTLRule(p, "", time.Hour); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestAddTTLRule_InvalidTTL(t *testing.T) {
	p := ttlPath(t)
	if err := AddTTLRule(p, "svc-a", 0); err == nil {
		t.Error("expected error for zero ttl")
	}
}

func TestLoadTTLList_NotFound(t *testing.T) {
	list, err := LoadTTLList("/nonexistent/ttl.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(list.Rules) != 0 {
		t.Errorf("expected empty list")
	}
}

func TestExpiredServices_DetectsExpired(t *testing.T) {
	p := ttlPath(t)
	_ = AddTTLRule(p, "svc-old", time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_ = AddTTLRule(p, "svc-new", 24*time.Hour)
	list, _ := LoadTTLList(p)
	expired := ExpiredServices(list)
	if len(expired) != 1 || expired[0] != "svc-old" {
		t.Errorf("expected [svc-old], got %v", expired)
	}
}

func TestExpiredServices_NoneExpired(t *testing.T) {
	list := TTLList{
		Rules: []TTLRule{
			{Service: "svc-a", TTL: 24 * time.Hour, CreatedAt: time.Now().UTC()},
		},
	}
	if exp := ExpiredServices(list); len(exp) != 0 {
		t.Errorf("expected no expired services, got %v", exp)
	}
}

func TestSaveTTLList_Persists(t *testing.T) {
	p := filepath.Join(t.TempDir(), "ttl.json")
	list := TTLList{
		Rules: []TTLRule{
			{Service: "svc-x", TTL: time.Hour, CreatedAt: time.Now().UTC()},
		},
	}
	if err := saveTTLList(p, list); err != nil {
		t.Fatalf("save error: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Errorf("file not written: %v", err)
	}
}
