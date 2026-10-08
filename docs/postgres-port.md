# Postgres Pivot

Goal: run the Codedock control plane on Postgres everywhere. Self-hosted provisions an embedded Postgres container automatically; cloud and HA point `CODEDOCK_DATABASE_URL` at external Postgres. There is one dialect and one schema. SQLite survives only as transitional scaffolding for not-yet-ported repository tests, plus the one-shot legacy importer.

## Decisions

- Plain Postgres 16 through `pgx/v5/stdlib` over `database/sql`.
- `schema/*.sql` is canonical Postgres DDL. `schema/sqlite/*.sql` is transitional and gets deleted when the last repository test ports.
- Repositories write `$n` placeholders directly. The `Rebind`/`DriverName` helpers and the driver parameter on the migration runner are transitional scaffolding and get deleted with the SQLite schema.
- `CODEDOCK_DATABASE_URL` unset means self-hosted: the daemon supervises a pinned `codedock-postgres` container (bind-mounted data under `dataDir/postgres`, 0600 password file, `unless-stopped`, health-gated boot). Set means use that Postgres.
- The supervisor never removes or recreates the container. Image upgrades are an explicit future operation, never automatic.
- A legacy `dataDir/codedock.db` is imported once into a fresh migrated Postgres, then renamed to `codedock.db.imported`. Import refuses a non-empty target. Nothing is ever deleted.
- `CODEDOCK_TEST_PG_URL` activates Postgres-backed tests; unset means skip. CI runs a `go-postgres` job with a real Postgres 16 service covering repositories, systemdb and daemon commands.

## Query porting pattern

Each repository port follows the settings-repo example (`server_settings.go` was first):

- `sqlx.NewDb(db, "sqlite")` becomes `sqlx.NewDb(db, "pgx")`.
- `?` placeholders become `$n`.
- `substr(x, -n)` becomes `right(x, n)` (`attention.go`, `migration_store.go`, `operation.go`).
- `char(10)` becomes `chr(10)` (`deployment_recovery.go`).
- `LIKE ... COLLATE NOCASE` becomes `ILIKE` (`deployment.go` search).
- `INSERT OR IGNORE` becomes `INSERT ... ON CONFLICT DO NOTHING`; `INSERT OR REPLACE` becomes `ON CONFLICT DO UPDATE`.
- `CURRENT_TIMESTAMP`, `||`, `COALESCE`, `ON CONFLICT(id) DO UPDATE`, `LIMIT`/`OFFSET` with bound parameters work as-is.
- timestamptz columns scan into Go `bool`/`string`/`time.Time` model fields; pgx parses RFC 3339 and `YYYY-MM-DD HH:MM:SS` on write.
- The `deployments.trigger` column is quoted (`"trigger"` is reserved in Postgres).

No `LastInsertId` in daemon code. IDs are strings, so no sequence handling. The user-database data browser (`handlers/databases`, `services/databases`) queries customer databases and is untouched by this pivot.

## DDL translation record

`schema/*.sql` was translated from the SQLite originals by six mechanical rules, verified by reverse-translation identity: `DATETIME` to `TIMESTAMPTZ`, `REAL` to `DOUBLE PRECISION`, `BOOLEAN DEFAULT 0/1` to `DEFAULT FALSE/TRUE`, `trigger` quoted, `is_active=1` to `is_active=TRUE`, plus three structural fixes Postgres requires: forward foreign keys in `001` moved to trailing `ALTER TABLE ... ADD FOREIGN KEY` statements, and `CREATE VIEW IF NOT EXISTS services` became `CREATE OR REPLACE VIEW services`.

## Test strategy during the port

- Ported tests open their own database (`CREATE DATABASE`, migrate, drop on cleanup) so parallel packages never share state. See `setupTestDatabaseURL` in `cmd/codedockd/commands/setup_test.go`.
- Unported repository tests keep `:memory:` SQLite plus `RunMigrationsDialect(db, DriverSQLite)` and stay green.
- Env-gated suites: migration runner, SQLite importer, supervisor lifecycle (needs Docker), daemon boot.

## Remaining work

- Port repositories domain by domain with Postgres-backed tests (58 files). Suggested order: auth, users, projects, app services, deployments, environments, databases, backups, then the rest.
- Delete scaffolding when the last port lands: `schema/sqlite/`, `dialect.go`, the driver parameter, `modernc.org/sqlite` from non-test code (the importer keeps it until then).
- Move single-node assumptions behind Postgres: cron mutex to advisory-lock leadership, backup/operations-reaper/attention-reconcile/update loops multi-node review.
- Migrate file-local shared state (`self-hosted.json`, backup staging) into Postgres or object storage; document what stays node-local.
- Per-service images, compose stack, Patroni guidance, HA runbook for hosted cells.

## Running the Postgres tests

```sh
docker run -d --name codedock-pg -e POSTGRES_USER=codedock -e POSTGRES_PASSWORD=codedock -e POSTGRES_DB=codedock_test -p 5432:5432 postgres:16-alpine
CODEDOCK_TEST_PG_URL=postgres://codedock:codedock@localhost:5432/codedock_test?sslmode=disable go test ./internal/repositories/... ./internal/engine/systemdb/... ./cmd/...
```

PowerShell:

```powershell
$env:CODEDOCK_TEST_PG_URL = "postgres://codedock:codedock@localhost:5432/codedock_test?sslmode=disable"
go test ./internal/repositories/... ./internal/engine/systemdb/... ./cmd/...
```
