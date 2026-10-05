---
title: API Access
description: Use login and project tokens supported by the current request guard.
---

Use `Authorization: Bearer TOKEN` with a valid JWT access token or a `vsl_tok_` project token. Project tokens are managed in project settings through `/api/projects/:projectId/tokens`. Project binding is enforced on routes with project or service authorization checks; it is not a universal boundary. Scope-only routes, such as server management with `server:write`, can authorize instance-wide operations without checking the token project. Grant only the scopes needed for the intended endpoints.

The request guard checks authentication, account state, configured IP restrictions, role and route-specific scopes. Project tokens cannot perform administrator or owner-only operations. Read the relevant endpoint's requirements before using a token in automation.

## Create a project token

Use the project's token settings as a project administrator. Copy the generated value when it is shown, store it securely and revoke it when no longer needed. Creation and deletion require the appropriate project role and environment-write scope; listing requires environment-read scope.

## Token types

The profile token endpoint creates a different token type. The current guard does not resolve those profile tokens as project bearer credentials. There is no generic `ap_` key or `X-API-Key` authentication path documented for this implementation.

For scripts, use the remote CLI login flow or a project token accepted by the target endpoint. Avoid recording token values in CI logs.

## Examples

```sh
curl "$CODEDOCK_URL/api/projects/$PROJECT_ID"   -H "Authorization: Bearer $CODEDOCK_TOKEN"
```

Browser cookie requests have CSRF requirements. See [API overview](/api/) for authentication and response behavior and [endpoint catalogue](/reference/api-endpoints/) for route-specific middleware wiring.
