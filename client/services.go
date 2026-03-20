package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	accord "github.com/Harshmaury/Accord/api"
)

// ServicesClient provides typed access to the /services API.
type ServicesClient struct{ c *Client }

// List returns all registered services.
func (s *ServicesClient) List(ctx context.Context) ([]accord.ServiceDTO, error) {
	var out []accord.ServiceDTO
	if err := s.c.get(ctx, "/services", &out); err != nil {
		return nil, fmt.Errorf("services.list: %w", err)
	}
	return out, nil
}

// Reset clears the fail count and maintenance state for a service.
func (s *ServicesClient) Reset(ctx context.Context, id string) (*accord.ServiceResetData, error) {
	var out accord.ServiceResetData
	if err := s.c.post(ctx, "/services/"+id+"/reset", nil, &out); err != nil {
		return nil, fmt.Errorf("services.reset %s: %w", id, err)
	}
	return &out, nil
}

// Register creates or updates a service in the Nexus state store.
func (s *ServicesClient) Register(ctx context.Context, req accord.ServiceRegisterRequest) (*accord.ServiceIdentity, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("services.register: marshal: %w", err)
	}
	var out accord.ServiceIdentity
	if err := s.c.post(ctx, "/services/register", bytes.NewReader(body), &out); err != nil {
		return nil, fmt.Errorf("services.register: %w", err)
	}
	return &out, nil
}

// ResetAll resets all given service IDs. Best-effort — errors are collected, not fatal.
func (s *ServicesClient) ResetAll(ctx context.Context, ids []string) map[string]error {
	errs := make(map[string]error)
	for _, id := range ids {
		if _, err := s.Reset(ctx, id); err != nil {
			errs[id] = err
		}
	}
	return errs
}
