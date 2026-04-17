package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Fetcher retrieves live service configuration from a remote endpoint.
type Fetcher struct {
	client  *http.Client
	baseURL string
}

// ServiceConfig holds the raw configuration retrieved from a live service.
type ServiceConfig struct {
	ServiceName string
	Fields      map[string]string
}

// NewFetcher creates a Fetcher targeting the given base URL.
func NewFetcher(baseURL string, timeout time.Duration) *Fetcher {
	return &Fetcher{
		client:  &http.Client{Timeout: timeout},
		baseURL: baseURL,
	}
}

// Fetch retrieves configuration for a named service.
func (f *Fetcher) Fetch(ctx context.Context, serviceName string) (*ServiceConfig, error) {
	url := fmt.Sprintf("%s/services/%s/config", f.baseURL, serviceName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", serviceName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: unexpected status %d", serviceName, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	fields, err := parseConfigBody(body)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &ServiceConfig{ServiceName: serviceName, Fields: fields}, nil
}
