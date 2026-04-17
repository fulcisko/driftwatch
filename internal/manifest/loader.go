package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Manifest represents a parsed service config manifest.
type Manifest struct {
	Name      string            `yaml:"name"`
	Version   string            `yaml:"version"`
	Namespace string            `yaml:"namespace"`
	Env       map[string]string `yaml:"env"`
	Image     string            `yaml:"image"`
	Replicas  int               `yaml:"replicas"`
}

// LoadFile reads and parses a YAML manifest from the given file path.
func LoadFile(path string) (*Manifest, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("reading manifest %q: %w", abs, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest %q: %w", abs, err)
	}

	if m.Name == "" {
		return nil, fmt.Errorf("manifest %q missing required field: name", abs)
	}

	return &m, nil
}

// LoadDir loads all YAML manifests from a directory (non-recursive).
func LoadDir(dir string) ([]*Manifest, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %q: %w", dir, err)
	}

	var manifests []*Manifest
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		m, err := LoadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, m)
	}
	return manifests, nil
}
