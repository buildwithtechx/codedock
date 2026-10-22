# Codedock Platform & Architecture Roadmap: Upstream Sync & Target Prioritization

> **Upstream OpenShip Sync:** `7f0b4f882` → `6887ba85a` → `ed662afe7`  
> **Sync Date:** October 2026  
> **Source Scope:** 680+ files changed, +29,500 / -3,600 lines across API, Dashboard, Adapters, Platform, and Core engines.

---

## Executive Summary

The latest OpenShip upstream release cycle (`6887ba85a..ed662afe7`) brings critical operational resilience, advanced multi-tenant infrastructure controls, and essential UX polish:

1. **Already Completed in Codedock (Phases 1–3 & Parity Sync)**:
   - **Build Artifact Retention & Instant Retry** (`deployment_retry.go`, `deployment-detail.tsx`, `deployment-history-card.tsx`).
   - **Compose Required Env Preflight Check** (`env_preflight.go`).
   - **Cloud Metadata SSRF Egress Firewall** (`build-metadata-policy.ts`, `metadata_firewall.go`).
   - **Continuous Terminal Authority Watcher & Boundary Hardening** (`terminal_authority.go`).
   - **Safe Compose Overlap Cutover & Port Settlement** (`cutover.ts`, `cutover_decision.go`).
   - **Docker Image Volume (`VOLUME`) Preservation** (`image_volumes.go`).
   - **Empty Storage Reset with Volume Retention** (`empty-volume-plan.ts`, `empty-volume-dialog.tsx`, `service_volumes.go`).
   - **Runtime Base Image Security Update Detection & Menu** (`updates.ts`, `runtime-image-update-notice.tsx`, `runtime-redeploy-menu.tsx`, `runtime_updates.go`).
   - **Scoped Private Git Submodule Authentication** (`git_submodule_auth.go`).
   - **ACME SSL 502/504 Diagnostic Precision** (`acme-summary.ts`, `acme_diagnostics.go`).
   - **Dashboard ↔ Backend Route & Schema Parity** (Project members, Redis data browser, Stack restart/stop, Job aliases).

2. **New Targets from Upstream Pull (`6887ba85a..ed662afe7`)**:
   - **Subdomain Wildcard TLS Certificate Inheritance & Renewal Conf Purge** (PR #1323 & #1326): Automatic inheritance of parent `*.domain.com` certificates, preventing duplicate ACME challenges and rate-limit burn.
   - **Per-Service Swap Disabling & Memory Safety** (PR #1415): Opt-in `swapDisabled` flag pinning `MemorySwap` to `Memory` for databases and cache engines.
   - **Background Restart Loop Detection & Incident Monitoring** (PR #1198 / #672): Non-blocking 60-second watch after deploy settling to flag crash-looping services before users notice.
   - **Fine-Grained Token Capabilities & Domain-Scoped Access Grants** (PR #1304): Granular project capabilities (`domains:read|write|delete`, `deployments:read|trigger`) for tokens and MCP tools.
   - **Live GitHub Stars Nav Counter** (PR #1412): Edge-cached GitHub repo stars badge in navigation.
   - **Catalog Automatic Credential Generation** (PR #1309): Pre-populates required passwords (`POSTGRES_PASSWORD`, etc.) to prevent blank password launch failures.
   - **Hosted Manage Plans vs Cloud Pricing Split** (PR #1397): Clear product model distinguishing self-hosted servers from managed hosted control planes and cloud capacity.

---

## Roadmap Priority Matrix

| Initiative | Priority | Domain | Complexity | Status in Codedock |
| :--- | :---: | :---: | :---: | :---: |
| [1. Build Artifact Retention & Instant Retry](#1-build-artifact-retention--instant-retry) | **P0** | Deployments | Medium | **COMPLETED** (`5a24c9f8`) |
| [2. Compose Required Env Preflight Check](#2-compose-required-env-preflight-check) | **P0** | Deployments | Low | **COMPLETED** (`5a24c9f8`) |
| [3. Cloud Metadata SSRF Egress Firewall](#3-cloud-metadata-ssrf-egress-firewall) | **P0** | Security / BuildKit | Low | **COMPLETED** (`c92c3b3c`) |
| [4. Terminal Authority & Tenant Boundary Hardening](#4-terminal-authority--tenant-boundary-hardening) | **P0** | Security / Terminal | Medium | **COMPLETED** (`c92c3b3c`) |
| [5. Safe Compose Zero-Downtime Cutover & Port Settlement](#5-safe-compose-zero-downtime-cutover--port-settlement) | **P1** | Runtime / Ingress | Medium | **COMPLETED** (`6adba2d4`) |
| [6. Docker Image Volume Preservation](#6-docker-image-volume-preservation) | **P1** | Storage | Low | **COMPLETED** (`5a24c9f8`) |
| [7. Empty Storage Reset with Volume Retention](#7-empty-storage-reset-with-volume-retention) | **P1** | Storage | Low | **COMPLETED** (`6adba2d4`) |
| [8. Runtime Base Image Update Detection & Redeploy Menu](#8-runtime-base-image-update-detection--redeploy-menu) | **P1** | Runtime / Updates | Medium | **COMPLETED** (`6adba2d4`) |
| [9. Scoped Private Git Submodule Auth](#9-scoped-private-git-submodule-auth) | **P1** | Git / GitHub App | Medium | **COMPLETED** (`c92c3b3c`) |
| [10. ACME SSL Challenge Error Classification](#10-acme-ssl-challenge-error-classification) | **P2** | Ingress / Certbot | Low | **COMPLETED** (`6adba2d4`) |
| [11. Subdomain Wildcard TLS Certificate Inheritance](#11-subdomain-wildcard-tls-certificate-inheritance) | **P1 (Immediate)** | Ingress / Nginx | Low | **COMPLETED** (Phase 4) |
| [12. Per-Service Swap Disabling & Memory Safety](#12-per-service-swap-disabling--memory-safety) | **P1 (Immediate)** | Runtime / Resources | Medium | **COMPLETED** (Phase 4) |
| [13. Background Restart Loop Watch & Incident Detection](#13-background-restart-loop-watch--incident-detection) | **P1 (High)** | Monitoring / Engine | Medium | **COMPLETED** (Phase 4) |
| [14. Live GitHub Stars Nav Counter](#14-live-github-stars-nav-counter) | **P2 (Polish)** | Marketing & Web | Low | **COMPLETED** (Phase 4) |
| [15. Catalog Automatic Credential Generation](#15-catalog-automatic-credential-generation) | **P2 (Polish)** | Catalog / UI | Low | **COMPLETED** (Phase 4) |
| [16. Fine-Grained Token Capabilities & Domain Scoping](#16-fine-grained-token-capabilities--domain-scoping) | **P2 (Strategic)** | Auth & Permissions | Medium | **COMPLETED** (Phase 5) |
| [17. Hosted Manage Plans vs Cloud Pricing Architecture](#17-hosted-manage-plans-vs-cloud-pricing-architecture) | **P2 (Strategic)** | Billing / Plans | Medium | **COMPLETED** (Phase 5) |
| [18. Docker Swarm Multi-Node Combined Orchestration](#18-docker-swarm-multi-node-combined-orchestration) | **P3 (Long-Term)** | Multi-Host Cluster | High | **RESEARCH** |

---

## Detailed New Initiatives & Specifications

### 11. Subdomain Wildcard TLS Certificate Inheritance

- **Upstream Reference**: PR #1323 (`fix/wildcard-cert-inheritance-subdomains`, commit `140dd1289`) & PR #1326 (`fix/cleanup-obsolete-certbot-renewal-conf`, commit `d56ba1498`).

- **Problem**: When users create subdomains (e.g. `api.example.com`, `admin.example.com`) after already provisioning a wildcard certificate (`*.example.com`), Nginx previously forced independent ACME challenges for every subdomain, risking Let's Encrypt rate limits and causing unnecessary issuance delays.
- **Architecture for Codedock**:
  - Add `parentWildcardDomain(domain string) string`: strips the leading label to derive `*.example.com`.
  - In Nginx certificate resolution: if a direct certificate for `sub.domain.com` does not exist on disk, check if `*.domain.com/fullchain.pem` exists. If present, point Nginx `ssl_certificate` and `ssl_certificate_key` to the parent wildcard.
  - When a custom manual certificate is uploaded replacing an ACME lineage, automatically remove obsolete Certbot renewal configuration files (`/etc/letsencrypt/renewal/<domain>.conf`) to avoid renewal conflicts.

---

### 12. Per-Service Swap Disabling & Memory Safety

- **Upstream Reference**: PR #1415 (`feat/GH-1168-disable-service-swap`, commit `95b8ef715` & `8e9530769`).

- **Problem**: Memory-intensive workloads like Redis, PostgreSQL, and Elasticsearch suffer major performance degradation when operating systems start swapping container memory to disk.
- **Architecture for Codedock**:
  - Add `SwapDisabled bool` to service resource configuration models (`models.AppService`, service inputs, compose parser).
  - In container creation (`container_manager.go` / Docker runtime): when `SwapDisabled` is true and a memory limit is set, configure container `HostConfig.MemorySwap` equal to `HostConfig.Memory` (which disables additional swap in Docker).
  - Expose a toggle **"Disable swap memory"** in the service resource limits panel and project topology inspector.

---

### 13. Background Restart Loop Watch & Incident Detection

- **Upstream Reference**: PR #1198 / PR #672 (`feat/672-background-restart-loop-watch`, commits `b2965c6da`, `f7ffc9bb0`, `f4224dc01`).

- **Problem**: Deployments that start containers successfully report "Live" or "Ready" immediately. If the container crashes 10–30 seconds later (e.g. bad database connection string or runtime syntax error), the dashboard still displays green until a manual refresh or user complaint.
- **Architecture for Codedock**:
  - After deployment settles as `READY` / `ACTIVE`, spawn a non-blocking 60-second background watchdog (`WatchPostDeployContainerHealth`).
  - If a container enters a crash loop (`RestartCount > 2` or exit status non-zero within the first minute), record `restartLoop` metadata with exit code and the last 20 log lines.
  - Surface a prominent **"Action Required: Restart Loop Detected"** warning banner on the project and service overview with quick links to logs and redeploy.

---

### 14. Live GitHub Stars Nav Counter

- **Upstream Reference**: PR #1412 (`feat/web-github-stars`, commit `1adbefb65`).

- **Problem**: Marketing site and web navigation lacked social proof showing live repository stars.
- **Architecture for Codedock**:
  - Backend endpoint `GET /api/github-stars` in Codedock server with a 15-minute in-memory cache and `Cache-Control: public, max-age=900`.
  - Fetches `stargazers_count` from GitHub API.
  - `apps/web/src/components/main/header.tsx` renders a live compact stars badge (e.g. `★ 1.4k stars`) beside the GitHub link.

---

### 15. Catalog Automatic Credential Generation

- **Upstream Reference**: PR #1309 (`fix/GH-1309-generate-required-catalog-passwords`, commit `7bdfc7092`).

- **Problem**: Launching one-click catalog apps (PostgreSQL, MySQL, Redis, MinIO) with blank passwords caused container boot failures or unauthenticated vulnerabilities.
- **Architecture for Codedock**:
  - When selecting a catalog template, auto-generate cryptographically secure alphanumeric passwords for required credentials (e.g. `POSTGRES_PASSWORD`, `MYSQL_ROOT_PASSWORD`, `REDIS_PASSWORD`).
  - Display generated values in the setup form with a single-click copy button, ensuring zero empty-password deployment rejections.

---

### 16. Fine-Grained Token Capabilities & Domain Scoping

- **Upstream Reference**: PR #1304 (`feat/domain-scoped-token-grants`, commit `aa8290575`).

- **Problem**: Personal access tokens only supported broad roles (`admin`, `write`, `read`). Automation and MCP integrations need least-privilege scoping.
- **Architecture for Codedock**:
  - Introduce granular capability permissions: `domains:read`, `domains:write`, `domains:delete`, `deployments:read`, `deployments:trigger`.
  - Support capability filters in token generation UI and route guards.

---

### 17. Hosted Manage Plans vs Cloud Pricing Architecture

- **Upstream Reference**: PR #1397 (`feat/manage-and-cloud-plans`, commit `f2aa42289` / `b4499b90d`).

- **Details**:
  - Separation of concerns:
    - **Self-Hosted Community**: 100% free open-source control plane.
    - **Manage (Hosted Control Plane)**: Connect user-owned servers (Hetzner, AWS, DigitalOcean) to a managed Codedock control plane ($5/mo for 2 servers + $3/mo per extra server).
    - **Cloud Spaces (Managed Compute)**: Fully managed container infrastructure with provisioned CPU/RAM/NVMe storage.

---

### 18. Docker Swarm Multi-Node Combined Orchestration

- **Upstream Reference**: `feat/swarm-combined-20261022` (commits `05125d334`, `7753ed94c`, `8bf6a87a0`).

- **Details**:
  - Integrates Docker Swarm clusters directly into Codedock runtime:
    - **Manager & Worker Node Discovery**: Automatic swarm init, join-tokens, and dynamic cluster topology observation.
    - **Overlay Ingress & Routing Mesh**: Built-in ingress mesh across all swarm nodes with service mesh discovery.
    - **Rolling Zero-Downtime Service Updates**: Updates deployed using Docker Swarm service spec (`Order: start-first` or `stop-first`, update parallelism and delay).
    - **Host Constraint Placement & Resource Limits**: Strict node labeling, memory limits, and rejection of unsupported per-node swap overrides in swarm mode.

---

## Next Immediate Action Plan

```text
Phase 4: Ingress & Runtime Resilience
  ├── [x] 1. Subdomain Wildcard TLS Certificate Inheritance & Conf Cleanup (PR #1323, #1326)
  ├── [x] 2. Per-Service Swap Disabling & Memory Limits (PR #1415)
  ├── [x] 3. Background Post-Deploy Restart Loop Watcher (PR #1198 / #672)
  ├── [x] 4. Live GitHub Stars Nav Counter in API & Web Header (PR #1412)
  ├── [x] 5. Catalog Automatic Credential Generation (PR #1309)
  ├── [x] 6. Service Connection Toast Copy Polish (PR #1319 / GH-1311)
  └── [x] 7. Actionable Multi-Service Operation Refusals (PR #1313 / GH-1308)

Phase 5: Permissions, Governance & Managed Plans
  ├── [x] 1. Fine-Grained Token Capabilities & Domain-Scoped Access Grants (PR #1304)
  │         ├── domains:read, domains:write, domains:delete
  │         └── deployments:read, deployments:trigger
  ├── [x] 2. Hosted Manage Plans vs Cloud Pricing Split Architecture (PR #1397)
  │         ├── Management billing contract & schemas (packages/contracts/src/management-billing.ts)
  │         └── Server slot capacity enforcement & provider checkout reconciliation
  ├── [x] 3. Connected Server Health Auditing & Tenant Boundary Checks (PR #1411)
  └── [x] 4. Bounded Background Service Incident Tracking & Attention Derivation (PR #1198)

Phase 6: Multi-Node Swarm Orchestration & Advanced Clustering
  ├── [ ] 1. Docker Swarm Overlay Network Integration (feat/swarm-combined-20261022)
  ├── [ ] 2. Multi-Host Node Discovery & Dynamic Scaling
  ├── [ ] 3. Swarm Deployment Runtime & Rolling Service Updates
  └── [ ] 4. Cross-Node Health Watch & Ingress Routing
```
