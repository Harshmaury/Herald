# WORKFLOW SESSION — ENGX-HERALD-P1-001

**Date:** 2026-03-21
**Repos:** Herald + Accord (apply to both)
**Blocks:** All ADR-039 service ZIPs — apply this first

## What changed

**Accord — api/upstream.go (new file)**
Typed DTOs for Atlas, Forge, Guardian, and Nexus /metrics.
All verified against source handler output before writing.
- AtlasProjectDTO, AtlasEdgeDTO, AtlasGraphDTO
- ForgeExecutionDTO
- GuardianFindingDTO, GuardianSummaryDTO, GuardianReportDTO
- NexusMetricsDTO

**Herald — client/client.go**
- Added NewForService(baseURL, token) constructor for non-Nexus upstreams
- Added NexusMetrics(), Atlas(), Forge(), Guardian() accessor methods
- WaitReady error message generalised ("service" not "engxd")
- All existing methods and signatures unchanged

**Herald — client/atlas.go (new)**   Projects(), Graph()
**Herald — client/forge.go (new)**   History(), ByTrace()
**Herald — client/guardian.go (new)** Findings()
**Herald — client/metrics.go (new)** Get() for Nexus /metrics

## Apply

```bash
# Accord first
cd ~/workspace/projects/engx/Accord
unzip -o /mnt/c/Users/harsh/Downloads/engx-drop/ENGX-HERALD-P1-001.zip "accord/*" -d .
go build ./...
git add api/upstream.go
git commit -m "feat(accord): ADR-039 — upstream DTOs for Atlas, Forge, Guardian, NexusMetrics"
git push origin main

# Herald second
cd ~/workspace/projects/engx/Herald
unzip -o /mnt/c/Users/harsh/Downloads/engx-drop/ENGX-HERALD-P1-001.zip "herald/*" -d .
go build ./...
git add client/client.go client/atlas.go client/forge.go client/guardian.go client/metrics.go
git commit -m "feat(herald): ADR-039 — Atlas, Forge, Guardian, NexusMetrics typed clients"
git push origin main
```

## Verify

```bash
go build ./...   # must be zero errors in both repos
# No behaviour change — existing callers unaffected
# New accessors are additive only
```
