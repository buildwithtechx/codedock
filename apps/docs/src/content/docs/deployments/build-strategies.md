---
title: Build Strategies
description: Choose Dockerfile, automatic Nixpacks builds or Cloud Native Buildpacks.
---

Choose a build engine in application settings. A prebuilt image avoids a source build; pin its tag or digest for reproducibility.

| Selection | Current execution path |
| --- | --- |
| Dockerfile | Docker builds the configured Dockerfile and context |
| Nixpacks | Runs the Nixpacks CLI in a builder container |
| Railpack | Currently uses the same Nixpacks path; not a separate Railpack CLI |
| Buildpacks | Runs `pack build` in a builder container |
| Auto | Chooses a source build strategy from service settings and detected files |

## Dockerfile

Use a Dockerfile when you need explicit dependencies, image layers or startup behavior. Configure the repository, root directory, Dockerfile path, exposed service port and start command as needed. Your application must listen on an interface reachable from outside its container.

```dockerfile
FROM node:24-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev
COPY . .
CMD ["node", "index.js"]
```

## Nixpacks and Railpack selection

The builder currently invokes `nixpacks build /app --name IMAGE`. Service install, build and start overrides become `--install-cmd`, `--build-cmd` and `--start-cmd`. Environment variables are passed using `--env`.

The Railpack label does not imply Railpack configuration files or flags are supported. Use a Dockerfile for requirements the Nixpacks path cannot express.

## Buildpacks

The Buildpacks selection runs `pack build` with a builder image. Server configuration can override the pack and builder images. Variables are passed to pack; command override behavior is different from the Nixpacks branch. Prefer an explicit Dockerfile if your workflow depends on exact command overrides.

## Builder failure and fallback

If a builder container fails, Codedock attempts a synthesized Dockerfile based on the detected stack. Read the deployment logs to identify which path produced the final image. A fallback success does not establish that all builder-specific configuration was applied.

Docker access, registry connectivity, disk space and compatible builder images are required. The daemon is a Go binary, but building and running applications requires container infrastructure.

## Troubleshooting

1. Verify the root directory and Dockerfile path against the repository.
2. Check install and build logs before changing the runtime start command.
3. Confirm your build can access its registry and dependency sources.
4. Match the configured service port to the listening port.
5. Use a Dockerfile when automatic detection produces the wrong result.
