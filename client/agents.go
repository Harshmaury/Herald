package client

import (
	"context"
	"fmt"

	accord "github.com/Harshmaury/Accord/api"
)

// AgentsClient provides typed access to the /agents API.
type AgentsClient struct{ c *Client }

// List returns all registered agents with online/offline status.
func (a *AgentsClient) List(ctx context.Context) ([]accord.AgentDTO, error) {
	var out []accord.AgentDTO
	if err := a.c.get(ctx, "/agents", &out); err != nil {
		return nil, fmt.Errorf("agents.list: %w", err)
	}
	return out, nil
}
