---
title: Projects API
description: Current project and environment routes and request fields.
---

Requests require authentication and permission on the relevant organization or project.

## List projects

`GET /api/projects?organizationId=ORG_ID&page=1&limit=10` requires `organizationId`. The paginated response contains `data.records`, `total`, `page` and `totalPages`.

```sh
curl "$CODEDOCK_URL/api/projects?organizationId=$ORGANIZATION_ID"   -H "Authorization: Bearer $CODEDOCK_TOKEN"
```

## Create a project

`POST /api/projects` accepts the project request, including `name`, optional `description`, `organizationId`, `serverId`, `environmentName` and `environmentType`. Other source-related fields are defined in the current [project DTO](https://github.com/buildwithtechx/codedock/blob/main/pkg/types/project.go).

```json
{"name":"My API","description":"Application stack","organizationId":"ORG_ID"}
```

If omitted for an authenticated user, the organization can be resolved to their default organization. Account limits can affect creation.

## Project operations

| Method | Path | Requirement |
| --- | --- | --- |
| GET | `/api/projects/:id` | Project access |
| DELETE | `/api/projects/:id` | Project owner |
| GET | `/api/projects/:id/environments` | Project access |
| POST | `/api/projects/:id/environments` | Project administrator |
| GET | `/api/projects/:id/apps` | Project access |
| GET | `/api/projects/:id/services` | Alias for app listing |
| GET | `/api/projects/:id/deployments` | Project access |
| POST | `/api/projects/:id/deploy` | Project administrator |
| GET | `/api/projects/:id/summary` | Canvas summary |

Environment creation fields are defined by the [environment handler](https://github.com/buildwithtechx/codedock/blob/main/internal/handlers/projects/environment.go). Delete an environment with `DELETE /api/environments/:id` after reviewing its resources.

See [service creation](/reference/api-services/) and [variables](/reference/api-environment-and-domains/).
