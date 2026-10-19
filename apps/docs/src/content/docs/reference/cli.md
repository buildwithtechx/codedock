---
title: CLI Reference
description: Separate the remote client, daemon commands and installer host wrapper.
---

There is one remote client named `codedock`: the Go binary in `cmd/codedock`. Install it with the install script (`curl -fsSL https://get.codedock.run/cli | sh`), via npm (`npm install -g codedock`, which ships this same binary), or from source. `codedockd` is the server daemon. The Linux installer also supplies a host management wrapper. Their command syntax and configuration are different.

## Remote client: codedock

Use the client from your workstation or CI runner. It connects to the daemon over HTTP. Login prompts for the server URL, email and password and saves its configuration to `~/.codedock/config.json`.

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
| `login`, `logout`, `me` | Account and saved client configuration (`whoami` aliases `me`) |
| `status` | Service ID |
| `servers` | `list`, `create`, `delete` |
| `projects` | `list`, `create`, `destroy` |
| `environments`, `env` | `list`, `create`, `destroy`; the `env` alias manages environments, not variables |
| `apps` | `list`, `create`, `destroy`, plus `secrets`, `domains`, `deployments` and `logs` |
| `apps secrets` | `list`, `set`; project variable operations |
| `db` | `list`, `create`, `destroy`, `import`, `backups` |
| `db backups` | `list`, `create`, `trigger`, `history` |
| `compose` | `analyze`, `deploy` |
| `deploy` | Local path or existing service ID |
| `version` | Binary version |

Use `codedock COMMAND --help` for resource flags. See [Compose import](/deployments/compose/) before using the Compose commands.

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
