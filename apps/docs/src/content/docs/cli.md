---
title: CLI Reference
description: Separate the remote client, daemon commands and installer host wrapper.
---

The repository has two remote clients named `codedock`: a Go binary in `cmd/codedock` and the npm client in `packages/cli`. Their commands and configuration differ. `codedockd` is the server daemon. The Linux installer also supplies a host management wrapper. Their command syntax and configuration are different.

## npm client: packages/cli

This section describes the npm client, not the Go binary. Use the client built from `packages/cli` from your workstation or CI runner. It connects to the daemon over HTTP. From the repository, build it with `npm run build:cli` and run `node packages/cli/dist/bin.js --help`.

```sh
node packages/cli/dist/bin.js login --server https://pilot.example.com --email owner@example.com
node packages/cli/dist/bin.js whoami
node packages/cli/dist/bin.js projects list
node packages/cli/dist/bin.js env set PROJECT_ID NODE_ENV=production
node packages/cli/dist/bin.js status SERVICE_ID
```

Login prompts for the password. Do not pass `--password`; that option is rejected. A token can be supplied with `--token`. Global options include `--server`, `--token`, `--json`, `--help` and `--version`. `CODEDOCK_SERVER_URL` and `CODEDOCK_TOKEN` can supply connection settings.

| Command | Subcommands / argument |
| --- | --- |
| `login` | Authenticate interactively or with a token |
| `logout` | Remove saved authentication |
| `whoami`, `me` | Current user |
| `status` | Service ID |
| `servers` | `list`, `create`, `delete` |
| `projects` | `list`, `create`, `delete` |
| `apps` | `list`, `create`, `delete`, `logs`, `deployments` |
| `env` | `list`, `get`, `set` |
| `db` | `list`, `create`, `backup` |
| `deploy` | Local path or existing service ID |
| `version` | Client version |

Read `codedock --help` for options. Creation commands use flags specific to the resource. For example, remote database creation accepts `--project`, `--type`, `--name` and `--server-id`.

## Go client: cmd/codedock

The compiled Go client is also named `codedock`. Its login prompts for the server URL, email and password and saves its configuration. It does not implement the npm client's global `--server`, `--token` or `--json` flags or environment-variable overrides.

```sh
codedock login
codedock me
codedock projects create --name "My API"
codedock environments list --project PROJECT_ID
codedock apps secrets set PORT=3000 --project PROJECT_ID
codedock db create --project PROJECT_ID --environment ENVIRONMENT_ID --name primary --engine postgres
```

| Command | Registered subcommands / argument |
| --- | --- |
| `login`, `logout`, `me` | Account and saved client configuration |
| `status` | Service ID |
| `projects` | `list`, `create`, `destroy` |
| `environments`, `env` | `list`, `create`, `destroy`; the `env` alias manages environments, not variables |
| `apps` | `list`, `create`, `destroy`, plus `secrets`, `domains`, `deployments` and `logs` |
| `apps secrets` | `list`, `set`; project variable operations |
| `db` | `list`, `create`, `destroy`, `import`, `backups` |
| `db backups` | `list`, `create`, `trigger`, `history` |
| `compose` | `analyze`, `deploy` |
| `deploy` | Local path or existing service ID |
| `version` | Binary version |

Use `codedock COMMAND --help` from this binary for resource flags. There is no `whoami` or `servers` command in this Go client. The `delete` and `db backup` spellings in the npm client do not apply here. See [Compose import](/deployments/compose/) before using the Go client's Compose commands.

## Daemon binary: codedockd

The daemon binary starts the server with no command or with `serve`. It also provides direct server-side management using the configured data directory and Docker access.

| Command | Purpose |
| --- | --- |
| `serve` | Start the server |
| `setup` | Interactive setup |
| `reset-password` | Reset administrator credentials |
| `config` | Configuration workflow |
| `deploy` | Git URL, `--image`, `--template` or `--archive` workflow |
| `restart` | Restart the control-plane container |
| `backup` | Create a control-plane data archive |
| `restore FILE` | Restore a control-plane archive |
| `diagnostics` | Server diagnostics |
| `mcp` | MCP stdio integration |
| `version` | Binary version |

Direct resource commands use a colon, not the remote client's space-separated syntax:

```sh
codedockd project:list
codedockd project:show PROJECT_ID
codedockd apps:list
codedockd db:list
codedockd env:list --project PROJECT_ID
codedockd deployment:list --service SERVICE_ID
codedockd deployment:logs DEPLOYMENT_ID
codedockd domain:list --project PROJECT_ID
```

Project, app and database families also provide `create` and `destroy`; variable commands provide `set` and `unset`; domains provide `add` and `remove`. These local operations use host privileges and are not authenticated remote API calls. Run them only against the intended installation and data directory.

## Linux installer wrapper

The installer's `codedockd` host wrapper manages the installed Docker service and forwards supported daemon operations. It adds installation lifecycle commands such as `status`, `logs`, `update`, `downgrade` and `uninstall`. Use its help output for the options on your installed version.

Do not assume host lifecycle commands are subcommands of a downloaded raw daemon binary. See [installation](/getting-started/installation/) for the supported path.

## Recovery

Daemon data archives preserve control-plane state and keys. They are separate from per-database or volume backups. Review the target before restoring or deleting resources, and verify application data after recovery.
