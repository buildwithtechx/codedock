---
title: Services API
description: Create and operate application resources inside an environment.
---

Application resources use `/apps` routes. Database resources have separate `/databases` endpoints.

## Create an application

`POST /api/environments/:id/apps` accepts an application service request. The environment in the route supplies the creation context.

```json
{"name":"web","projectId":"PROJECT_ID","repositoryUrl":"https://github.com/your-org/your-app","branch":"main","rootDirectory":".","runtimeMode":"web","buildEngine":"dockerfile","dockerfilePath":"Dockerfile","internalPort":3000}
```

For a prebuilt container, set `imageRef` instead of a source repository. Configure variables separately. The [service DTO](https://github.com/buildwithtechx/codedock/blob/main/pkg/types/service.go) defines install, build and start commands, static output, health checks and resource limits.

## Operations

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/environments/:id/apps` | List environment applications |
| GET | `/api/apps?organizationId=ORG_ID` | List organization applications |
| GET | `/api/apps/:id` | Inspect an application |
| PUT | `/api/apps/:id` | Update settings; administrator access |
| DELETE | `/api/apps/:id` | Delete; owner access |
| POST | `/api/apps/:id/stop` | Stop the runtime |
| POST | `/api/apps/:id/restart` | Restart the runtime |
| POST | `/api/apps/:id/redeploy` | Redeploy |

Stop, restart and redeploy require administrator access. Updating settings and triggering a deployment are distinct operations.

## Related resources

Volumes, service webhooks and log drains use `/api/apps/:id/volumes`, `/webhooks` and `/log-drains`. List with `GET`, create with `POST`, and delete the specific child ID with `DELETE`. Service variables use `/api/services/:serviceId/variables`.

See [endpoint catalogue](/reference/api-endpoints/) for all child paths and current handler definitions, [deployments](/reference/api-deployments/), and [Compose import](/deployments/compose/).
