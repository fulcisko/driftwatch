package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Owner represents a team or individual responsible for a service.
type Owner struct {
	Service string `json:"service"`
	Team    string `json:"team"`
	Contact string `json:"contact"`
}

// OwnershipMap maps service names to their owners.
type OwnershipMap struct {
	Owners []Owner `json:"owners"`
}

// AddOwner adds or updates an owner entry for a service.
func AddOwner(path, service, team, contact string) error {
	if service == "" || team == "" {
		return fmt.Errorf("service and team are required")
	}
	om, _ := LoadOwnership(path)
	for i, o := range om.Owners {
		if o.Service == service {
			om.Owners[i] = Owner{Service: service, Team: team, Contact: contact}
			return saveOwnership(path, om)
		}
	}
	om.Owners = append(om.Owners, Owner{Service: service, Team: team, Contact: contact})
	return saveOwnership(path, om)
}

// RemoveOwner removes the ownership entry for a service.
func RemoveOwner(path, service string) error {
	om, err := LoadOwnership(path)
	if err != nil {
		return err
	}
	filtered := om.Owners[:0]
	for _, o := range om.Owners {
		if o.Service != service {
			filtered = append(filtered, o)
		}
	}
	if len(filtered) == len(om.Owners) {
		return fmt.Errorf("service %q not found in ownership map", service)
	}
	om.Owners = filtered
	return saveOwnership(path, om)
}

// LoadOwnership loads the ownership map from disk.
func LoadOwnership(path string) (OwnershipMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return OwnershipMap{}, nil
	}
	var om OwnershipMap
	if err := json.Unmarshal(data, &om); err != nil {
		return OwnershipMap{}, fmt.Errorf("parse ownership: %w", err)
	}
	return om, nil
}

// LookupOwner returns the owner for a given service, if any.
func LookupOwner(om OwnershipMap, service string) (Owner, bool) {
	for _, o := range om.Owners {
		if strings.EqualFold(o.Service, service) {
			return o, true
		}
	}
	return Owner{}, false
}

func saveOwnership(path string, om OwnershipMap) error {
	data, err := json.MarshalIndent(om, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
