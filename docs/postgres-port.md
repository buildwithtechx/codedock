# Postgres Pivot

Goal: run the Codedock control plane on Postgres everywhere. Self-hosted provisions an embedded Postgres container automatically; cloud and HA point `CODEDOCK_DATABASE_URL` at external Postgres. There is one dialect and one schema. All SQLite code and dependencies are deleted; there is no legacy support.

## Decisions

- Plain Postgres 16 through `pgx/v5/stdlib` over `database/sql`.
- `schema/*.sql` is canonical Postgres DDL.
- Repositories write `$n` placeholders directly through `sqlx.NewDb(db, "pgx")`. There is no `Rebind` helper and no driver parameter on the migration runner.
- `CODEDOCK_DATABASE_URL` unset means self-hosted: the daemon supervises a pinned `codedock-postgres` container (bind-mounted data under `dataDir/postgres`, 0600 password file, `unless-stopped`, health-gated boot). Set means use that Postgres.
- The supervisor never removes or recreates the container. Image upgrades are an explicit future operation, never automatic.
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

## Test strategy

- Every Postgres-backed test opens its own database (`CREATE DATABASE`, migrate, drop on cleanup) via `internal/testdb.Open` so parallel packages never share state.
- Env-gated suites: migration runner, supervisor lifecycle (needs Docker), daemon boot. All skip cleanly when `CODEDOCK_TEST_PG_URL` is unset.

## Remaining work

- Per-service images for hosted cells.
- Production Patroni deployment beyond the single-Postgres reference cell in `compose.ha.yml` (guidance: `docs/ha-runbook.md`).

Done: advisory-lock scheduler leadership with entry reconcile (`internal/engine/leadership`), multi-node review of backup/operations/autoscaling/attention/update loops, cross-node operation cancel, `self_hosted_config` table replacing `self-hosted.json`, `CODEDOCK_VAULT_KEY` for shared vault keys, reference HA compose stack and runbook.

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
