---
title: Docker Compose Import
description: Analyze Compose definitions and configure the resulting Codedock resources.
---

Compose import translates a definition into database and application resources. It is not equivalent to running `docker compose up` with the original file.

## Analyze first

Use `POST /api/compose/analyze` with JSON fields `composeContent` and `projectId`. The response contains `appServices` and `databases`; it does not deploy them.

```json
{"projectId":"PROJECT_ID","composeContent":"services:\n  cache:\n    image: redis:7\n"}
```

The parser identifies database engines from image names and otherwise creates an application request. Image tags provide database versions. A `build` field selects the Dockerfile build strategy, but does not carry a complete build context into the created service.

## Import resources

`POST /api/compose/deploy` accepts multipart fields `file` (the Compose YAML) and `projectId`. The upload limit is 50 MB. Use a project you administer.

```sh
curl "$CODEDOCK_URL/api/compose/deploy"   -H "Authorization: Bearer $CODEDOCK_TOKEN"   -F "projectId=$PROJECT_ID" -F "file=@compose.yaml"
```

The endpoint creates resources and returns a count. Database creation can start containers immediately during import. Review the Analyze response before calling Deploy; import is not a deferred preview. It does not preserve the entire Compose runtime configuration. Import is not atomic: if a later creation fails, earlier resources can already exist. Check for duplicates before retrying.

## Environment and override rules

The current parser does not translate Compose `environment`, `env_file`, variable interpolation, `ports`, `volumes`, `depends_on`, or override-file merging into Codedock resource settings. There is no Compose environment precedence applied during import.

Set variables, credentials, ports, volumes and dependencies explicitly on the resulting resources. Do not assume a `.env` beside the uploaded file is loaded. Never upload production secrets expecting Compose interpolation to protect or inject them.

For a full Compose stack whose configuration must remain intact, operate it directly with Docker Compose and manage its configuration separately. Codedock import is a resource conversion workflow.

See [environment variables](/deployments/environment-variables/) and [build strategies](/deployments/build-strategies/).
