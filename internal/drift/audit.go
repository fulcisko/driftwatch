package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Service   string    `json:"service"`
	User      string    `json:"user"`
	Detail    string    `json:"detail"`
}

func AppendAuditEvent(path, action, service, user, detail string) error {
	events, _ := LoadAuditLog(path)
	events = append(events, AuditEvent{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Service:   service,
		User:      user,
		Detail:    detail,
	})
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal audit log: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadAuditLog(path string) ([]AuditEvent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []AuditEvent{}, nil
		}
		return nil, fmt.Errorf("read audit log: %w", err)
	}
	var events []AuditEvent
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("parse audit log: %w", err)
	}
	return events, nil
}

func FilterAuditLog(events []AuditEvent, service, action string) []AuditEvent {
	var out []AuditEvent
	for _, e := range events {
		if service != "" && e.Service != service {
			continue
		}
		if action != "" && e.Action != action {
			continue
		}
		out = append(out, e)
	}
	return out
}

func FormatAuditLog(events []AuditEvent) string {
	if len(events) == 0 {
		return "no audit events found\n"
	}
	var out string
	for _, e := range events {
		out += fmt.Sprintf("[%s] %s | service=%s user=%s detail=%s\n",
			e.Timestamp.Format(time.RFC3339), e.Action, e.Service, e.User, e.Detail)
	}
	return out
}
