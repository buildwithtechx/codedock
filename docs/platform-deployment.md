# Application, Compose and topology workflows

Application setup accepts Git repositories or container images and selects a project and environment. Repository inspection proposes framework, package manager, commands, port and static output defaults; every proposal remains editable. Connected repository credentials belong to the authenticated user. Inspection errors preserve the draft. Review shows the submitted configuration and variables before saving or deploying.

Application creation uses a stable request ID, so retrying an interrupted save reconciles the same application. Deployment progress includes preparation, build, startup, readiness and routing. A replacement keeps the previous containers until candidates are ready. Cancellation or failed readiness restores the previous runtime and metadata. Services with persistent mounts pause the previous containers during replacement. A durable journal recovers interrupted replacements when the daemon restarts; unresolved cleanup blocks another replacement.

## Compose stacks

Project Compose setup saves a complete resolved stack configuration encrypted in the database. Review validates supported fields, dependencies, interpolation, published ports and ownership before saving. Deployment prepares every image before activation and delegates dependency and health handling to Docker Compose. Each service has an observed container state, health and exit code.

Named volumes and networks belong to the stack and survive retries. Build and image preparation failures leave existing workloads running. Activation failure or cancellation can leave a partially updated stack; the interface reports this and allows observation and retry. Removing a service from the saved configuration does not delete its running container or data.

Uploaded relative build contexts require an explicit public HTTPS Git repository, branch and Compose directory. Supported repository hosts are GitHub, GitLab and Bitbucket. Host bind mounts, anonymous volumes, external resources, file-based environment/secrets, privileged containers, host networking and cluster placement are rejected before saving. Private repository build contexts require a future credential adapter.

The container image includes Docker Compose. Binary installations need the Docker Compose CLI plugin on the daemon host. Stack execution uses the same Docker endpoint as the daemon. No new required environment variables were introduced.

## Operational canvas

Select an environment to view applications, databases, stack services, domains, variable bindings and dependencies. Resource selection shows observed runtime state, logs, metrics and applicable lifecycle controls. Stack operations remain at stack level. Dependency edits require review and apply; the server checks project permissions, environment membership, cycles and a revision of the saved configuration. Conflicts preserve the pending draft for review and retry.

Database variable bindings use stable resource IDs and replace the selected application variable on apply. They take effect on its next deployment. The reviewed environment revision prevents applying a binding against changed configuration. Dragging nodes only changes the local layout.

Application setup inherits the project Docker destination; SSH Docker application deployments are supported. Kubernetes and native destinations are configured separately in application build settings after creation. Compose stack execution still requires local Docker. Managed Hetzner provisioning remains unimplemented.

## Verification

Run `npm run fmt`, `npm run typecheck`, dashboard tests and `go test ./...` for the combined batch. Disposable Docker recovery checks are opt-in: set `CODEDOCK_TEST_LIVE_DOCKER=1` and run `go test ./internal/engine/compose ./internal/engine/deploy -run TestLive -v`. They create isolated containers, a temporary network and a temporary volume, then remove only those fixture resources.

Compose behavior follows the official [configuration reference](https://docs.docker.com/reference/cli/docker/compose/config/), [activation reference](https://docs.docker.com/reference/cli/docker/compose/up/) and [Git build context reference](https://docs.docker.com/build/concepts/context/#git-repositories).
