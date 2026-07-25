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

## 3. Features Completed ✅

The old `multiserver.md` roadmap is now officially completed!

- [x] **Agent-based Worker Daemon:** `codedockw` implemented and installer script created.
- [x] **Multi-Server Database Schema:** `servers` table linking to `projects`.
- [x] **Organizations & Teams:** Full RBAC and team invites.
- [x] **Private Docker Registries:** UI and Backend fully wired.
- [x] **PR Previews:** GitHub/GitLab webhooks creating ephemeral environments.
- [x] **Volume Backups:** Integrated snapshotting to S3.
- [x] **Terminal & Metrics:** Real-time streaming via WebSockets.
- [x] **Docker Image Deploy Flow:** UI modal implemented.

## 4. New Findings & Next Steps (What we can improve on)

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
