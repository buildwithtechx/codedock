# Multi-Server & Codedock Cloud Roadmap

This document outlines the architectural changes, features, and steps required to transform Codedock from a single-node manager into a **Distributed Fleet Manager** (Multi-Server) and **Codedock Cloud SaaS**.

## 1. Architectural Strategy: ✅ Agent-based Worker Daemon (DECIDED)

**Decision: Worker Daemon** — chosen because Codedock Cloud is the long-term goal.

### Why Worker Daemon wins for SaaS

| Concern                    | Agentless SSH                                                   | Worker Daemon ✅                                                       |
| -------------------------- | --------------------------------------------------------------- | ---------------------------------------------------------------------- |
| **Security**               | Control plane dials INTO user's server (requires open SSH port) | Worker dials OUT to `api.codedock.dev` (outbound only, firewall-friendly) |
| **Deployment reliability** | Deployment dies if connection drops                             | Deployment continues locally, reports back on reconnect                |
| **Real-time metrics**      | SSH polling — expensive at scale                                | Persistent WebSocket push — cheap at scale                             |
| **Monetisation gate**      | No natural gate                                                 | License key validated on worker registration                           |
| **Enterprise acceptance**  | Blocked by most firewalls                                       | Accepted (same model as GitHub Actions runners)                        |

### Architecture

```text
User's VPS                          Codedock Cloud (api.codedock.dev)
┌─────────────────────┐             ┌──────────────────────────┐
│  codedockw            │──WebSocket──▶  Control Plane (Go API)  │
│  - Runs deployments │◀─Commands───│  - Dashboard UI          │
│  - Streams logs     │──Metrics───▶│  - Billing (Stripe)      │
│  - Reports health   │             │  - License validation    │
└─────────────────────┘             └──────────────────────────┘
```

### How it works

1. User creates an account on `app.codedock.dev` and gets a **license/registration key**.
2. User runs a one-liner on their VPS: `curl -sL get.codedock.dev | bash -s -- --key <LICENSE_KEY>`
3. The script installs `codedockw` as a systemd service.
4. `codedockw` dials out to `api.codedock.dev` via WebSocket and registers itself.
5. The user's server appears in their dashboard. All deployments, logs, and metrics flow through the persistent WebSocket tunnel.

### Worker Binary (`cmd/codedockw/`)

- Written in Go, single static binary (~15 MB stripped).
- Connects to control plane via `gorilla/websocket`.
- Executes deployment commands received from the control plane using the local Docker socket.
- Streams container logs and metrics back via the same WebSocket.
- Reconnects automatically with exponential backoff if the connection drops.

---

## 2. Database Schema Changes

We need to make Codedock aware of physical servers. The **project** is the server boundary — all apps and databases inside a project inherit the server automatically. Users never pick a server when creating individual apps or databases.

- **`servers` table (New):**
  - `id` (UUID)
  - `user_id` — owner (for Codedock Cloud multi-tenancy)
  - `name` (e.g., "EU Production")
  - `ip_address` (e.g., "198.51.100.1")
  - `status` (`online`, `offline`, `provisioning`)
  - `worker_token` — the secret the `codedockw` binary uses to authenticate
  - `last_seen_at` — heartbeat timestamp
  - `metrics` (JSON — latest CPU/RAM/Disk snapshot pushed by the worker)

- **`projects` table (Update):**
  - Add `server_id` FK → `servers.id`
  - `NULL` means local (for single-node self-hosted installs — backward compatible)

- **`app_services` & `databases` tables — NO CHANGE ✅**
  - They resolve their server by joining through their parent project.
  - No `server_id` column needed here. The user never selects a server per-app or per-database.

---

## 3. Core Engine Updates (`internal/engine/`)

The Deployer engine currently assumes `localhost`. It needs to become **Project-Server-Aware**.

- **Command Routing via Worker WebSocket:**
  Instead of executing Docker commands directly, the deployer checks the project's `server_id`. If a server is attached, it serialises the deployment command as a JSON message and sends it down the worker's persistent WebSocket connection. The worker executes it locally and streams results back.

  ```text
  Deployer.Deploy(app)
    → look up app.Project.ServerID
    → if nil: run locally (existing behaviour, no change)
    → if set:  serialize command → send to WorkerHub → worker executes → stream back
  ```

- **WorkerHub (`internal/engine/worker_hub.go` — New):**
  Maintains a registry of live `server_id → WebSocket connection` pairs. When a worker connects, it registers here. When the deployer needs a remote server, it looks up the live connection from the hub.

- **No SFTP / SSH needed:** The worker has direct access to its own Docker socket. Build context (Dockerfile, source) is sent as a binary payload over the WebSocket, not via SFTP.

---

## 4. Networking & Traefik Routing

Currently, the single Traefik container routes all traffic. In a multi-server setup:

- **Distributed Proxies:** Every Worker Node must run its own instance of Traefik.
- **DNS Resolution:** When an app is deployed to `Server B`, the dashboard must show the user the IP address of `Server B` so they can point their Custom Domain's A-Record to the correct worker node.
- **Wildcard Domains:** If the user has a wildcard domain (e.g., `*.apps.mycodedock.com`), the DNS A-Record for the wildcard must point to the specific server hosting those apps, or we need a central load balancer.

---

## 5. Frontend & UI Integrations

- **Servers Dashboard (`/dashboard/servers` — New):**
  - List all connected worker servers with live status (Online/Offline), last seen, CPU/RAM.
  - "Add Server" page — shows the one-liner install command pre-filled with a fresh `worker_token`.
  - Per-server detail page: resource graphs, connected projects, events log.

- **Project Creation (Update — minor):**
  - Add an optional **"Deploy to Server"** dropdown. Defaults to "Local" for self-hosted installs.
  - Once set on the project, all apps and databases inside automatically deploy to that server. No per-app or per-database server selection ever shown. ✅

- **App/Database Creation — NO CHANGE ✅**
  - Server is fully inherited from the project. Users never see a server dropdown here.

- **Metrics UI:**
  - Per-server resource graphs on the Servers page.
  - Per-project resource graphs continue to work as before (metrics aggregated from the project's server).

---

## 6. Codedock Cloud (SaaS) Considerations

Once Multi-Server is built, launching Codedock Cloud is trivial:

1. Host the **Codedock Control Plane** centrally at `app.codedock.dev`.
2. Users create accounts (user auth already exists in the codebase).
3. Users spin up any VPS (DigitalOcean, Hetzner, AWS EC2, bare-metal), run the one-liner install command from their Servers dashboard, and the server appears live in seconds.
4. Users create a Project, select the server to deploy it to, and add their apps/databases as normal.
5. **Billing:** Stripe integration gates plan limits (e.g., Free = 1 server, Pro = 5 servers, Enterprise = unlimited).

---

## TODO Checklist

### Backend — Models & Repositories

- [x] Create `servers` model (`internal/models/server.go`)
- [x] Create `ServerRepository` (`internal/repositories/server.go`) — CRUD + `GetByWorkerToken`
- [x] Add `server_id` (nullable FK) to `projects` table only

### Backend — Services & Handlers

- [x] Create `ServerService` (`internal/services/server.go`) — enforces plan-based server limits (hobby = 1)
- [x] Create `ServerHandler` (`internal/handlers/server.go`) — REST endpoints
- [x] Add server routes to `internal/http/routes.go`
- [x] Plan-based project limit enforced in `ProjectHandler` (hobby = 2 projects)

### Backend — Worker Engine

- [x] Create `WorkerHub` (`internal/engine/worker_hub.go`)
- [x] Create Worker WebSocket endpoint (`/ws/worker`) — authenticate with `worker_token`, register in hub
- [x] Update `Deployer` — route to WorkerHub if `project.ServerID != nil`
- [x] WorkerHub handles heartbeats and updates `servers.last_seen_at` + `servers.status`

### Worker Binary

- [x] Scaffold `cmd/codedockw/` — separate Go entrypoint
- [x] Worker connects via WebSocket using `worker_token`
- [x] Worker executes deployment commands on local Docker socket
- [x] Worker streams logs and metrics back over the WebSocket
- [x] Worker reconnects with exponential backoff
- [x] Worker installs Traefik on first boot

### Frontend

- [x] Servers dashboard page (`/servers`) — list, status, last seen, metrics
- [x] Add Server flow — one-liner install command with pre-filled worker token
- [x] Project creation — optional "Deploy to Server" dropdown
- [x] App and Database creation — NO changes (server inherited from project) ✅
- [x] Git providers wired in dashboard (GitHub, GitLab, Bitbucket, Gitea)
- [x] `useGit.ts` hooks (`useGitStatus`, `useListGitRepos`, `useListGitBranches`, `useConnectGit`, `useDisconnectGit`)
- [x] `CreateGitAppModal` — deploy from connected provider or public URL
- [ ] Per-server resource graphs (CPU/RAM/Disk over time) — UI skeleton exists, needs real data wiring
- [ ] Docker Image deployment flow — currently shows `alert('coming soon!')`
- [ ] Private Docker Registry UI — backend done (`Registry` model + handler), dashboard pages not built

### Codedock Cloud SaaS

- [x] User registration + email verification
- [x] Stripe subscription billing
- [x] Plan-based limits validated server-side (server count, project count)
- [ ] Hosted deployment of control plane (`app.codedock.dev`)

---

## 8. Competitive Analysis (Codedock vs Dokploy)

### Features Dokploy has that Codedock doesn't (yet)

| Feature | Status | Notes |
|---|---|---|
| **PR Previews** | ❌ Not started | Webhook handler + ephemeral env lifecycle + preview domain routing |
| **More Git Providers** | ✅ Done | GitHub, GitLab, Bitbucket, Gitea — backend + dashboard fully wired |
| **Private Docker Registries** | ⏳ Backend only | Model + repo + service + handler exist; dashboard UI not built |
| **Docker Image Deploy Flow** | ⏳ Partial | `AppService.imageRef` field exists in backend; frontend shows `alert('coming soon!')` |
| **Organizations & Teams** | ❌ Not started | Full RBAC: `organizations` table, `organization_users` with roles |
| **Volume Backups** | ❌ Not started | Extend BackupService to support Docker volume snapshots to S3 |

### What's left (hardest → easiest)

1. **Organizations & Teams (RBAC)** — schema migration + service layer + full UI
2. **PR Previews** — GitHub webhook, ephemeral env lifecycle, preview routing
3. **Volume Backups** — BackupService extension + S3 upload path for volumes
4. **Private Registry UI** — backend already done, just needs list/create/delete pages per project
5. **Docker Image Deploy Flow** — modal similar to `CreateGitAppModal` but for `imageRef`
6. **Per-server resource graphs** — hook up existing metrics WebSocket to the graph UI
7. **Hosted deployment** (`app.codedock.dev`) — ops/infra

### Features we both have, but we did MUCH better (Where Dokploy went wrong)

- **Multi-Server Clustering:** Docker Swarm (them) vs WebSocket Worker Daemon (us) — outbound-only, no open ports.
- **Background Tasks:** Separate Node.js microservices (them) vs Go goroutines inside one binary (us).
- **Database Models:** Per-engine tables (them) vs one unified `database.go` with `engine_type` (us).

### Type Safety between Backend and Workers

- [x] Moved to shared Go structs — control plane and worker binary use the same native types, no raw JSON drift.
