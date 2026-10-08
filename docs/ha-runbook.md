# Codedock HA runbook

Reference cell: `compose.ha.yml` runs Postgres 16, two `codedockd` API replicas behind an edge Traefik on host port 8090, and the daemon-managed workload Traefik on host ports 80/443. All replicas share one Docker daemon on a single host.

## Environment contract

Every replica must share the same control-plane database and vault key:

- `CODEDOCK_DATABASE_URL` points at Postgres. Unset means single-node embedded mode and must never be relied on in HA.
- `CODEDOCK_VAULT_KEY` is 64 hex characters, identical on all replicas. Without it each node generates its own key file and encrypted columns become unreadable across nodes.
- `CODEDOCK_PG_PASSWORD` seeds the cell Postgres. Keep it URL-safe because it is interpolated into `CODEDOCK_DATABASE_URL`.

JWT, refresh and telemetry secrets live in the `self_hosted_config` table and converge automatically. Explicit `CODEDOCK_JWT_SECRET`, `CODEDOCK_REFRESH_SECRET` and `CODEDOCK_WILDCARD_DOMAIN` values still take precedence when set.

## How leadership works

One elector key (`scheduled-execution`) holds a Postgres session advisory lock. Exactly one replica is leader:

- The leader runs cron jobs, backup schedules, Docker cleanup, disk checks and the update checker. Entries reconcile from the database every 30 seconds, so registrations made through any replica converge without restarts.
- Followers serve API traffic and keep their schedulers stopped.
- On leader loss the lock releases with the session and a contender takes over within seconds. In-flight cron and backup executions are not migrated; the next schedule fires on the new leader.

Operations, autoscaling and attention need no leader. Operation apply and autoscale reservations are compare-and-swap updates in Postgres, so concurrent replicas stay correct and exactly one wins. Cancelling an operation from any replica records `CANCELLING`; the executor observes it at its next progress report and interrupts. Attention evaluates on read with a per-node throttle, which is safe to run concurrently.

## What stays node-local

- Backup archives and staging under `CODEDOCK_DATA_DIR/backups`. The reference cell shares one volume; a multi-host cell must use shared storage or S3 destinations with local copies disabled.
- Workload execution itself. Cron and backup runs exec into containers on the shared Docker daemon, so replicas must reach the same daemon.
- Traefik access logs, TSDB and Loki data. Metrics and log workers run on every replica; duplicate scrapes are harmless.
- The daemon-managed `codedock-traefik` container is a singleton per Docker daemon. Concurrent ensures from replicas are idempotent and converge.

## Postgres failover

The reference cell (`compose.ha.yml`) runs a single Postgres. The production cell (`compose.ha-patroni.yml`) runs etcd ×3, Patroni ×3 and HAProxy: the daemon connects through HAProxy port 5432, which routes to the current primary via Patroni REST health checks. Replication is synchronous to one standby in strict mode, so an acknowledged write survives a primary loss.

Boot the production cell with the same variables as the reference cell plus `PATRONI_SUPERUSER_PASSWORD` and `PATRONI_REPLICATION_PASSWORD`:

```sh
docker compose -f compose.ha-patroni.yml up -d
```

Failover behavior, verified live:

- Killing the primary promotes the sync standby within 30 seconds. Writes through HAProxy resume on the new primary with no lost rows.
- A restarted node rejoins as a streaming replica automatically via `pg_rewind`.
- Advisory locks live on the primary and do not carry over, so a database failover triggers a scheduler re-election automatically. Expect one missed schedule window at most.
- Check topology from any Patroni node: query `http://localhost:8008/cluster` for member roles and states.

Back up the control plane with `pg_dump` through HAProxy port 5432 on a schedule outside the cell. Restore into an empty database, then boot one API replica first so migrations settle before the others join. A managed Postgres service is an acceptable substitute for the Patroni tier; point `CODEDOCK_DATABASE_URL` at its primary endpoint.

## Images

First-party images publish to GHCR on every version tag with `latest` and multi-arch (`linux/amd64`, `linux/arm64`) builds:

- `ghcr.io/buildwithtechx/codedock` (self-hosted daemon, `Dockerfile`) consumed by the bootstrap installer.
- `ghcr.io/buildwithtechx/codedock-cloud` (`Dockerfile.cloud`) consumed by both HA compose cells via `CODEDOCK_VERSION` (default `latest`).

Pin `CODEDOCK_VERSION` to a tag for production cells. The Patroni image builds locally from `docker/patroni` (stock `postgres:16` plus pinned Patroni and driver); edge and etcd use pinned stock images.

## Operating the cell

- Scale API replicas with `docker compose -f compose.ha.yml up --scale api=3 -d`. The elector keeps a single scheduler regardless of replica count.
- Check leadership with `SELECT locktype, pid, granted FROM pg_locks WHERE locktype = 'advisory'` on the cell database. Exactly one granted session lock under the scheduler namespace is healthy.
- Readiness is `GET /healthz` on any replica. The edge Traefik only routes to healthy replicas.
- Rotate the vault key by re-encrypting every encrypted column first. There is no online rotation; plan a maintenance window, dump, rotate, restore.
- Never run two cells against one database. The elector assumes one cell per database; a second cell would split scheduled execution.
