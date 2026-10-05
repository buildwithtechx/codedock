---
title: Variables and Domains API
description: Project variables, service variables and domain request formats.
---

## Project variables

`GET /api/projects/:id/env` reads a map of variables and requires environment-read scope. `POST` or `PUT` to the same path sets supplied values and requires project administrator access and environment-write scope.

```json
{"variables":{"NODE_ENV":"production","PORT":"3000"}}
```

The setter updates supplied keys individually; it is not an atomic replacement of the whole map.

## Service variables

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/services/:serviceId/variables` | List variables |
| GET | `/api/services/:serviceId/env-suggestions` | Suggested variables |
| POST | `/api/services/:serviceId/variables` | Create a variable |
| PUT | `/api/services/:serviceId/variables/:id` | Update the variable |
| DELETE | `/api/services/:serviceId/variables/:id` | Delete the variable |

Creation and update accept variable fields such as `key` and `value`. Project, service and environment association is resolved or validated by the handler. Writes require administrator access.

```json
{"key":"NODE_ENV","value":"production"}
```

## Domains

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/domains` | List visible domains |
| GET | `/api/services/:id/domains` | List service domains |
| POST | `/api/services/:id/domains` | Create |
| DELETE | `/api/domains/:id` | Delete |
| GET or POST | `/api/domains/:id/verify` | Verify |

Creation accepts `domainName`, optional `redirectTo` and `pathPrefix`. DNS provider-backed creation requires administrator or owner access in addition to service access.

```json
{"domainName":"app.example.com","pathPrefix":"/"}
```

See [domains and TLS](/operations/domains-and-dns/) for DNS and Traefik requirements and [Compose import](/deployments/compose/) for configuration that is not carried over from a Compose file.
