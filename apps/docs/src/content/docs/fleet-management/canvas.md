---
title: Canvas Topology
description: Explore the services and relationships in an environment.
---

Open a project and select Canvas to see summaries of its default environment. The current page does not offer an environment selector. The view keeps application and database resources together so you can navigate the stack from a shared context.

## Use the map

1. Choose the project whose default environment you want to inspect.
2. Identify application and database nodes and their current state.
3. Inspect the names, engine types and status summaries shown on nodes. Node selection does not currently open resource details; use the project resource pages for configuration and deployment history.
4. Zoom and pan to follow the relevant part of the environment.

Connections are drawn manually in the Canvas UI. Codedock does not discover service dependencies for this map, and drawn edges are unverified local UI state rather than persisted network configuration. They are not a packet capture, a proof that network traffic is flowing, or a guarantee that every runtime dependency has been discovered. Inspect service variables, Docker networking and logs when diagnosing connectivity.

## Data endpoints

| Endpoint | Purpose |
| --- | --- |
| `GET /api/canvas/projects` | Project summaries for Canvas |
| `GET /api/projects/:id/summary` | A project's Canvas summary |
| `GET /api/environments/:id/canvas` | Environment topology data |

These requests require authentication and access to the underlying project resources. Refresh after changing a deployment or resource configuration.

See [projects and environments](/projects/overview/) and [worker servers](/fleet-management/servers/).
