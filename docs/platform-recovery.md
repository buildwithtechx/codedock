# Backup recovery and cluster setup

## Backup policies and operations

Backup policies keep the existing database or named-volume producer, local/S3 destination, schedule, timezone, timeout and retention settings. Service policies can require a verified backup before an application deployment. A disabled required policy or failed producer/verification blocks that deployment. The hook covers applications, archives, Compose stack activation and SSH Docker deployments. A database policy without an application owner gates deployments in its project; an application policy gates only that application. Remote producers use the project Docker engine. Retention limits can be cleared. Control-plane snapshots use pg_dump against the control-plane database.

Manual policy runs use durable operations with phases, bounded logs and cancellation. Daemon restart marks active operations interrupted and unfinished backup records failed. Scheduled jobs use the same durable, cancellable operation workflow. The saved policy owner must retain active project admin access at scheduling and execution. Manual runs appear under Backups; operation history belongs to the requesting user.

Completed records can be protected for thirty days through the interface, or until a chosen Unix expiry through the API. Protection and retention/manual deletion claims are mutually exclusive database updates. A policy containing protected or active records cannot be deleted. Cleanup failures preserve a failed record for inspection rather than reporting a usable archive.

Archives have a SHA256 checksum, byte count and verification timestamp. Verification checks stored bytes; it does not prove an application-consistent snapshot or successful recovery of every database object. Named-volume backups can include files being written by running applications. Quiesce workloads when application consistency is required.

## Reviewed restore

Prepare checks archive availability, checksum, source engine and target runtime identity. Database targets must use the same engine in the same authorized project and run on local Docker. The interface lists eligible databases by name. A separate target must already exist; this workflow does not provision it automatically.

Named-volume targets must already exist, belong to the registered service/container and have no running consumers. Apply stages the exact reviewed archive in a private temporary file and checks compressed tar structure before extraction. Paths escaping the volume, links and special files are rejected. Existing symbolic links in the destination also prevent extraction. Archives are limited to 1 GiB compressed and 10 GiB declared file content.

Apply requires an opaque confirmation that expires after ten minutes and matches the operation owner, archive checksum and target snapshot. A changed archive, recreated target or changed policy requires another review. Target locks prevent conflicting restores and lifecycle operations.

Database restore completion checks command exit and target runtime/health. Volume completion checks extraction exit and input delivery. Interrupting an overwrite can leave partial data; interruption is not rollback. The dialog keeps progress, errors and cancellation visible while execution runs. The older direct restore endpoints return HTTP 410.

## Cluster infrastructure

The project Clusters tab lets instance administrators review installation, worker joins, release upgrades and removal on registered SSH servers. Each node requires its private IPv4 address, interface and independently verified SHA256 SSH host fingerprint. Preflight checks Linux, systemd, required tools, capacity, assigned private IP and existing K3s ownership. Addresses overlapping the default K3s pod/service networks are rejected.

The topology has one control-plane server and additional agents. It provides no control-plane failover. An exact K3s release is required. The downloaded official installer is bound to review by checksum, and join tokens are generated automatically and encrypted in storage. Upgrade reviews reject downgrades and skipped minor releases. Review effects identify installation changes and potentially destructive removal; apply requires an expiring confirmation.

Only clusters bearing the matching Codedock ownership marker can be changed or uninstalled. Server and cluster locks prevent conflicting operations. Removal processes agents before the control plane; retrying removal tolerates nodes already removed. Cancellation attempts to terminate the owned remote process group and reports possible partial installation. An upgrade writes an encrypted journal before changing releases and verifies every node reports the reviewed version and readiness. Failed or interrupted upgrades resume forward to the reviewed release. Explicitly cancelled upgrades remain paused for a reviewed recovery action. Installation/join interruptions remain visible for reviewed retry.

Administrators must restrict network access according to the [K3s requirements](https://docs.k3s.io/installation/requirements). Preflight verifies address assignment; it does not configure firewall rules or prove connectivity between every node. Release installation follows the [official K3s installation workflow](https://docs.k3s.io/quick-start).

Applications can select a ready project cluster from build settings through an expiring reviewed operation. Normal Git/image and archive deployments then use that destination. Git builds publish unique image tags to a project registry before applying cluster resources. Deployment manifests include owned namespaces, secret-backed environment variables, node placement, replicas, readiness, Services/Ingress and explicit persistent claims. Workers do not require an HTTP listener. Multiple replicas with storage require ReadWriteMany claims.

Cluster workloads support observed pod state, live metrics when the metrics API is available, logs, command execution, stop and restart. Dependency readiness and canvas observations dispatch to the selected runtime. Encrypted recovery journals preserve previous resources before mutation; startup and failed-rollout recovery restore configuration while retaining persistent data. A minute-based reconciler reapplies the last committed encrypted manifest under service/cluster locks, restores replica drift, respects stopped workloads and skips newly configured destinations until deployment. Missing or incomplete metrics remain an explicit error. Existing claim class, access modes and capacity cannot change as part of an application rollout.

## Cluster databases

Ready clusters expose reviewed operator installation and database creation in the Clusters tab. Operator reviews download exact cert-manager, CloudNativePG and Barman plugin release manifests, bind their checksums to confirmation and refuse existing foreign-owned resources. Installation proceeds in dependency order and waits for controller readiness. PostgreSQL uses exact image versions, generated encrypted credentials, private cluster DNS and one or three instances. Three instances require distinct nodes and persistent storage on each node. Database cards query observed readiness instead of displaying a saved READY state after runtime failure.

PostgreSQL S3 destinations supply encrypted credentials to owned Kubernetes Secrets. The Barman plugin archives WAL, schedules daily backups and supports reviewed manual backups with progress and cancellation. Recovery requires a completed base backup and creates a separate database with new credentials and a separate backup prefix. A chosen UTC recovery time is bound to review; unavailable WAL causes recovery/readiness failure. Source data is preserved. Upstream configuration follows the [Barman plugin usage](https://cloudnative-pg.io/plugin-barman-cloud/docs/usage/) and [installation requirements](https://cloudnative-pg.io/plugin-barman-cloud/docs/installation/).

Redis supports one instance or a fixed primary with two replicas. It uses generated credentials, AOF persistence and readiness checks that include replica synchronization. Three-instance mode requires three nodes. This mode does not provide Sentinel failover or sharding. Reviewed recovery reapplies the owned workload against retained data; PostgreSQL S3 recovery is unavailable for Redis. Failed or interrupted database operations remain visible and can be reviewed for recovery; cancellation preserves persistent claims and can leave partial resources.

## Native releases

The bare runtime selects a registered SSH server with a pinned host fingerprint and an HTTPS executable artifact bound to its SHA256 checksum. It executes `./app` with reviewed arguments as an isolated service user under systemd. Environment variables are resolved automatically. A private persistent directory is exposed as `CODEDOCK_DATA_DIR`; CPU/memory limits, service logs, command execution and lifecycle controls are supported. HTTP processes use an unprivileged native port and readiness URL; workers require a stable active service. Managed Docker domains, Docker volume mounts and multiple replicas are rejected.

Activation records the previous release before switching the owned service. Failed activation and startup recovery restore the prior release while retaining persistent data. Server locks prevent conflicting changes. Native CPU usage is explicitly unavailable rather than reported as zero. No Hetzner servers are provisioned and no new environment variables are required.

## Verification

Repository and service checks cover encrypted payloads/tokens, review expiry, stale revisions, target conflicts, protection/deletion races, interrupted-state recovery, required-backup failures and unsafe archive input. The restore dialog check exercises prepare, explicit approval and submission of the exact opaque confirmation.

Live K3s installation, upgrades, network restrictions and recovery need disposable registered servers. No remote cluster installation was performed during this batch. Live operator/database replication and native release recovery were not exercised. Migration cutover/rollback, traffic analytics, desktop lifecycle, catalogue expansion and exhaustive REST schemas remain unchecked in the roadmap.
