# docker/

Container build contexts and runtime configs consumed by the compose cells.
Host-side installer scripts live in `bootstrap/` instead.

## patroni/

Production Postgres tier for `compose.ha-patroni.yml`, built locally (no published image):

- `Dockerfile` — stock `postgres:16` plus pinned Patroni 4.0.4 and `psycopg2-binary`
- `patroni.yml` — cluster template; per-node addresses and passwords arrive via `PATRONI_*` environment variables, which Patroni maps onto config keys
- `entrypoint.sh` — fixes data-dir ownership, then runs Patroni as `postgres`
- `post_bootstrap.py` — creates the `codedock` role (safely quoted password from `CODEDOCK_PG_PASSWORD`) and database after init
- `haproxy.cfg` — routes `:5432` to the primary and `:5433` to replicas via Patroni REST checks, plus a stats page on `:8404`

Topology and failover behavior: `docs/ha-runbook.md`.
