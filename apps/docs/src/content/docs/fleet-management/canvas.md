---
title: Canvas Topology
description: Explore the services and relationships in an environment.
---

Open a project environment and select Canvas to see its service topology. The view keeps application and database resources together so you can navigate the stack from a shared context.

## Use the map

1. Choose the project and environment you want to inspect.
2. Identify application and database nodes and their current state.
3. Select a resource to open its details, deployment history or configuration.
4. Zoom and pan to follow the relevant part of the environment.

Connections represent relationships known to Codedock. They are not a packet capture, a proof that network traffic is flowing, or a guarantee that every runtime dependency has been discovered. Inspect service variables, Docker networking and logs when diagnosing connectivity.

## Data endpoints

| Endpoint | Purpose |
| --- | --- |
| `GET /api/canvas/projects` | Project summaries for Canvas |
| `GET /api/projects/:id/summary` | A project's Canvas summary |
| `GET /api/environments/:id/canvas` | Environment topology data |

These requests require authentication and access to the underlying project resources. Refresh after changing a deployment or resource configuration.

See [projects and environments](/projects/overview/) and [worker servers](/fleet-management/servers/).
