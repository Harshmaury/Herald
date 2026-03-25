// @herald-project: herald
// @herald-path: SERVICE-CONTRACT.md
# SERVICE-CONTRACT.md — Herald
# @version: 0.1.5
# @updated: 2026-03-25

**Type:** Library · **Module:** `github.com/Harshmaury/Herald` · **Domain:** Shared protocol

---

## Code

```
client/client.go      New(), NewForService() — constructor with token + address
client/services.go    Services().List/Register/Reset/ResetAll
client/projects.go    Projects().List/Get/Start/Stop/Exists
client/agents.go      Agents().List
client/events.go      Events().List/Since
client/health.go      Health().Get/IsReady · WaitReady(ctx, addr, pollInterval)
client/atlas.go       Atlas().Projects/Graph/Services
client/forge.go       Forge().History
client/guardian.go    Guardian().Findings
client/metrics.go     NexusMetrics().Get
client/navigator.go   Navigator().Graph
```

---

## Contract

No HTTP server. All methods accept `context.Context` as first param. All return typed Accord DTOs.

**Retry policy (invariant — changes require ADR):**

| Condition | Behavior |
|-----------|----------|
| `ErrDaemonUnavailable` | Retry 3×: 0ms → 200ms → 500ms |
| 4xx (any) | No retry |
| 5xx | Retry (transient) |

**Error handling:** always `*accord.Error`. Switch on `Code`, never on `err.Error()` string.

**`WaitReady`:** polls `GET /health` every `pollInterval` until success or ctx cancelled. Returns `ctx.Err()` on timeout. Replaces all `sleep N` before platform start.

**Usage:**
```go
// Nexus
c := herald.New(nexusAddr, herald.WithToken(token))
svcs, err := c.Services().List(ctx)
evts, err := c.Events().Since(ctx, sinceID, 100)

// Non-Nexus upstreams
ac := herald.NewForService(atlasAddr, serviceToken)
projs, err := ac.Atlas().Projects(ctx)
```

---

## Control

No runtime behavior. No goroutines, no background polling. Thin HTTP adapter only.

---

## Context

Single entry point for all inter-service HTTP communication. No platform component calls any service via raw `http.Get` or `http.Post`. Retry policy and error wrapping are centralised here.
