# Postgres Port Audit

Goal: run the Codedock control plane on plain Postgres for multi-node HA and hosted cells, while keeping SQLite as the self-hosted default. This document is the Phase 0 audit: exact gap inventory, porting strategy, and phase plan.

## Decisions

- Plain Postgres 16, connected through `pgx/v5/stdlib` over `database/sql`. No new dependencies: `pgx`, `pgx/stdlib` and `sqlx` are already imported.
- Single query source. Repositories keep writing `?` placeholders; `internal/repositories/dialect.go` rebinds to `$n` for Postgres via `sqlx.Rebind` at the boundary. SQLite behavior is unchanged.
- Per-dialect schema. SQLite keeps `schema/*.sql`; Postgres gets translated DDL in Phase 1. The migration runner gains a dialect-aware path with the same checksum discipline.
- Env-gated proof. `CODEDOCK_TEST_PG_URL` activates Postgres-backed repository tests; unset means skip. CI runs a `go-postgres` job with a real Postgres 16 service.

## Query-level gaps

| Gap | Sites | Postgres form |
| --- | --- | --- |
| `?` placeholders | 366 lines across all repositories | `Rebind` to `$n` |
| `substr(x, -n)` tail truncation | `attention.go`, `migration_store.go`, `operation.go` | `right(x, n)` |
| `char(10)` newline | `deployment_recovery.go` | `chr(10)` |
| `LIKE ... COLLATE NOCASE` | `deployment.go` search | `ILIKE` |
| `CURRENT_TIMESTAMP`, `\|\|`, `COALESCE`, `ON CONFLICT` | widespread | works as-is |

No `LastInsertId` in daemon code. No boolean literals in SQL. `LIMIT`/`OFFSET` accept bound parameters in both dialects.

## DDL gaps

Timestamps are mixed `DATETIME`/`TEXT` in SQLite; Postgres maps them to `timestamptz`. Go `time.Time` scans from both. Booleans map to `BOOLEAN`, blobs to `BYTEA`, text stays `TEXT`. IDs are strings (UUIDs), so no sequence handling. `INTEGER PRIMARY KEY` tables, if any carry auto-increment semantics, map to `GENERATED ALWAYS AS IDENTITY`.

## Single-node assumptions

- `internal/engine/cron` guards scheduling with an in-process mutex; backup, operations-reaper, attention-reconcile and update loops assume one daemon. Phase 4 moves these behind Postgres advisory-lock leadership.
- `CODEDOCK_DATA_DIR` holds the SQLite file, vault material, `self-hosted.json`, Traefik logs and backup staging. Phase 3 moves shared state into Postgres or object storage and documents what stays node-local.
- The migration runner executes whole-file SQL in one `Exec` inside a transaction; the Postgres path must preserve atomic apply plus checksum verification.

## Phases

- Phase 0 (this document): audit, `dialect.go`, dialect unit tests, Postgres probe test, `go-postgres` CI job.
- Phase 1: connection factory plus `CODEDOCK_DATABASE_URL`, translated Postgres schema, dialect-aware migration runner.
- Phase 2: repository port domain by domain, each with Postgres-backed tests; SQLite stays green; dual-dialect CI.
- Phase 3: file-local state migration.
- Phase 4: advisory-lock leadership and multi-node safety review.
- Phase 5: per-service images, compose stack, Patroni guidance, HA runbook.

## Running the Postgres tests

```sh
docker run -d --name codedock-pg -e POSTGRES_USER=codedock -e POSTGRES_PASSWORD=codedock -e POSTGRES_DB=codedock_test -p 5432:5432 postgres:16-alpine
CODEDOCK_TEST_PG_URL=postgres://codedock:codedock@localhost:5432/codedock_test?sslmode=disable go test ./internal/repositories/...
```

PowerShell:

```powershell
$env:CODEDOCK_TEST_PG_URL = "postgres://codedock:codedock@localhost:5432/codedock_test?sslmode=disable"
go test ./internal/repositories/...
```
