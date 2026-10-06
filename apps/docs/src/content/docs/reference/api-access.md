---
title: API Access
description: Use login and project tokens supported by the current request guard.
---

Use `Authorization: Bearer TOKEN` with a valid JWT access token, a `vpt_` personal token or a `vsl_tok_` project token. Project tokens are managed in project settings through `/api/projects/:projectId/tokens`. Project binding is enforced on routes with project or service authorization checks; it is not a universal boundary. Scope-only routes, such as server management with `server:write`, can authorize instance-wide operations without checking the token project. Grant only the scopes needed for the intended endpoints.

The request guard checks authentication, account state, configured IP restrictions, role and route-specific scopes. Project tokens cannot perform administrator or owner-only operations. Read the relevant endpoint's requirements before using a token in automation.

## Create a project token

Use the project's token settings as a project administrator. Copy the generated value when it is shown, store it securely and revoke it when no longer needed. Creation and deletion require the appropriate project role and environment-write scope; listing requires environment-read scope.

## Token types

The profile token endpoint creates `vpt_` personal bearer tokens with their own access level and project selection. There is no generic `ap_` key or `X-API-Key` authentication path documented for this implementation.

For scripts, use the remote CLI login flow or a project token accepted by the target endpoint. Avoid recording token values in CI logs.

## Examples

```sh
curl "$CODEDOCK_URL/api/projects/$PROJECT_ID"   -H "Authorization: Bearer $CODEDOCK_TOKEN"
```

Browser cookie requests have CSRF requirements. See [API overview](/api/) for authentication and response behavior and [endpoint catalogue](/reference/api-endpoints/) for route-specific middleware wiring.

## Personal API tokens

Create personal tokens in profile token management. Send the returned `vpt_` credential as `Authorization: Bearer TOKEN`. Credentials are stored as SHA-256 hashes and shown once. The API returns the token metadata in `data.token` and the one-time credential in `data.plain`.

**Read** tokens allow GET, HEAD and OPTIONS. **Read and write** tokens additionally permit mutations within their owner's current resource permissions. Expired, deleted tokens and tokens belonging to inactive or deleted users are rejected. The UI's **No expiration** choice creates a token without an expiry; revoke it when no longer needed.

**Specific projects** limits access to the selected projects and their services, environments, databases, deployments and backups. Unscoped collection endpoints are denied for these tokens; request a selected project's resources through its scoped routes. References to other projects in query parameters or JSON bodies are rejected. Project scope narrows the owner's existing permissions and grants no additional membership.

Account changes, token management, interactive terminals and instance role-restricted operations require a user session. Existing project API tokens retain their separate project and operation scopes.
