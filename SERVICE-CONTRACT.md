# SERVICE-CONTRACT.md — Herald
# @version: 1.0.0
# @updated: 2026-03-20

**Module:** `github.com/Harshmaury/Herald`
**Type:** Library — HTTP client, no server, no daemon
**Role:** Single entry point for all engxd HTTP API communication
**Depends on:** Accord v0.1.0+

---

## Contract Definition

Herald owns the answer to: *"how do I talk to engxd?"*

No platform component calls engxd directly. All HTTP calls go through herald.
This means: one retry policy, one version check, one error type, one place to fix.

```go
c := client.New("http://127.0.0.1:8080")
svcs, err := c.Services().List(ctx)      // typed, retried, version-checked
proj, err := c.Projects().Get(ctx, "atlas")
ok   := c.Health().IsReady(ctx)
```

---

## Retry Policy (invariant — do not change without ADR)

| Condition | Behaviour |
|---|---|
| `ErrDaemonUnavailable` (connect refused, timeout) | Retry 3×: 0ms → 200ms → 500ms |
| `ErrNotFound` (404) | No retry — resource doesn't exist |
| `ErrUnauthorized` (401) | No retry — token wrong |
| `ErrInvalidInput` (400) | No retry — fix the request |
| Any 4xx | No retry |
| Any 5xx | Retry (treated as transient) |

Retry policy changes require an ADR. Callers depend on timing behaviour.

---

## Error Handling Contract

Herald always returns `*accord.Error` on failure. Callers switch on `Code`:

```go
svcs, err := c.Services().List(ctx)
if err != nil {
    var apiErr *accord.Error
    if errors.As(err, &apiErr) {
        switch apiErr.Code {
        case accord.ErrDaemonUnavailable:
            // engxd not running
        case accord.ErrUnauthorized:
            // bad token
        case accord.ErrNotFound:
            // resource gone
        }
    }
}
```

**Never match on `err.Error()` string content.** Strings are human-readable
and may change. `accord.ErrorCode` values are permanent API surface.

---

## WaitReady Contract

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := client.WaitReady(ctx, "http://127.0.0.1:8080", 500*time.Millisecond)
```

- Polls `GET /health` every `pollInterval` until success or ctx cancelled
- Returns nil when daemon is ready
- Returns `ctx.Err()` wrapped in a descriptive error on timeout
- Replaces all `sleep N` calls before platform start

---

## API Coverage

| Scope | Methods | engxd endpoint |
|---|---|---|
| `c.Services()` | `List`, `Register`, `Reset`, `ResetAll` | `GET/POST /services` |
| `c.Projects()` | `List`, `Get`, `Start`, `Stop`, `Exists` | `GET/POST /projects` |
| `c.Agents()` | `List` | `GET /agents` |
| `c.Events()` | `List`, `Since` | `GET /events` |
| `c.Health()` | `Get`, `IsReady` | `GET /health` |
| `client.WaitReady()` | (package-level) | `GET /health` (poll) |

Socket protocol (`sendCommand`) is NOT part of herald.
Socket commands (project.start, drop.approve, etc.) stay in engx directly.

---

## Consumer Migration Guide

Replace raw `http.Get` calls in engx one command group at a time:

```go
// BEFORE (engx main.go — duplicated 20+ times)
resp, err := http.Get(*httpAddr + "/services")
var result struct { OK bool; Data []struct{...} }
json.Unmarshal(body, &result)

// AFTER (herald)
c := client.New(*httpAddr)
svcs, err := c.Services().List(ctx)
```

Migration order:
1. `platform` commands (highest value — fixes desired=running bug)
2. `doctor` command (removes 5 raw http.Get calls)
3. `services` / `agents` / `events` commands
4. `build` / `check` / `trace` commands

---

## What Herald Is NOT

- Not a daemon — no goroutines, no background polling
- Not a socket client — daemon socket protocol stays in engx
- Not aware of business logic — thin HTTP adapter only
- Not a replacement for the SSE stream client (that is cmd_follow.go)

---

## Breakage Prevention Checklist

Before merging any PR to herald:

- [ ] accord version in go.mod is pinned to a released tag (no pseudo-versions)
- [ ] All new methods have context.Context as first parameter
- [ ] All new methods return typed accord DTOs, not map[string]any
- [ ] Retry policy not changed without ADR
- [ ] Error wrapping preserves *accord.Error type (errors.As compatible)
- [ ] No business logic added — herald is a transport layer only
