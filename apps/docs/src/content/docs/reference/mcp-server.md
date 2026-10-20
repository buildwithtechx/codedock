---
title: Model Context Protocol (MCP) Server
description: Connect AI assistants like Claude Desktop, Cursor, and Antigravity to Codedock.
---

Codedock includes native support for the **Model Context Protocol (MCP)**, allowing AI coding assistants and agentic workflows (such as Claude Desktop, Cursor, and Antigravity) to directly query server metrics, inspect containers, trigger deployments, and manage PaaS resources.

## What is MCP?

MCP is an open standard that enables AI models to securely discover and invoke tools provided by external systems. With Codedock's MCP integration, you can ask your AI assistant questions like:

- *"Why did the latest deployment for `api-server` fail?"*
- *"Check the CPU and memory utilization on worker-02."*
- *"Run a SELECT query on my staging PostgreSQL database to check user counts."*
- *"Trigger a redeployment of the web service."*

## Transport Modes

Codedock supports both standard MCP transport protocols:

### 1. Stdio Mode (Desktop & Local CLI)

Used when running your AI client on the same machine or via SSH:

```bash
codedockd mcp
```

### 2. Server-Sent Events (SSE) Mode (Remote API)

Used by web agents, cloud IDEs, and remote integrations:

- **SSE Connection**: `GET /api/mcp/sse` (Authenticated)
- **Message Dispatch**: `POST /api/mcp/messages` (Authenticated)

## Client Configuration Examples

### Claude Desktop

Add Codedock to your `claude_desktop_config.json` (located at `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS or `%APPDATA%\Claude\` on Windows):

```json
{
  "mcpServers": {
    "codedock": {
      "command": "codedockd",
      "args": ["mcp"],
      "env": {
        "CODEDOCK_API_URL": "https://codedock.yourdomain.com",
        "CODEDOCK_TOKEN": "your-personal-access-token"
      }
    }
  }
}
```

### Cursor & AI IDEs

In your IDE's MCP settings, add a new MCP Server:

- **Type**: `command` (or `sse`)
- **Command**: `codedockd mcp`
- **Environment**: Supply `CODEDOCK_API_URL` and `CODEDOCK_TOKEN`.

## Available MCP Tools

Codedock exposes a rich set of operational tools to connected AI models:

| Tool Name | Description |
| --- | --- |
| `list_projects` | Lists all accessible projects and current environment statuses. |
| `get_project_summary` | Retrieves detailed application, database, and domain topologies for a project. |
| `get_service_logs` | Fetches recent stdout/stderr runtime logs for a container with optional tail limits. |
| `get_server_metrics` | Returns real-time CPU, RAM, disk, and load metrics for cluster hosts. |
| `trigger_deployment` | Initiates a fresh container build and rollout for an application service. |
| `rollback_deployment` | Instantly rolls back a service to its previous stable release. |
| `query_database` | Executes read-only SQL queries against project databases via SQL Studio. |
| `get_environment_vars` | Inspects service environment variable keys (sensitive values masked by default). |
| `set_environment_var` | Updates or injects environment variables and triggers a service reload. |

## Security & Access Control

- **Token Scopes**: MCP interactions inherit the exact permissions and project roles of the personal access token used to authenticate.
- **Audit Logging**: Every action invoked by an AI model through MCP is recorded in the Codedock Audit Log with the `mcp-bridge` actor attribution.
- **Server Toggle**: Administrators can disable MCP globally at any time under **Server Settings → AI & MCP → Enable MCP Server**.
