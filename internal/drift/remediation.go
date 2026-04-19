package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type RemediationAction string

const (
	ActionApply  RemediationAction = "apply"
	ActionRevert RemediationAction = "revert"
	ActionIgnore RemediationAction = "ignore"
)

type RemediationEntry struct {
	Service   string            `json:"service"`
	Key       string            `json:"key"`
	Action    RemediationAction `json:"action"`
	Note      string            `json:"note,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type RemediationLog struct {
	Entries []RemediationEntry `json:"entries"`
}

func AddRemediation(path, service, key string, action RemediationAction, note string) error {
	log, _ := LoadRemediations(path)
	log.Entries = append(log.Entries, RemediationEntry{
		Service:   service,
		Key:       key,
		Action:    action,
		Note:      note,
		CreatedAt: time.Now().UTC(),
	})
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal remediation: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadRemediations(path string) (RemediationLog, error) {
	var log RemediationLog
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return log, nil
	}
	if err != nil {
		return log, fmt.Errorf("read remediation log: %w", err)
	}
	if err := json.Unmarshal(data, &log); err != nil {
		return log, fmt.Errorf("parse remediation log: %w", err)
	}
	return log, nil
}

func FilterRemediations(log RemediationLog, service string) []RemediationEntry {
	var out []RemediationEntry
	for _, e := range log.Entries {
		if service == "" || e.Service == service {
			out = append(out, e)
		}
	}
	return out
}
