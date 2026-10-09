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
- `codedock version` - View CLI version

Full reference with every flag: [cmd/codedock/README.md](../../cmd/codedock/README.md).

## How it works

On install, the `postinstall` script downloads `codedock_<os>_<arch>.tar.gz` from the GitHub release matching this package's version and extracts the binary next to the launcher. Set `CODEDOCK_SKIP_POSTINSTALL=1` to skip the download (for example when running `npm ci` inside this monorepo). Releasing a new CLI version means building the Go binaries, tagging the matching release, and publishing this package at the same version.
