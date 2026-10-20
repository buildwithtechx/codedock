---
title: Rollbacks & Attention Health
description: Rollback capacity retention, instant recovery, and resource attention issues.
---

Codedock is designed for resilient operations, combining zero-downtime instant rollbacks with automated **Attention Health** diagnostics that actively monitor container crash loops, resource exhaustion, and Out-Of-Memory (OOM) events.

## Zero-Downtime Rollbacks

When a new code deployment introduces an unexpected bug or regression, you can revert to a known-good release version instantly:

1. Open your Service and navigate to the **Deployments** tab.
2. Locate the previous successful deployment.
3. Click **Rollback to this Release** (or trigger via `POST /api/deployments/:id/rollback`).
4. Codedock switches Traefik ingress traffic to the previous image snapshot and starts the healthy container before stopping the faulty version, ensuring zero downtime.

## Rollback Retention Policies

To ensure rollback images remain available without filling up your server's disk, Codedock provides granular **Rollback Retention Settings**:

- **Rollback Capacity**: Inspect current retained image count and storage consumption via `GET /api/projects/:id/rollback-capacity`.
- **Retention Period (`rollback_retention_days`)**: Set how many days rollback snapshots are retained before automatic pruning (e.g., 7, 14, or 30 days).
- **Max Retained Images**: Configure a hard ceiling on historical image versions kept per service.
- **Configure Policies**: Update settings in **Project Settings → Rollback Retention** or via `PATCH /api/projects/:id/rollback-settings`.

### Pinned Deployments
To protect critical stable releases (such as major milestones or baseline versions) from being deleted by automatic retention sweeps, click **Pin Deployment** (`POST /api/deployments/:id/pin`). Pinned deployments are permanently exempted from automated pruning.

## Attention Issues & Health Diagnostics

Codedock includes an active health detection engine (`attention_issues`) that analyzes container runtime events and flags potential operational issues before they lead to outages.

```
┌────────────────────────────────────────────────────────┐
│                   Attention Dashboard                  │
├───────────────────┬──────────────┬─────────────────────┤
│ Issue Type        │ Service      │ Diagnostic Details  │
├───────────────────┼──────────────┼─────────────────────┤
│ OOMKilled         │ api-server   │ Memory limit 512MB  │
│ CrashLoopBackOff  │ payment-svc  │ Exit code 137       │
│ DiskPressure      │ worker-01    │ Storage above 85%   │
└───────────────────┴──────────────┴─────────────────────┘
```

### Monitored Conditions

1. **Out-Of-Memory (OOM) Crashes**: Detects when Linux kills a container due to memory starvation (`State.OOMKilled`). Codedock alerts you with the exact memory ceiling and recommends increased memory allocations.
2. **Crash Loops**: Identifies containers that exit repeatedly with non-zero exit codes within a short timeframe, linking directly to container crash logs.
3. **Resource Saturation**: Warns when CPU throttling or disk I/O bottlenecks cause latency spikes in web requests.
4. **Ingress Misconfigurations**: Highlights services whose exposed ports do not match container listening ports.
