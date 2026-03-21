// @herald-project: Herald
// @herald-path: client/gate.go
// GateClient provides typed access to the Gate identity API (ADR-042).
// Gate is the sole identity authority — never call it directly from service code.
//
// Usage:
//   c := client.NewForService("http://127.0.0.1:8088", serviceToken)
//   claim, err := c.Gate().Validate(ctx, identityToken)
//   if err != nil || !claim.Valid { /* deny */ }
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	accord "github.com/Harshmaury/Accord/api"
)

// GateClient provides typed access to the Gate /gate API.
type GateClient struct{ c *Client }

// Validate calls POST /gate/validate with the given identity token.
// Returns the validated claim on success, or an error if the token is
// invalid, expired, or revoked.
// Callers should use the returned claim's HasScope method for authorization.
func (g *GateClient) Validate(ctx context.Context, token string) (*accord.GateValidateResponse, error) {
	body, err := json.Marshal(accord.GateValidateRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("gate.validate: marshal: %w", err)
	}
	var out accord.GateValidateResponse
	if err := g.c.post(ctx, "/gate/validate", bytes.NewReader(body), &out); err != nil {
		return nil, fmt.Errorf("gate.validate: %w", err)
	}
	return &out, nil
}

// PublicKey fetches Gate's Ed25519 public key from GET /gate/public-key.
// Services call this at startup to cache the key for local token verification.
// Re-call if local signature verification fails — handles key rotation.
func (g *GateClient) PublicKey(ctx context.Context) (*accord.GatePublicKeyDTO, error) {
	var out accord.GatePublicKeyDTO
	if err := g.c.get(ctx, "/gate/public-key", &out); err != nil {
		return nil, fmt.Errorf("gate.public-key: %w", err)
	}
	return &out, nil
}

// RevokeToken calls POST /gate/revoke to revoke a token by JTI.
// Requires a Gate client token with ScopeAdmin.
func (g *GateClient) RevokeToken(ctx context.Context, jti string) error {
	body, err := json.Marshal(map[string]string{"jti": jti})
	if err != nil {
		return fmt.Errorf("gate.revoke: marshal: %w", err)
	}
	var out struct{}
	if err := g.c.post(ctx, "/gate/revoke", bytes.NewReader(body), &out); err != nil {
		return fmt.Errorf("gate.revoke: %w", err)
	}
	return nil
}
