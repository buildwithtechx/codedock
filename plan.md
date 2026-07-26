# Codedock vs Dokploy (Competitive Analysis & Roadmap)

After completing the Multi-Server architecture, Organizations, Private Registries, and PR Previews, Codedock has fundamentally surpassed many of the initial roadmap goals.

A thorough scan of the `dokploy` codebase (UI and Backend) has revealed key architectural differences, areas where Codedock is already vastly superior, and a few minor features we can plan next to achieve total dominance.

## 1. Architectural Superiority

| Concern | Dokploy | Codedock (Winner) |
| :--- | :--- | :--- |
| **Backend Architecture** | Next.js API Routes + Node.js + BullMQ for queues | Single Go Binary (`codedockd`) + Goroutines |
| **Resource Footprint** | Heavy (requires Node.js, Redis for queues, etc.) | Extremely light (~15MB stripped binary) |
| **Clustering / Multi-Node** | **Docker Swarm** (Requires open ports, complex networking) | **WebSocket Worker Daemon** (`codedockw`) — Outbound only, zero open ports, firewall-friendly |
| **Database Support** | Postgres, MySQL, MariaDB, Mongo, Redis, LibSQL | **Postgres, MySQL, MariaDB, Mongo, Redis, Clickhouse, Kafka, RabbitMQ, Nats, Dragonfly, KeyDB** |
| **Deployment Speed** | Slow (Node.js overhead) | Lightning fast (Native Go + WebSocket tunneling) |

## 2. Security & Multi-Tenancy

Dokploy is primarily designed for self-hosting on a single or trusted cluster.
**Codedock is built for SaaS (`app.codedock.dev`) from day one.**

- **Terminal Isolation:** Codedock strictly enforces project-level RBAC (`ProjectService.HasPermission`) even on WebSocket terminal upgrades. This ensures true multi-tenant security where a logged-in user can only attach to containers in projects they own or belong to via their Organization.
- **Worker Authentication:** Codedock uses unique `worker_tokens` per server, meaning a compromised worker node cannot affect the control plane or other users' nodes.

## 3. New Findings & Next Steps (What we can improve on)

While Codedock is architecturally stronger, Dokploy has some UI/UX polishing that we should match or exceed:

1. **Data Browser Expansion:**
   - Currently, Codedock supports data browsing for Postgres/MySQL variants and Redis. We need to implement viewers for Mongo, Clickhouse, and messaging queues (Kafka/RabbitMQ), which currently show "currently unsupported".
2. **Advanced Traefik / Network UI:**
   - Dokploy has a dedicated `traefik.tsx` settings page for advanced routing rules. Codedock's Traefik configuration is currently mostly automatic. Exposing a UI for custom Traefik middlewares (like basic auth, rate limiting) per service would be highly valuable.
3. **Server SSH Keys Management:**
   - Dokploy has a dedicated UI for managing SSH keys for accessing servers. While `codedockw` avoids SSH entirely for deployments, allowing users to inject their public SSH keys into the worker nodes via the Codedock UI could be a nice "Server Management" feature.
4. **Hosted Deployment (Ops):**
   - The final remaining step is launching `app.codedock.dev`. The `Dockerfile.cloud` and `docker-compose.prod.yml` are ready.

## Conclusion

Codedock's decision to use a **Go-based WebSocket Worker Daemon** instead of Node.js + Docker Swarm makes it infinitely more scalable, secure, and easier to monetize as a SaaS. The focus should now shift to user acquisition, ops (launching the cloud version), and expanding the supported UI modules (Data Browsers, custom Traefik middlewares).

## 4. Security & Static Audit TODOs

A recent static audit revealed several critical authorization and configuration issues that must be addressed before SaaS deployment.

### Critical Security Issues

- [x] **Organization authorization bypass:** Add membership/role middleware to `internal/http/routes.go:159` and ownership checks to `internal/handlers/organization.go:52`.
- [x] **Global backup control exposed to all users:** Add user/project authorization to backup routes (`internal/http/routes.go:241`) and handlers (`internal/handlers/backup.go:24`).
- [x] **Scheduled-task authorization bypass:** Add project/service authorization to scheduled task routes (`internal/http/routes.go:294`) and handlers (`internal/handlers/scheduled_tasks.go:22`).
- [x] **Domain management authorization bypass:** Add service ownership checks to domain handlers (`internal/handlers/domain.go:22`).
- [x] **Environment deletion authorization bypass:** Add project-role middleware and ownership verify for `DELETE /api/environments/:id` (`internal/handlers/environment.go:54`).
- [x] **Unauthenticated server metrics WebSocket:** Add authentication and server ownership checks to the metrics WebSocket (`internal/handlers/server_metrics_ws.go:27`).
- [x] **Secrets returned in API responses:** Redact database passwords, registry tokens, worker tokens, and AI API keys before returning models.
- [x] **WebSocket origin checks disabled:** Enforce strict origin checks for Terminal, Worker, and Metrics WebSockets to prevent CSWSH.
- [x] **Worker token passed in query string:** Move worker token authentication to HTTP headers (`internal/handlers/worker_ws.go:32`, `bootstrap/worker.sh:88`).
- [x] **Password-reset URLs trust attacker-controlled headers:** Use a configured base URL instead of `Origin`/`Referer` headers (`internal/services/auth_service.go`).

### High-risk Correctness/Security Issues

- [x] Environment and Database listings filter by project existence, not user membership (`internal/handlers/app_services.go:114`, `internal/handlers/database.go:42`).
- [x] DNS records are globally readable/mutable by any authenticated user (`internal/handlers/dns.go:21`).
- [x] Project creation must verify organization membership before creation (`internal/handlers/project.go:49`).
- [x] Password reset tokens should be single-use and revoke existing sessions.
- [x] Enforce consistent password strength on registration and reset.
- [x] Remove the dangerous fallback in `baseAuth` that grants admin identity if token service is nil (`internal/http/middleware/auth.go:112`).
- [x] Project-level authorization must strictly protect the database query endpoint (`internal/services/database_query.go:18`).
- [x] Log-drain SSRF validation is vulnerable to DNS rebinding (`internal/handlers/app_services.go:368`).

### Bootstrap/Update Problems

- [x] Installer executes remote Docker installer code directly and downloads scripts without checksums (`bootstrap/install.sh`).
- [x] Default compose image tag is `latest`; needs pinning for reproducibility.
- [x] `docker-compose.prod.yml` uses unpinned `watchtower:latest`.
- [x] Add path-traversal protection to `restore.sh` when extracting tar archives.
- [x] Worker installer embeds token directly in systemd unit.
- [x] Fix naming mismatch: systemd installer references `codedock` while compose container is `codedock-control-plane`.
- [x] Fix uninstall script binary path (`codedockctl` vs `codedockd`).

### Dashboard Issues

- [x] Move access/refresh tokens from localStorage to HttpOnly cookies to prevent XSS theft (`apps/dashboard/src/stores/authStore.ts`).
- [x] Remove JWT from URL query string in live-log WebSocket (`live-logs-viewer.tsx`).
- [x] Terminal WebSocket has no explicit auth header and relies purely on cookies (needs CSRF/CSWSH mitigation).
- [x] Dashboard route layout has no client-side auth guard (`apps/dashboard/src/routes/_dashboard.tsx`).
- [x] Backup UI exposes restore functionality which lacks backend authorization.

### TODOs/Placeholders/Stale Values

- [x] Implement CLI `status` functionality.
- [x] Update MCP version from hardcoded `1.0.0` (`internal/http/bridge.go:90`).
- [x] Update Bootstrap server hardcoded version `0.1.0`.
- [x] Implement "future" worker log streaming behavior (`internal/engine/worker_hub.go:120`).
- [x] Refactor large files exceeding 350 lines (e.g., `git_service.go`, `app_services.go`, `github-integration.tsx`).

## 5. Structural Cleanup & Domain Boundaries

The overall structure of Codedock is reasonable, but both the Go backend and the React dashboard need cleanup and stronger domain boundaries to prevent monolithic God-classes and ensure maintainability.

### Go Backend Structure

The current structure follows the intended layered architecture (`cmd/`, `internal/`, `pkg/`), but main structural problems exist:

- [ ] **Broad Handlers:** `internal/handlers/` is becoming too broad. Split by domain (e.g., `auth/`, `users/`, `projects/`, `deployments/`, `backups/`).
- [ ] **Bloated Services:** `internal/services/` has too many unrelated responsibilities. Group related services similarly to handlers.
- [ ] **Mixed Engine Package:** `internal/engine/` is a large mixed package. Separate infrastructure concerns (e.g., `docker/`, `cron/`, `networking/`, `observability/`).
- [ ] **CLI Mix:** `cmd/codedockd/` mixes CLI parsing, setup, deployment, DB ops, and password management. Move command implementations into `cmd/codedockd/commands/`.
- [ ] **Script Overlap:** `bootstrap/` and `scripts/` overlap in responsibilities (upgrade, restore, install, backup). Define one canonical lifecycle location.
- [ ] **Public SDK Leakage:** `pkg/http/` imports internal models, making it less reusable as a public SDK.

**Note:** Large files violating the 350-line rule (`git_service.go`, `app_services.go`, `backup_manager.go`, `backup.go`, `service.go`) should be refactored as part of this structural split.

### Dashboard Structure

The React dashboard structure is generally good, but suffers from inconsistent domain ownership. The biggest structural improvement is to make backend authorization domain-specific and colocate dashboard API/hooks/types with each feature.

- [ ] **Colocate Domain Concepts:** `services/`, `hooks/`, and `interfaces/` duplicate domain concepts. Colocate them within `features/<domain>/` (e.g., `features/projects/api.ts`, `hooks.ts`, `types.ts`).
- [ ] **Feature Granularity:** `features/instance/` is too broad. Split into `features/settings/`, `backups/`, `dns/`, `notifications/`, `users/`, `registries/`.
- [ ] **Fix Kebab-case Naming:** Rename files violating conventions (e.g., `authStore.ts` -> `auth-store.ts`, `useAuth.ts` -> `use-auth.ts`, `apiClient.ts` -> `api-client.ts`).
- [ ] **API Client Refactor:** Separate token/auth refresh behavior from `api-client.ts` so it isn't overly central.
- [x] **Route Protection:** Centralize route protection in the dashboard layout or router context, rather than relying on API failures to trigger redirects.
- [x] **Build Artifacts:** Stop tracking `tsconfig.tsbuildinfo` in the dashboard directory.
