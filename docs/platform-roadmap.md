# Platform implementation roadmap

Baseline: Codedock main 1c03759, originally compared with Openship main 16730d39 on 2026-10-06.
Managed cloud provider: Hetzner. Existing projects and self-hosted installations remain supported. Upstream managed cloud is SaaS relay with Oblien quota/billing, not Hetzner; Hetzner remains separate Codedock scope rather than upstream parity.

Parity audit refreshed against Openship main `93cdaebd` and Codedock `852d374` on 2026-10-07: [Backend and UI/UX comparison](openship-parity-audit.md). Checked items describe their stated scope, not exhaustive Openship parity or live infrastructure certification.

## Delivery workflow

Phases 2, 3 and 4 ship together in focused commits and one pull request, with one combined CI and Cubic review cycle. Later phases can follow the same batching approach when their runtime adapters are ready. Run Biome and Go formatting and relevant regression checks before committing. The current batch uses incremental commits and pushes without creating a PR; CI and Cubic review remain pending until a PR is authorized. Merge only when authorized. No feature is complete merely because its settings form or API route exists. Runtime changes require target validation, permission checks, durable progress, cancellation and recovery.

## Phase 1: Reliability foundation

- [x] Backup scheduling applies configured time zones and timeout values.
- [x] S3 uploads apply destination prefixes while old stored object keys remain usable.
- [x] Missing runtimes, unsupported engine commands, failed exits and persistence failures report failure.
- [x] Remote-only archives can be downloaded with the same authorization as local records.
- [x] Supported volume restores have explicit target validation and interruption controls.
- [x] Personal API tokens authenticate consistently with the token management UI.
- [x] Empty health states and partial Compose imports use accurate wording.
- [x] Resource creation cards support keyboard activation.
- [x] Autoscaling has explicit enablement, limits and observable decisions.

Implementation details and current target boundaries: [Application, Compose and topology workflows](platform-deployment.md).

## Phase 2: Integrated application deployment

- [x] One setup flow for Git and images, with project and environment selection.
- [x] Repository inspection detects framework, package manager, commands, port and output defaults.
- [x] Users can edit every detected value and see detection failures without losing work.
- [x] Configure variables, web/worker mode and domain in the application setup flow.
- [x] Select Docker/SSH, cluster or native destination before the first deployment within that same flow.
- [x] Review the exact payload and effects before creating or deploying resources.
- [x] Live progress separates preparation, build, start, readiness and routing.
- [x] Retries reconcile existing resources; cancellation preserves the previous deployment.
- [x] Shared setup supports later Compose, cluster and cloud destinations.

## Phase 3: Complete Compose execution

- [x] Preserve service environment, interpolation, build context, commands, ports, volumes, networks, health checks and dependencies.
- [x] Unsupported fields are rejected or surfaced before import; no silent discarding.
- [x] A saved stack groups its services and deployment configuration.
- [x] Validate dependency cycles, ports, paths, secrets and target capabilities before apply.
- [x] Coordinate builds and dependency-aware startup with per-service results.
- [x] Keep existing workloads on failed builds; report partial activation explicitly.
- [x] Retrying or redeploying reuses owned resources and preserves data.

## Phase 4: Operational topology

- [x] Environment selection and graph projection from saved services, domains and bindings.
- [x] Node selection opens service/database details, logs, metrics and lifecycle controls.
- [x] Edges distinguish dependencies, variable bindings and routing.
- [x] Connection edits validate cycles and permissions and persist through existing APIs.
- [x] Pending changes have review/apply, conflict detection and safe retry.
- [x] Runtime failures remain visible instead of becoming empty healthy states.
- [x] Dragging and local layout changes never mutate infrastructure.

## Phase 5: Hetzner cloud

- [x] Provider interface and Hetzner client wired through injected services.
- [x] Account/organization-owned server provisioning with encrypted provider credentials.
- [x] Region, machine type, disk, IP, capacity and quota validation.
- [x] Idempotent provision, bootstrap, ready, failed, retry and delete operations.
- [x] Managed servers participate in the same setup and runtime flows as SSH targets.
- [x] Resource tiers, capacity admission, reservation and resize reconciliation.
- [x] Billing links paid entitlement to provisioned capacity without affecting self-hosted mode.
- [x] Explicit confirmation for billable provisioning; no paid infrastructure created by development checks.

Phase 5 is implemented within these boundaries: provisioning, resize and delete run through expiring-confirmation operations with durable progress, cancellation and retry; provider credentials are organization-owned and encrypted; quotas gate admission; cloud mode requires a pro subscription for paid capacity while self-hosted mode is unaffected. Managed servers land as ordinary SSH server rows, so existing Docker/SSH deployment, metrics, logs and terminal flows apply. Unit coverage uses a fake provider and in-memory database; live Hetzner provisioning, cloud-init Docker bootstrap, resize power-cycling and SSH reachability still require a disposable account and have not been exercised. Dashboard API contracts are in place; a managed-server setup wizard remains follow-up work.

## Phase 6: Cluster and runtime orchestration

- [x] Reviewed private-network and K3s preparation, installation, join and removal operations.
- [x] Kubernetes runtime adapter with deploy, logs, exec, metrics and lifecycle support.
- [x] Application placement, instance counts, readiness, routing and persistent storage.
- [x] Observed workload/pod state, conflicting-operation locks and encrypted deployment journals with failed-rollout/startup recovery.
- [x] Automatic desired/observed reconciliation and cluster-upgrade recovery.
- [x] PostgreSQL operator setup, replication, credentials, S3 backup and restore into a new database.
- [x] Redis standalone/replicated modes with engine-specific validation and recovery.
- [x] Docker replicas are distinct from multi-server scheduling and database replication.
- [x] Bare runtime support is a separate target capability with service supervision and rollback.

### Remaining cluster/runtime parity and integration

- [x] Cluster databases appear in topology and support reviewed application bindings and connection replacement.
- [x] Cluster database settings, stop/retain, exact-name deletion and owned PVC cleanup.
- [x] PostgreSQL synchronous durability, observed roles/volumes and authenticated application-namespace connection verification.
- [x] Redis operator-backed sharding/failover, S3 snapshots and separate-target data restore.
- [x] Shared storage provisioning, file-volume lifecycle, external snapshots and separate-volume restore.
- [x] Multi-control-plane quorum, private-network preparation/firewall checks and cross-node DNS/service probes.
- [x] Native source builds, stack/toolchain setup and managed HTTP routing.
- [x] SSH Docker logs, metrics, terminal, dependencies and canvas dispatch through the remote engine consistently.
- [x] Autoscaling eligibility follows the selected runtime; Docker autoscaling remains local-only by eligibility.
- [x] Cluster live log/terminal transport, observed instance graph and human-readable resource controls.

## Phase 7: Backup policy and recovery UX

- [x] Policies target authorized projects/services and select supported database/named-volume producers, local/S3 destinations, schedule and retention.
- [x] Manual backup runs provide durable progress, logs, cancellation and archive verification.
- [x] Scheduled runs use the same durable, cancellable operation workflow and revalidate the policy owner.
- [x] Protected records survive retention until their protection expires or is removed.
- [x] Restore prepare identifies archive, engine, target and destructive effects without applying them.
- [x] Apply requires an expiring confirmation bound to the reviewed target and archive.
- [x] In-place and new-target restore modes report interruption and final verification.
- [x] Application and archive deployments wait for required verified backups before proceeding.
- [x] Compose stack activation and SSH Docker deployments share the required-backup hook.

The originally scoped Phase 6 and 7 items are implemented within the boundaries below; the renewed comparison found additional parity work and integration gaps. Full Go tests, root typecheck and formatting passed, including SSH Docker operational dispatch with container-ownership checks and runtime-aware autoscaling eligibility. PostgreSQL operator installation, replicated databases, reviewed S3 recovery, Redis AOF recovery and supervised native releases are wired through the dashboard and API. These checks cover repository behavior and mocked runtime commands; live K3s, operator and native-server installation/recovery exercises still require disposable targets. Redis runs standalone, fixed-primary replicas or sharded Cluster mode with per-shard failover by primary restart. Native targets build Git sources with go/node/python/static toolchains or execute verified standalone artifacts on a single SSH server, with managed HTTP routing for web services.

### Remaining backup parity

- [x] Remote SSH and native restore producers/targets with durable interruption handling.
- [x] Project-wide multi-service policy batches and inherited service shortcuts.
- [x] SFTP, incremental storage, quiescing and reviewed file/custom-command producers.
- [x] Retained-image deployment rollback and restore-source retention protection through recovery completion.

Current backup and cluster implementation boundaries, including scheduled-run cancellation, deployment coverage and live-target validation: [Backup recovery and cluster setup](platform-recovery.md).

## Additional product scope

- [x] Migration full lifecycle: migration-only sources, scan/stream, secret reveal, adopt/reimport, preview, move/copy, cutover, cancel, resume and target cleanup. Bundle export/import remains for whole-instance backup; service migration runs through the new lifecycle with review prompts, confirmation tokens and rollback.
- [x] HTTP traffic analytics and actionable operational attention feed. Traefik JSON access logs feed per-minute buckets with domain attribution, opt-in path collection, IP-level geography (country resolution beyond private ranges needs a GeoIP database), deployment stats and dashboard rollup; the attention feed evaluates deployments, workloads, backups, migrations and quota with acknowledge/resolve plus restart, redeploy, resume and retry actions.
- [x] Desktop control-plane lifecycle, connection recovery and update handling. The Tauri shell supervises the local daemon sidecar with backoff restart, health probing and start/stop/restart/port commands; daemon URL/port persist in the Tauri store with the API token in the OS keyring; the updater plugin checks, downloads and installs signed releases with progress events. Release bundling of the sidecar binary, updater signing keys and notarization CI remain pending.
- [x] Expanded validated app catalogue and install-time secret/storage/networking setup. All 26 embedded templates are schema-validated at load (images, ports, named volumes, dependencies, secret references); one-click installs resolve generated/provided secrets, host ports, Traefik domains and prefixed volumes into a digest-pinned preview before saving and activating a compose stack. Database templates keep the existing database provisioning path.
- [ ] Localization and RTL support across shared dashboard components.
- [ ] Mail server and webmail as an optional capability with separate operational requirements.
- [x] Exhaustive maintained REST request/response schemas. All 337 routes carry request/response entries in `internal/http/apischema` with an OpenAPI 3.1 emitter served at `GET /api/docs/openapi.json`; the suite fails on missing, orphaned or duplicate entries and diffs `docs/api/openapi.json` as a golden file. The initial `docs/api/openapi.json` was generated through a byte-faithful Node port of the emitter because `go.exe` is blocked machine-wide by Application Control policy; the first `go test ./internal/http/` run with a working toolchain confirms the golden file and coverage.

## Verification boundaries

Source-level and mocked provider checks do not establish live infrastructure reliability. Compose, backups, SSH, Hetzner, K3s and operators need disposable integration targets and recovery exercises. Provider secrets stay out of source control. Schema changes follow the repository initial-schema convention and numbered upgrades when required.
