# SERVICE-CONTRACT.md — Herald

**Type:** Library (no runtime, no HTTP server)
**Module:** github.com/Harshmaury/Herald
**Version:** 0.1.0
**Depends on:** Accord v0.1.0

---

## What Herald Is

Herald is the typed HTTP client for the Nexus API.
Every component that needs to talk to engxd imports herald.

```go
import "github.com/Harshmaury/Herald/client"

c := client.New("http://127.0.0.1:8080")
svcs, err := c.Services().List(ctx)
projects, err := c.Projects().List(ctx)
health, err := c.Health().Get(ctx)
```

---

## Features

- Typed request/response via Accord — compile-time contract with engxd
- Retry with exponential backoff on ErrDaemonUnavailable (3 attempts: 0/200ms/500ms)
- No retry on 4xx errors — bad input should not be retried
- WaitReady(ctx, addr, interval) — polls /health until daemon is up
- X-Service-Token header injection via WithToken() option
- API version header check on every response

---

## Consumer Map

| Consumer | Uses herald for |
|---|---|
| engx (CLI) | services, projects, agents, events, health |
| Atlas | project registration check |
| Forge | project start/stop, event streaming |
| ZP | health check only |
| External tools | any — herald is the public SDK |

---

## What Herald Is NOT

- Not a daemon — no goroutines, no background work
- Not a socket client — socket protocol stays in engx/main.go (daemon commands)
- Not aware of business logic — thin adapter only
