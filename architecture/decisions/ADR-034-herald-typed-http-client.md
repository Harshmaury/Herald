# ADR-034 — Herald: Typed HTTP Client for the Nexus API

**Status:** Accepted
**Date:** 2026-03-20
**Author:** Harsh Maury
**Scope:** New project — github.com/Harshmaury/Herald
**Depends on:** ADR-033 (Accord), ADR-003 (HTTP/JSON)

---

## Context

Every platform component that talks to engxd makes raw `http.Get`/`http.Post`
calls with inline JSON structs. This is duplicated ~20 times in cmd/engx/main.go
alone. No retry logic, no API version check, no structured error handling.

---

## Decision

Create `github.com/Harshmaury/Herald` — typed HTTP client for the Nexus API.

```go
c := client.New("http://127.0.0.1:8080", client.WithToken(token))
svcs, err := c.Services().List(ctx)
proj, err := c.Projects().Get(ctx, id)
ok   := c.Health().IsReady(ctx)

// Error handling — switch on code, not string
if e, ok := err.(*accord.Error); ok && e.Code == accord.ErrNotFound {
    // handle cleanly
}
```

### Retry Policy
- 3 attempts: 0ms / 200ms / 500ms backoff
- Retry only on ErrDaemonUnavailable
- Never retry on 4xx

### WaitReady
```go
client.WaitReady(ctx, "http://127.0.0.1:8080", 500*time.Millisecond)
```
Replaces manual sleep loops in platform start.

### Migration path for engx — replace raw http.Get calls per command group:
1. platform commands
2. doctor command
3. services/agents/events commands
4. build/check/trace commands

sendCommand (socket protocol) is NOT replaced — stays for daemon commands.

---

## Compliance

| ADR | Status |
|-----|--------|
| ADR-003 | ✅ HTTP/JSON on 127.0.0.1 — client enforces no cross-service calls |
| ADR-033 | ✅ Uses Accord types throughout |
