# codedock

Command line interface for Codedock self-hosted PaaS.

## Installation

```sh
npm install -g codedock
```

Or run directly without installation:

```sh
npx codedock --help
```

A Go build installs the same `codedock` command name (`curl -fsSL https://get.codedock.run/cli | sh`) - pick one, since the two builds differ in command shape.

## Quick Start

### 1. Authenticate

```sh
codedock login --server http://your-codedock-host:8080
```

### 2. View Status

```sh
codedock whoami
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
- `codedock whoami` - Display authenticated user
- `codedock status <service-id>` - Check app service status
- `codedock servers [list|create|delete]` - Manage servers
- `codedock projects [list|create|delete]` - Manage projects
- `codedock apps [list|create|delete|logs|deployments]` - Manage apps
- `codedock env [list|get|set]` - Manage project environment variables
- `codedock db [list|create|backup]` - Manage databases and backups
- `codedock deploy [path|service-id]` - Trigger deployment
- `codedock version` - View CLI version
