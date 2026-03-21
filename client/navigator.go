// @herald-project: Herald
// @herald-path: client/navigator.go
// NavigatorClient provides typed access to the Navigator topology API.
// ADR-039 gap closure: Guardian navigator.go raw HTTP replaced by this client.
// Navigator speaks the standard {ok, data} envelope — Herald handles it uniformly.
package client

import (
	"context"
	"fmt"

	accord "github.com/Harshmaury/Accord/api"
)

// NavigatorClient provides typed access to the Navigator /topology API.
type NavigatorClient struct{ c *Client }

// Graph returns the full workspace topology from Navigator GET /topology/graph.
func (n *NavigatorClient) Graph(ctx context.Context) (*accord.NavigatorGraphDTO, error) {
	var out accord.NavigatorGraphDTO
	if err := n.c.get(ctx, "/topology/graph", &out); err != nil {
		return nil, fmt.Errorf("navigator.graph: %w", err)
	}
	return &out, nil
}
