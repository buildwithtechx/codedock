# Codedock Security & Architectural Roadmap

## Marketing Website (`apps/web`) Enhancements (Dokploy Parity & Multi-Page Architecture)

### 1. Header & Navigation Mega-Menus (`apps/web/src/components/header.astro`)

- [ ] **Features Dropdown**: Add interactive mega-menu links for Application Deployment, Databases, Monitoring, and AI Deployment.
- [ ] **Solutions Dropdown**: Add links for Enterprise, Self-Hosted PaaS, and Agencies.
- [ ] **Resources Dropdown**: Add links for One-Click Templates, Comparison Suite (vs Coolify, vs Dokploy, vs Vercel, etc.), Blog/Changelog, and FAQ.

### 2. Dedicated Comparison Suite (`apps/web/src/pages/vs/`)

- [ ] **Codedock vs. Coolify (`src/pages/vs/coolify.astro`)**: Compare sub-30MB Go binary daemon vs Node.js runtime, native SQL data browser, and Canvas view.
- [ ] **Codedock vs. Dokploy (`src/pages/vs/dokploy.astro`)**: Highlight single-binary architecture, embedded Docker engine deployer, and zero-dependency footprint.
- [ ] **Codedock vs. Vercel (`src/pages/vs/vercel.astro`)**: Contrast zero egress fees, self-hosted data ownership, and unlimited container deploys.
- [ ] **Codedock vs. Render (`src/pages/vs/render.astro`)**: Show cost savings, custom hardware flexibility, and direct server access.
- [ ] **Codedock vs. Portainer (`src/pages/vs/portainer.astro`)**: Compare developer-first PaaS DX vs generic Docker management.
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

### 5. Multi-Column Footer & Homepage Polish

- [ ] **Multi-Column Footer (`apps/web/src/components/footer.astro`)**: Update to 5-column layout (Product, Enterprise, Solutions, Compare & Learn, Company) matching Dokploy footer.
- [ ] **Homepage Visual Elevation (`apps/web/src/pages/index.astro`)**: Upgrade hero with install script copy box, GitHub/Discord CTA badges, glassmorphic control plane preview, and glowing dark mesh background.
