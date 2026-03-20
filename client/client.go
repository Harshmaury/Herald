// Package client provides a typed HTTP client for the Nexus API.
// Import herald — never call engxd directly from application code.
//
// Usage:
//   c := client.New("http://127.0.0.1:8080")
//   svcs, err := c.Services().List(ctx)
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	accord "github.com/Harshmaury/Accord/api"
)

const (
	defaultTimeout    = 10 * time.Second
	defaultMaxRetries = 3
)

// Client is the Nexus API HTTP client. Safe for concurrent use.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	maxRetries int
}

// Option configures a Client.
type Option func(*Client)

// WithToken sets the X-Service-Token header on every request.
func WithToken(token string) Option { return func(c *Client) { c.token = token } }

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithRetries sets the number of retry attempts on transient failures.
func WithRetries(n int) Option { return func(c *Client) { c.maxRetries = n } }

// New creates a Client for the Nexus API at baseURL.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		maxRetries: defaultMaxRetries,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Services returns the services API client.
func (c *Client) Services() *ServicesClient { return &ServicesClient{c} }

// Projects returns the projects API client.
func (c *Client) Projects() *ProjectsClient { return &ProjectsClient{c} }

// Agents returns the agents API client.
func (c *Client) Agents() *AgentsClient { return &AgentsClient{c} }

// Events returns the events API client.
func (c *Client) Events() *EventsClient { return &EventsClient{c} }

// Health returns the health API client.
func (c *Client) Health() *HealthClient { return &HealthClient{c} }

// WaitReady polls GET /health until engxd responds or ctx is cancelled.
func WaitReady(ctx context.Context, baseURL string, pollInterval time.Duration) error {
	c := New(baseURL, WithRetries(1), WithTimeout(2*time.Second))
	for {
		var h accord.HealthData
		if err := c.get(ctx, "/health", &h); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("engxd not ready at %s: %w", baseURL, ctx.Err())
		case <-time.After(pollInterval):
		}
	}
}

// ── HTTP PRIMITIVES ──────────────────────────────────────────────────────────

func (c *Client) get(ctx context.Context, path string, out any) error {
	return c.doWithRetry(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) post(ctx context.Context, path string, body io.Reader, out any) error {
	return c.doWithRetry(ctx, http.MethodPost, path, body, out)
}

func (c *Client) doWithRetry(ctx context.Context, method, path string, body io.Reader, out any) error {
	backoff := []time.Duration{0, 200 * time.Millisecond, 500 * time.Millisecond}
	retries := c.maxRetries
	if retries > len(backoff) {
		retries = len(backoff)
	}
	var lastErr error
	for i := 0; i < retries; i++ {
		if backoff[i] > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff[i]):
			}
		}
		err := c.do(ctx, method, path, body, out)
		if err == nil {
			return nil
		}
		if e, ok := err.(*accord.Error); ok && e.Code != accord.ErrDaemonUnavailable {
			return err // 4xx — do not retry
		}
		lastErr = err
	}
	return fmt.Errorf("after %d retries: %w", retries, lastErr)
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("X-Service-Token", c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &accord.Error{Code: accord.ErrDaemonUnavailable,
			Message: fmt.Sprintf("cannot reach engxd at %s: %v", c.baseURL, err)}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	var env struct {
		OK    bool            `json:"ok"`
		Error string          `json:"error,omitempty"`
		Data  json.RawMessage `json:"data,omitempty"`
	}
	if err := json.Unmarshal(data, &env); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if !env.OK {
		code := accord.ErrInternal
		switch resp.StatusCode {
		case http.StatusNotFound:
			code = accord.ErrNotFound
		case http.StatusUnauthorized:
			code = accord.ErrUnauthorized
		case http.StatusBadRequest:
			code = accord.ErrInvalidInput
		}
		return &accord.Error{Code: code, Message: env.Error}
	}
	if out != nil && len(env.Data) > 0 {
		return json.Unmarshal(env.Data, out)
	}
	return nil
}
