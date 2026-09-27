# Codedock Security & Architectural Roadmap

## Marketing Website (`apps/web`) Exhaustive Architecture (Dokploy Parity & UI Elevation)

### 1. Navigation Mega-Menus (`apps/web/src/components/header.astro`)

- [ ] **Features Dropdown**: Add links for Application Deployment, Databases, Monitoring & Logs, and AI Deployment.
- [ ] **Solutions Dropdown**: Add links for Enterprise, Self-Hosted PaaS, and Agencies.
- [ ] **Resources Dropdown**: Add links for One-Click Templates, Comparison Suite, Blog/Changelog, Philosophy, and Docs.

### 2. Full Comparison Suite (`apps/web/src/pages/vs/`)

- [ ] **Comparison Hub (`src/pages/vs/index.astro`)**: Hub comparing Codedock against all major alternatives.
- [ ] **Codedock vs. Coolify (`src/pages/vs/coolify.astro`)**: Compare sub-30MB Go binary daemon vs Node.js runtime, native SQL data browser, and Canvas view.
- [ ] **Codedock vs. Dokploy (`src/pages/vs/dokploy.astro`)**: Highlight single-binary architecture, embedded Docker engine deployer, and zero-dependency footprint.
- [ ] **Codedock vs. Vercel (`src/pages/vs/vercel.astro`)**: Contrast zero egress fees, self-hosted data ownership, and unlimited container deploys.
- [ ] **Codedock vs. Render (`src/pages/vs/render.astro`)**: Show cost savings, custom hardware flexibility, and direct server access.
- [ ] **Codedock vs. Portainer (`src/pages/vs/portainer.astro`)**: Compare developer-first PaaS DX vs generic Docker management.
- [ ] **Codedock vs. CapRover (`src/pages/vs/caprover.astro`)**: Highlight modern dashboard, built-in backups, and automated Traefik SSL.
- [ ] **Codedock vs. Dokku (`src/pages/vs/dokku.astro`)**: Highlight web UI, multi-server fleet control, and visual topology canvas.

### 3. Feature Deep-Dive Pages (`apps/web/src/pages/features/`)

- [ ] **Application Deployment (`src/pages/features/application-deployment.astro`)**: Git-native, Dockerfile, Compose, and Railpack build engine details.
- [ ] **Database Management (`src/pages/features/databases.astro`)**: Built-in SQL/NoSQL data browser, connection pooling, and automated S3 backups.
- [ ] **System Monitoring (`src/pages/features/monitoring.astro`)**: Live container metrics, WebSocket log tailing, and resource alerts.
- [ ] **AI Deployment (`src/pages/features/ai-deployment.astro`)**: Local LLM and vector database deployment (Ollama, Qdrant, PGVector).

### 4. Solutions & Templates (`apps/web/src/pages/`)

- [ ] **One-Click Templates Library (`src/pages/templates.astro`)**: Interactive catalog for instant Postgres, MySQL, Redis, MongoDB, Supabase, MinIO, N8N, and Grafana deployment.
- [ ] **Self-Hosted PaaS Solution (`src/pages/solutions/self-hosted.astro`)**: Deep dive into single-binary daemon installation and architecture.
- [ ] **Enterprise & Fleet Management (`src/pages/solutions/enterprise.astro`)**: Yamux tunnel agentless orchestration, RBAC, and telemetry metering.
- [ ] **Agencies & Teams (`src/pages/solutions/agencies.astro`)**: Multi-project isolation and client environment management.

### 5. Multi-Column Footer & Homepage Polish

- [ ] **5-Column Mega Footer (`apps/web/src/components/footer.astro`)**: Update to 5-column layout (Product, Enterprise, Solutions, Compare & Learn, Company) matching Dokploy footer.
- [ ] **Homepage Visual Elevation (`apps/web/src/pages/index.astro`)**: Upgrade hero with install script copy box, GitHub/Discord CTA badges, glassmorphic control plane preview, and glowing dark mesh background.

## Documentation (`apps/docs`) Roadmap (Dokploy Parity & Uncollapsed Navigation)

- [x] **Uncollapsed Sidebar Config (`apps/docs/astro.config.mjs`)**: Set `collapsed: false` on all Starlight sidebar groups so navigation categories stay permanently open.
- [ ] **Dokploy-Style Navigation Hierarchy (`apps/docs/astro.config.mjs`)**: Reorganize documentation into clean sections: Getting Started, Deployments, Databases, Storage & Backups, Networking & SSL, Fleet Management, Security & Operations, and Reference.
- [ ] **Missing Feature Docs (`apps/docs/src/content/docs/`)**: Write comprehensive docs for missing/sparse topics:
  - Docker Compose deployments & environment override rules
  - Railpack & Nixpacks build engine configuration
  - Native Database Data Browser & SQL Studio usage
  - Automated S3 / Cloudflare R2 / MinIO backup scheduling & restores
  - Canvas topology view & interactive dependency mapping
  - Traefik routing, custom domains, and automatic Let's Encrypt SSL
  - Yamux tunnel fleet management & agentless server setup
  - `codedockd` CLI commands and full REST API endpoint specifications
