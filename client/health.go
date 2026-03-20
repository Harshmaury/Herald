package client

import (
	"context"
	"fmt"

	accord "github.com/Harshmaury/Accord/api"
)

// HealthClient provides typed access to the /health API.
type HealthClient struct{ c *Client }

// Get returns the current daemon health.
func (h *HealthClient) Get(ctx context.Context) (*accord.HealthData, error) {
	var out accord.HealthData
	if err := h.c.get(ctx, "/health", &out); err != nil {
		return nil, fmt.Errorf("health.get: %w", err)
	}
	return &out, nil
}

// IsReady returns true if the daemon is reachable and healthy.
func (h *HealthClient) IsReady(ctx context.Context) bool {
	_, err := h.Get(ctx)
	return err == nil
}
