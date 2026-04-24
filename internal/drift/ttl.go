package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TTLRule defines a time-to-live rule for a service's drift tolerance.
type TTLRule struct {
	Service   string        `json:"service"`
	TTL       time.Duration `json:"ttl_ns"` // stored as nanoseconds
	CreatedAt time.Time     `json:"created_at"`
}

// TTLList holds all TTL rules.
type TTLList struct {
	Rules []TTLRule `json:"rules"`
}

// AddTTLRule adds or updates a TTL rule for the given service.
func AddTTLRule(path, service string, ttl time.Duration) error {
	if service == "" {
		return fmt.Errorf("service name is required")
	}
	if ttl <= 0 {
		return fmt.Errorf("ttl must be positive")
	}
	list, _ := LoadTTLList(path)
	updated := false
	for i, r := range list.Rules {
		if r.Service == service {
			list.Rules[i].TTL = ttl
			list.Rules[i].CreatedAt = time.Now().UTC()
			updated = true
			break
		}
	}
	if !updated {
		list.Rules = append(list.Rules, TTLRule{
			Service:   service,
			TTL:       ttl,
			CreatedAt: time.Now().UTC(),
		})
	}
	return saveTTLList(path, list)
}

// LoadTTLList loads the TTL rules from disk.
func LoadTTLList(path string) (TTLList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TTLList{}, nil
		}
		return TTLList{}, err
	}
	var list TTLList
	if err := json.Unmarshal(data, &list); err != nil {
		return TTLList{}, err
	}
	return list, nil
}

// ExpiredServices returns the names of services whose TTL has elapsed
// since the rule was created.
func ExpiredServices(list TTLList) []string {
	now := time.Now().UTC()
	var expired []string
	for _, r := range list.Rules {
		if now.After(r.CreatedAt.Add(r.TTL)) {
			expired = append(expired, r.Service)
		}
	}
	return expired
}

func saveTTLList(path string, list TTLList) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
