# Codedock Website Roadmap

The marketing site uses React and TanStack Start in `apps/web`. Documentation uses Astro and Starlight in `apps/docs`. Site copy describes the current implementation; infrastructure features are not created by adding marketing pages.

## Marketing website

- [x] Features, Solutions and Resources desktop menus and accessible mobile navigation.
- [x] Feature pages for application deployment, databases, monitoring and AI workloads.
- [x] Solutions for self-hosting, teams and fleet management, and agencies.
- [x] Comparison hub and Coolify, Dokploy, Vercel, Render, Portainer, CapRover and Dokku pages, with official sources and balanced fit guidance.
- [x] Searchable, category-filtered recipe catalogue for PostgreSQL, MySQL, Redis, MongoDB, Supabase, MinIO, n8n, Grafana, Ollama and Qdrant.
- [x] Five-column footer with Product, Enterprise, Solutions, Compare & Learn and Company links.
- [x] Homepage install-command copying, GitHub and community links, dark hero and illustrative control-plane preview.
- [x] Correct Apache-2.0 licensing and remove unsupported statistics, testimonials, prices, release versions, binary-size and Yamux claims.
- [x] Page-specific metadata and canonical URLs, generated route tree and production-only development tools exclusion.
- [x] Website interaction checks and CI coverage for web types, tests, web builds and documentation builds.

Routes live in `apps/web/src/routes`; shared pages and content live in `apps/web/src/features`. `routeTree.gen.ts` is generated rather than edited manually.

## Documentation

- [x] All sidebar groups uncollapsed.
- [x] Getting Started, Deployments, Databases, Storage & Backups, Networking & SSL, Fleet Management, Security & Operations and Reference groups.
- [x] Compose conversion workflow, ignored environment/volume/port/dependency fields and retry behavior.
- [x] Dockerfile, Nixpacks, Railpack-selection and Buildpacks execution paths and fallback behavior.
- [x] Relational data browser, SQL Studio, engine support and local connection requirements.
- [x] Backup configuration, cron scheduling, S3-compatible R2/MinIO destinations, download and restore workflows.
- [x] Canvas navigation and the limits of relationship discovery.
- [x] Traefik routing, DNS, TLS email and HTTP-01 certificate requirements.
- [x] Direct SSH worker setup and unsupported tunnel transports.
- [x] Remote CLI, daemon CLI and installer wrapper command distinctions.
- [x] Current API authentication, core request examples and a catalogue of 210 registered API routes.
- [ ] Exhaustive REST request/response schemas or a maintained OpenAPI specification for every endpoint. The catalogue is a route inventory, not that contract.

## Separate implementation follow-ups

These are explicitly outside the marketing and documentation changes:

- A complete Compose runtime importer preserving overrides, variables, ports, volumes and dependencies.
- An actual Railpack CLI build path distinct from Nixpacks.
- One-click recipe provisioning for the complete catalogue, with validated secrets, storage and networking.
- Additional database browsing, remote query connectivity and compatible backup/restore metadata for all advertised engines.
- Profile personal-access-token authentication support in the request guard.
- Yamux/reverse-tunnel fleet transport and verified usage-based metering, if retained as product goals.

The original roadmap's zero-dependency, sub-30MB, blanket cost-saving and instant-template claims are not acceptance criteria without implementation and evidence.
