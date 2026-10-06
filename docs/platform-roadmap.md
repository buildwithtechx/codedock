# Platform implementation roadmap

Baseline: Codedock main 1c03759, compared with Openship main 16730d39 on 2026-10-06.
Managed cloud provider: Hetzner. Existing projects and self-hosted installations remain supported.

## Delivery workflow

Each phase ships in focused commits and pull requests. Run Biome and Go formatting, relevant regression checks, CI and Cubic review before considering a phase complete. Merge only when authorized. No feature is complete merely because its settings form or API route exists. Runtime changes require target validation, permission checks, durable progress, cancellation and recovery.

## Phase 1: Reliability foundation

- [x] Backup scheduling applies configured time zones and timeout values.
- [x] S3 uploads apply destination prefixes while old stored object keys remain usable.
- [x] Missing runtimes, unsupported engine commands, failed exits and persistence failures report failure.
- [x] Remote-only archives can be downloaded with the same authorization as local records.
- [ ] Supported volume restores have explicit target validation and interruption controls.
- [ ] Personal API tokens authenticate consistently with the token management UI.
- [x] Empty health states and partial Compose imports use accurate wording.
- [x] Resource creation cards support keyboard activation.
- [ ] Autoscaling has explicit enablement, limits and observable decisions.

## Phase 2: Integrated application deployment

- [ ] One setup flow for Git and images, with project and environment selection.
- [ ] Repository inspection detects framework, package manager, commands, port and output defaults.
- [ ] Users can edit every detected value and see detection failures without losing work.
- [ ] Configure variables, runtime, target and domain in the same flow.
- [ ] Review the exact payload and effects before creating or deploying resources.
- [ ] Live progress separates preparation, build, start, readiness and routing.
- [ ] Retries reconcile existing resources; cancellation preserves the previous deployment.
- [ ] Shared setup supports later Compose, cluster and cloud destinations.

## Phase 3: Complete Compose execution

- [ ] Preserve service environment, interpolation, build context, commands, ports, volumes, networks, health checks and dependencies.
- [ ] Unsupported fields are rejected or surfaced before import; no silent discarding.
- [ ] A saved stack groups its services and deployment configuration.
- [ ] Validate dependency cycles, ports, paths, secrets and target capabilities before apply.
- [ ] Coordinate builds and dependency-aware startup with per-service results.
- [ ] Keep existing workloads on failed builds; report partial activation explicitly.
- [ ] Retrying or redeploying reuses owned resources and preserves data.

## Phase 4: Operational topology

- [ ] Environment selection and graph projection from saved services, domains and bindings.
- [ ] Node selection opens service/database details, logs, metrics and lifecycle controls.
- [ ] Edges distinguish dependencies, variable bindings and routing.
- [ ] Connection edits validate cycles and permissions and persist through existing APIs.
- [ ] Pending changes have review/apply, conflict detection and safe retry.
- [ ] Runtime failures remain visible instead of becoming empty healthy states.
- [ ] Dragging and local layout changes never mutate infrastructure.

## Phase 5: Hetzner cloud

- [ ] Provider interface and Hetzner client wired through injected services.
- [ ] Account/organization-owned server provisioning with encrypted provider credentials.
- [ ] Region, machine type, disk, IP, capacity and quota validation.
- [ ] Idempotent provision, bootstrap, ready, failed, retry and delete operations.
- [ ] Managed servers participate in the same setup and runtime flows as SSH targets.
- [ ] Resource tiers, capacity admission, reservation and resize reconciliation.
- [ ] Billing links paid entitlement to provisioned capacity without affecting self-hosted mode.
- [ ] Explicit confirmation for billable provisioning; no paid infrastructure created by development checks.

## Phase 6: Cluster and runtime orchestration

- [ ] Reviewed private-network and K3s preparation, installation, join and removal operations.
- [ ] Kubernetes runtime adapter with deploy, logs, exec, metrics and lifecycle support.
- [ ] Application placement, instance counts, readiness, routing and persistent storage.
- [ ] Desired/observed reconciliation, upgrade recovery and conflicting-operation locks.
- [ ] PostgreSQL operator setup, replication, credentials, S3 backup and restore into a new database.
- [ ] Redis standalone/replicated modes with engine-specific validation and recovery.
- [ ] Docker replicas are distinct from multi-server scheduling and database replication.
- [ ] Bare runtime support is a separate target capability with service supervision and rollback.

## Phase 7: Backup policy and recovery UX

- [ ] Policies target projects/services and select producer, destination, schedule and retention.
- [ ] Durable runs provide live progress, logs, cancellation and archive verification.
- [ ] Protected records survive retention until their protection expires or is removed.
- [ ] Restore prepare identifies archive, engine, target and destructive effects without applying them.
- [ ] Apply requires an expiring confirmation bound to the reviewed target and archive.
- [ ] In-place and new-target restore modes report interruption and final verification.
- [ ] Pre-deployment policies wait for verified backups before deployment proceeds.

## Additional product scope

- [ ] Migration discovery, adoption, service-selected copy, cutover and rollback.
- [ ] HTTP traffic analytics and actionable operational attention feed.
- [ ] Desktop control-plane lifecycle, connection recovery and update handling.
- [ ] Expanded validated app catalogue and install-time secret/storage/networking setup.
- [ ] Localization and RTL support across shared dashboard components.
- [ ] Mail server and webmail as an optional capability with separate operational requirements.
- [ ] Exhaustive maintained REST request/response schemas.

## Verification boundaries

Source-level and mocked provider checks do not establish live infrastructure reliability. Compose, backups, SSH, Hetzner, K3s and operators need disposable integration targets and recovery exercises. Provider secrets stay out of source control. Schema changes follow the repository initial-schema convention and numbered upgrades when required.
