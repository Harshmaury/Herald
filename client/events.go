package client

import (
	"context"
	"fmt"
	"strconv"

	accord "github.com/Harshmaury/Accord/api"
)

// EventsClient provides typed access to the /events API.
type EventsClient struct{ c *Client }

// List returns recent platform events.
func (e *EventsClient) List(ctx context.Context, limit int) ([]accord.EventDTO, error) {
	var out []accord.EventDTO
	path := fmt.Sprintf("/events?limit=%d", limit)
	if err := e.c.get(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("events.list: %w", err)
	}
	return out, nil
}

// Since returns events with ID greater than sinceID.
func (e *EventsClient) Since(ctx context.Context, sinceID int64, limit int) ([]accord.EventDTO, error) {
	var out []accord.EventDTO
	path := "/events?since=" + strconv.FormatInt(sinceID, 10) + "&limit=" + strconv.Itoa(limit)
	if err := e.c.get(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("events.since: %w", err)
	}
	return out, nil
}

// ByTrace returns all events for a specific trace ID.
func (e *EventsClient) ByTrace(ctx context.Context, traceID string) ([]accord.EventDTO, error) {
	var out []accord.EventDTO
	if err := e.c.get(ctx, "/events?trace="+traceID, &out); err != nil {
		return nil, fmt.Errorf("events.bytrace: %w", err)
	}
	return out, nil
}
