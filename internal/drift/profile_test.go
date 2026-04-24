package drift

import (
	"path/filepath"
	"testing"
)

func profilePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "profiles.json")
}

func TestSaveAndLoadProfile(t *testing.T) {
	path := profilePath(t)
	p := Profile{
		Name:        "strict",
		Description: "High severity only",
		MinSeverity: "high",
		IgnoreKeys:  []string{"version"},
	}
	if err := SaveProfile(path, p); err != nil {
		t.Fatalf("SaveProfile: %v", err)
	}
	profiles, err := LoadProfiles(path)
	if err != nil {
		t.Fatalf("LoadProfiles: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	if profiles[0].Name != "strict" {
		t.Errorf("expected name 'strict', got %q", profiles[0].Name)
	}
	if profiles[0].MinSeverity != "high" {
		t.Errorf("expected min_severity 'high', got %q", profiles[0].MinSeverity)
	}
}

func TestSaveProfile_MissingName(t *testing.T) {
	path := profilePath(t)
	err := SaveProfile(path, Profile{Description: "no name"})
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestSaveProfile_UpdatesExisting(t *testing.T) {
	path := profilePath(t)
	p := Profile{Name: "dev", MinSeverity: "low"}
	_ = SaveProfile(path, p)
	p.MinSeverity = "medium"
	p.Description = "updated"
	if err := SaveProfile(path, p); err != nil {
		t.Fatalf("SaveProfile update: %v", err)
	}
	profiles, _ := LoadProfiles(path)
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile after update, got %d", len(profiles))
	}
	if profiles[0].MinSeverity != "medium" {
		t.Errorf("expected updated min_severity 'medium', got %q", profiles[0].MinSeverity)
	}
}

func TestGetProfile_Found(t *testing.T) {
	path := profilePath(t)
	_ = SaveProfile(path, Profile{Name: "prod", MinSeverity: "high"})
	p, ok := GetProfile(path, "prod")
	if !ok {
		t.Fatal("expected profile to be found")
	}
	if p.Name != "prod" {
		t.Errorf("expected 'prod', got %q", p.Name)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	path := profilePath(t)
	_, ok := GetProfile(path, "missing")
	if ok {
		t.Error("expected profile not to be found")
	}
}

func TestRemoveProfile_Success(t *testing.T) {
	path := profilePath(t)
	_ = SaveProfile(path, Profile{Name: "temp", MinSeverity: "low"})
	if err := RemoveProfile(path, "temp"); err != nil {
		t.Fatalf("RemoveProfile: %v", err)
	}
	profiles, _ := LoadProfiles(path)
	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles after removal, got %d", len(profiles))
	}
}

func TestRemoveProfile_NotFound(t *testing.T) {
	path := profilePath(t)
	err := RemoveProfile(path, "ghost")
	if err == nil {
		t.Error("expected error when removing non-existent profile")
	}
}

func TestLoadProfiles_NotFound(t *testing.T) {
	path := profilePath(t) + ".missing"
	profiles, err := LoadProfiles(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected empty profiles, got %d", len(profiles))
	}
}
