---
title: Deployments API
description: Trigger deployments and inspect their history, logs and metrics.
---

## Trigger a service

`POST /api/services/:serviceId/deploy` triggers deployment for a service. It requires administrator access to the project. Read the response and deployment history to track completion; accepting the request is not a completed rollout.

```sh
curl -X POST "$CODEDOCK_URL/api/services/$SERVICE_ID/deploy"   -H "Authorization: Bearer $CODEDOCK_TOKEN"
```

`POST /api/projects/:id/deploy` triggers the project's deployment workflow.

## Inspect operations

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/services/:serviceId/deployments` | Service deployment history |
| GET | `/api/services/:serviceId/previews` | Pull-request previews |
| GET | `/api/projects/:id/deployments` | Project history |
| GET | `/api/deployments` | Organization history |
| GET | `/api/deployments/:id/logs` | Deployment logs; logs-read scope |
| GET | `/api/deployments/:id/explain` | Failure explanation |
| POST | `/api/deployments/:id/rollback` | Rollback workflow |
| GET | `/api/services/:serviceId/metrics` | Current service metrics |
| GET | `/api/services/:serviceId/metrics/historical` | Historical metrics |
| GET | `/api/services/:serviceId/logs/historical` | Historical runtime logs |

Live service logs use the WebSocket handler at `GET /api/services/:serviceId/logs`. They are not an SSE JSON list. WebSocket authentication and access checks occur in the handler.

## Archives

`POST /api/deploy/archive` accepts a multipart archive workflow. Use the supported tar archive format and fields in the [archive handler](https://github.com/buildwithtechx/codedock/blob/main/internal/handlers/deployments/archive.go); do not assume ZIP uploads are accepted.

Review the [deployment handler](https://github.com/buildwithtechx/codedock/blob/main/internal/handlers/deployments/deployment.go) for query parameters and response models. See [build strategies](/deployments/build-strategies/) for the execution paths.
