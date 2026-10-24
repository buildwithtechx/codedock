# codedock

Command line interface for Codedock self-hosted PaaS. This package installs the prebuilt Go binary for your platform — no Go toolchain needed.

## Installation

```sh
npm install -g codedock
```

Prefer a direct binary? Use the install script instead:

```sh
curl -fsSL https://get.codedock.run/cli | sh
```

## Quick Start

### 1. Authenticate

```sh
codedock login --server http://your-codedock-host:8080
```

### 2. View Status

```sh
codedock me
codedock servers list
codedock projects list
```

### 3. Deploy

```sh
codedock deploy .
```

## Commands

- `codedock login` - Authenticate with your server
- `codedock logout` - Clear credentials
- `codedock me` - Display authenticated user (`whoami` works too)
- `codedock status [service-id]` - Check app service status
- `codedock servers [list|create|delete]` - Manage servers
- `codedock projects [list|create|destroy]` - Manage projects
- `codedock environments [list|create|destroy]` - Manage environments
- `codedock apps [list|create|destroy]` - Manage apps, plus `secrets`, `domains`, `deployments`, `logs`
- `codedock db [list|create|destroy|import]` - Manage databases, plus `backups`
- `codedock compose [analyze|deploy]` - Deploy compose files
- `codedock deploy [path|service-id]` - Trigger deployment
- `codedock mcp [serve|status|tools]` - Stdio MCP bridge and tool inspection
- `codedock ai [ask|confirm]` - Query the instance-aware Codedock AI Copilot
- `codedock version` - View CLI version

Full reference with every flag: [cmd/codedock/README.md](../../cmd/codedock/README.md).

## Programmatic SDK

`codedock` includes the programmatic client SDK with zero runtime dependencies. It provides TypeScript types and client surfaces for actions, ingress, networks, clusters, auth, and billing.

### Installation

```sh
npm install codedock
```

### Usage

```ts
import { createCodedockClient } from 'codedock/client';

const client = createCodedockClient({
  baseUrl: 'https://panel.example.com/api',
  token: process.env.CODEDOCK_TOKEN, // personal access token (dashboard → profile → tokens)
  // tenantId: 'acme',               // optional tenant scoping
});

// Access typed API surfaces
const runs = await client.actions.listRuns();
const workflows = await client.actions.listWorkflows();
const ingressRoutes = await client.ingress.listRoutes();
const clusters = await client.clusters.listClusters();
```

### Client Surfaces

- `client.actions` — background workflows, runs, triggers, and execution status.
- `client.ingress` — domains, certificates, routing rules, and TLS settings.
- `client.networks` — virtual overlays, port bindings, and internal service discovery.
- `client.clusters` — remote node pools, cluster state, and heartbeat metrics.
- `client.auth` — session checks, profile management, and credentials.
- `client.billing` — subscription tier, usage quotas, and invoice history.

### Error Handling

The SDK exposes structured error classes:

```ts
import {
  CodedockApiError,
  CodedockBillingGateError,
  CodedockRateLimitError,
} from 'codedock';

try {
  await client.actions.triggerRun({ workflowId: 'wf-123' });
} catch (err) {
  if (err instanceof CodedockRateLimitError) {
    console.warn(`Rate limited. Retry after ${err.retryAfterSeconds}s`);
  } else if (err instanceof CodedockBillingGateError) {
    console.error('Subscription tier required:', err.message);
  } else if (err instanceof CodedockApiError) {
    console.error(`API Error (${err.status}):`, err.message);
  }
}
```

### Personal Access Tokens

To authenticate programmatic requests, generate a Personal Access Token in the Codedock dashboard under **Profile → Personal Access Tokens**. Tokens are hashed, scoped, and support granular role permissions.

## How it works

On install, the `postinstall` script downloads `codedock_<os>_<arch>.tar.gz` from the GitHub release matching this package's version and extracts the binary next to the launcher. Set `CODEDOCK_SKIP_POSTINSTALL=1` to skip the download (for example when running `npm ci` inside this monorepo). Releasing a new CLI version means building the Go binaries, tagging the matching release, and publishing this package at the same version.
