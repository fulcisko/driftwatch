package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type NotifyChannel string

const (
	ChannelSlack  NotifyChannel = "slack"
	ChannelEmail  NotifyChannel = "email"
	ChannelWebhook NotifyChannel = "webhook"
)

type NotifyRule struct {
	Channel   NotifyChannel `json:"channel"`
	Target    string        `json:"target"`
	MinSeverity string      `json:"min_severity"`
	Services  []string      `json:"services,omitempty"`
}

type NotifyEvent struct {
	Service   string        `json:"service"`
	Severity  string        `json:"severity"`
	Message   string        `json:"message"`
	Timestamp time.Time     `json:"timestamp"`
	Channel   NotifyChannel `json:"channel"`
	Target    string        `json:"target"`
}

func GenerateNotifyEvents(results []CompareResult, rules []NotifyRule) []NotifyEvent {
	var events []NotifyEvent
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		sev := string(MaxSeverity(r.Diffs))
		for _, rule := range rules {
			if !severityMeetsMin(sev, rule.MinSeverity) {
				continue
			}
			if len(rule.Services) > 0 && !containsService(rule.Services, r.Service) {
				continue
			}
			events = append(events, NotifyEvent{
				Service:   r.Service,
				Severity:  sev,
				Message:   fmt.Sprintf("drift detected in %s: %d diff(s)", r.Service, len(r.Diffs)),
				Timestamp: time.Now().UTC(),
				Channel:   rule.Channel,
				Target:    rule.Target,
			})
		}
	}
	return events
}

func SaveNotifyEvents(path string, events []NotifyEvent) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(events)
}

func LoadNotifyEvents(path string) ([]NotifyEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var events []NotifyEvent
	return events, json.NewDecoder(f).Decode(&events)
}

func severityMeetsMin(sev, min string) bool {
	order := map[string]int{"none": 0, "low": 1, "medium": 2, "high": 3}
	return order[sev] >= order[min]
}

func containsService(list []string, svc string) bool {
	for _, s := range list {
		if s == svc {
			return true
		}
	}
	return false
}
