---
title: REST API
description: Authentication, response formats and current endpoint references.
---

The daemon exposes `/api/*`. Examples use the unversioned path; `/api/v1/*` is rewritten to the same API. Replace resource placeholders with IDs from your own instance.

## Authentication

Use a valid login access token or a supported project token in `Authorization: Bearer TOKEN`. Project tokens are created through project settings and start with `vsl_tok_`. Scope and project-role checks vary by endpoint; project tokens do not grant administrator access. A login JWT retains the authenticated account role, including administrator or owner privileges where applicable.

```sh
export CODEDOCK_URL="https://pilot.example.com"
export CODEDOCK_TOKEN="YOUR_VALID_TOKEN"
curl "$CODEDOCK_URL/api/organizations"   -H "Authorization: Bearer $CODEDOCK_TOKEN"
```

Profile personal access tokens are a separate stored token type. The current request guard recognizes JWTs and `vsl_tok_` project tokens; do not assume a profile token works as an API bearer token.

## Browser sessions and CSRF

Browser requests use the authentication cookie. Bootstrap `GET /api/auth/csrf`, retain its cookie and send the returned token in `X-CSRF-Token` for unsafe requests. Auth routes still require CSRF handling when a bearer header is present. Bearer requests to other routes bypass that CSRF check.

For sign-in, send `email` and `password` to `POST /api/auth/signin`; retain tokens or cookies from the response. Do not put passwords in URLs or commit credentials to a script.

## Response envelope

Most JSON handlers use:

```json
{"status":"success","message":"Operation successful","data":{},"path":"/api/example"}
```

Errors use `status: "error"` and a message. Middleware errors can use a different shape. Paginated `data` contains `records`, `total`, `page` and `totalPages`. Some responses are raw JSON, downloads, WebSocket upgrades or `204` without a body; check the endpoint rather than always assuming an envelope.

| Status | Typical meaning |
| --- | --- |
| 400 | Invalid fields or missing parameters |
| 401 | Missing or invalid authentication |
| 403 | Insufficient permission, scope or CSRF token |
| 404 | Missing resource or disabled capability |
| 409 | Resource conflict |
| 429 | Rate limit exceeded |
| 500 | Operation failed on the server |

## References

- [All registered endpoints](/reference/api-endpoints/)
- [Access and token restrictions](/reference/api-access/)
- [Projects and environments](/reference/api-projects/)
- [Application services](/reference/api-services/)
- [Deployments and metrics](/reference/api-deployments/)
- [Databases and data operations](/reference/api-databases/)
- [Variables and domains](/reference/api-environment-and-domains/)

Billing endpoints are enabled only in cloud mode. Self-hosted instances return `404` for billing. Health, authentication, signed webhooks and WebSocket handlers have their own access behavior.

The endpoint catalogue identifies route wiring and links to the current handler definitions. It is not an OpenAPI contract with exhaustive request and response schemas; that remains separate documentation work.
