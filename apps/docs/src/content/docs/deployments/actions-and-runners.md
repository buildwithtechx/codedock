---
title: Actions & CI/CD Runners
description: Native containerized workflows and self-hosted CI runner management.
---

Codedock includes a built-in **Actions Engine** and **Runner Catalog** that allows you to execute containerized CI/CD pipelines, automated build scripts, and self-hosted GitHub Actions or GitLab CI runners directly on your server fleet.

## Actions Architecture

```
[ Git Push / Webhook ] ───► [ Codedock Actions Engine ] ───► [ Isolated Runner Container ]
                                       │                               │
                                       ▼                               ▼
                             [ Live Log Streaming ]         [ Build Artifact / Deployment ]
```

## Workflows & Pipelines

You can define custom automated workflows that run on code pushes, schedules, or manual triggers:

- **Pipeline Definition**: Write declarative workflow steps executed inside ephemeral Docker containers.
- **Workflow Triggers**: Trigger workflows from Git webhooks, project deployment events, or manual clicks.
- **Live Execution Streaming**: Watch build and test output live in the Codedock dashboard via WebSocket log streams.
- **Run Controls**: Cancel in-flight jobs (`POST /api/actions/runs/:id/cancel`) or retry failed workflow runs (`POST /api/actions/runs/:id/retry`) with one click.

## Runner Catalog & Self-Hosted CI Runners

Instead of paying for expensive cloud CI runner minutes, Codedock allows you to run self-hosted runners on your own server fleet:

### Supported Runner Types

- **GitHub Actions Runners**: Self-hosted `actions-runner` containers paired directly with your GitHub organizations or repositories.
- **GitLab CI Runners**: Native GitLab runner instances executing jobs for your self-hosted or cloud GitLab projects.
- **Custom Docker Runners**: Generic automation runners for executing bash scripts, security scans, or automated database seeders.

### Registering a Runner

1. In the Codedock dashboard, navigate to **Actions → Runners**.
2. Click **Register Runner** (or use `POST /api/actions/runners`).
3. Select your runner platform (GitHub, GitLab, or Generic).
4. Enter your registration token and target host server.
5. Codedock automatically pulls the runner container image, registers with your Git provider, and starts polling for CI jobs.

## Host Inspection & Resource Isolation

Before dispatching intensive builds, Codedock inspects runner destinations (`POST /api/actions/runners/inspect`):

- **Resource Verification**: Verifies available CPU cores, free RAM, and storage on the host.
- **Ephemeral Cleanups**: Automatically prunes stale build caches, dangling build images, and leftover work directories after job completion.
- **Concurrency Caps**: Sets limits on simultaneous job execution per runner to prevent build jobs from starving production web applications of compute power.
