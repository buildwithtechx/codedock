# Codedock

**Self-hosted platform for shipping apps with built-in CI/CD.** Connect a server, push your code, and Codedock handles the build, rollout, domains, and certificates. Manage it all from the dashboard, desktop app, or CLI.

[![Latest release](https://img.shields.io/github/v/release/buildwithtechx/codedock?color=0b7285)](https://github.com/buildwithtechx/codedock/releases)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)
[![Website](https://img.shields.io/badge/website-codedock.run-0b7285)](https://codedock.run)

[Quick Start](#quick-start) · [How It Works](#how-it-works) · [Interfaces](#interfaces) · [Docs](https://docs.codedock.run) · [Contributing](CONTRIBUTING.md)

---

## Quick Start

On a fresh Linux server (Ubuntu/Debian recommended), one command installs the daemon as a container:

```bash
curl -fsSL https://get.codedock.run | bash
```

Open the dashboard at `http://your-server-ip:8080`. The installer provisions the embedded Postgres database, the Traefik edge, and all generated secrets for you.

Prefer binaries? Download `codedockd` and `codedock` from the [releases page](https://github.com/buildwithtechx/codedock/releases) and run `codedockd serve`.

**Ship something:**

```bash
codedock login
codedock project list
codedock deploy <service-id>
```

## How It Works

Give Codedock a **Git repository**, a **local directory**, or a **compose file**. From there the pipeline runs itself:

1. **Inspect.** Codedock reads your manifests, lockfiles, and Dockerfiles to pick the stack, builder, and start command. No config files needed.
2. **Build.** Your code becomes a Docker image via Dockerfile, Railpack, Nixpacks, or Buildpacks. Each deploy is snapshotted, so a rollback restores the exact image that shipped.
3. **Launch.** Containers start under supervision with health checks, and traffic swaps over with zero downtime.
4. **Expose.** Traefik routes your domain to the new containers and provisions a Let's Encrypt certificate that renews itself.
5. **Automate.** Git webhooks redeploy on push and spin up a preview environment for every pull request.

Databases, backups, TLS, and monitoring live in the same control plane as your apps.

## Interfaces

- **Web dashboard** — deploy, monitor, and manage every resource from the UI the daemon serves.
- **Desktop app** — native Tauri shell around the dashboard in `apps/desktop`.
- **`codedockd`** — the server daemon ([README](./cmd/codedockd/README.md)).
- **`codedock`** — the remote CLI that runs on your machine ([README](./cmd/codedock/README.md)).
- **REST API** — versioned HTTP API with a maintained spec in `docs/api/openapi.json`.

## Features

| | |
| --- | --- |
| **CI/CD included** | Redeploy on push, per-PR previews, one-click rollbacks |
| **Flexible builds** | Dockerfile, Railpack, Nixpacks, Buildpacks, monorepo-aware |
| **Managed data** | Postgres, MySQL, MongoDB, Redis, workers, WebSockets, storage |
| **One-click apps** | 35 templates with secrets, ports, domains, and volumes resolved |
| **TLS everywhere** | Let's Encrypt issuance and renewal, wildcard domains |
| **Backups** | Scheduled database and volume backups to S3/R2/MinIO with restore |
| **Observability** | Streaming build logs, container metrics, HTTP traffic analytics |
| **Growth path** | Autoscaling policies plus multi-replica HA cells |
| **No lock-in** | Plain Docker containers you can run anywhere, with or without Codedock |
| **Compose-native** | Bring your own compose files and run them unchanged |

## Deploy Anywhere

- **Any VPS** — Hetzner, DigitalOcean, and the rest
- **Dedicated servers** — bare metal, colo, homelab
- **Connected servers** — spread workloads across machines over SSH
- **HA cells** — multi-replica control plane on Postgres (`compose.ha.yml`)

One workflow, whatever hardware sits underneath.

## Local development

Requirements: Go 1.26+, Node.js 22+, Docker.

```bash
git clone https://github.com/buildwithtechx/codedock.git
cd codedock
cp .env.example .env
cp apps/dashboard/.env.example apps/dashboard/.env
npm run dev
```

See [CONTRIBUTING](./CONTRIBUTING.md) for the full guide.

## Status

Pre-release and moving fast. Expect breaking changes to APIs and schema until 1.0.

**Up next:** Codedock Cloud, multi-node hardening, and guided setup flows for the app catalogue.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

If you find a vulnerability, tell us privately — not through a public issue, PR, or discussion:

- **Preferred:** open a private [security advisory](https://github.com/buildwithtechx/codedock/security/advisories/new)
- Scope and handling process: [SECURITY.md](SECURITY.md)

## License

Apache-2.0. See [LICENSE](./LICENSE).
