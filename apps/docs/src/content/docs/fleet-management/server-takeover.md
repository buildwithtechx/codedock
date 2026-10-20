---
title: Server Takeover & Docker Adoption
description: Scan existing Linux hosts and adopt running Docker containers into Codedock without downtime.
---

Codedock's **Server Takeover** feature allows you to seamlessly onboard existing Linux servers that are already running production Docker containers. You can adopt and manage running workloads through Codedock without restarting containers or incurring downtime.

## Why Use Server Takeover?

- **Migrating from Other Platforms**: Easily migrate workloads from manual `docker-compose` setups, Coolify, CapRover, Portainer, or legacy VPS instances.
- **Zero-Downtime Transition**: Existing containers continue serving web traffic while Codedock inspects their configuration and takes over future lifecycle management.
- **Automated Configuration Extraction**: Codedock automatically reads container environment variables, mounted volumes, open ports, and restart policies.

## How Adoption Works

```
1. Scan Host ───────► 2. Review Detected Apps ───────► 3. Execute Adoption
(Inspect Docker)       (Map Ports & Volumes)           (Import without Restart)
```

### 1. Scan the Target Host

1. In the Codedock dashboard, navigate to **Fleet Management → Servers**.
2. Select the server you wish to inspect and click **Scan for Existing Workloads** (or trigger via `POST /api/system/takeover/scan`).
3. Codedock connects to the host's Docker engine and inspects all running and stopped containers.
4. The scanner identifies:
   - Container names, image repositories, and tag digests
   - Bound host ports and exposed internal ports
   - Docker volume mappings and host bind mounts
   - Environment variables (excluding sensitive system variables)
   - Network attachments and container dependencies

### 2. Review Detected Services

Once the scan completes, Codedock displays an interactive review screen listing all discovered containers:
- **Filter Exclusions**: Deselect system containers (such as monitoring agents, cloud-init processes, or diagnostic tools) that you do not want Codedock to manage.
- **Assign Projects**: Choose which Codedock Project and Environment each container should be imported into.
- **Configure Ingress**: Review host port bindings and assign custom domains or Traefik routing rules.

### 3. Adopt Workloads

Click **Adopt Selected Services** (or trigger via `POST /api/system/takeover/adopt`):
- Codedock records the service specifications, environment variables, and storage volume bindings in the PostgreSQL metadata store.
- Existing containers are left actively running without interruption.
- Traefik routing rules are applied dynamically, allowing you to route SSL-secured custom domains to adopted services immediately.

## Post-Adoption Management

Once adopted, your applications receive the full benefits of the Codedock platform:
- **Continuous Deployment**: Connect a Git repository or webhook to trigger automatic builds and updates on future code commits.
- **Live Logs & Metrics**: Stream stdout/stderr container logs and inspect CPU/RAM usage in real time.
- **Automated Backups**: If you adopted database containers (such as PostgreSQL, MySQL, or Redis), configure scheduled backups to S3 or Cloudflare R2 with one click.
- **Terminal Access**: Open interactive web terminal sessions directly into running containers via the dashboard.
