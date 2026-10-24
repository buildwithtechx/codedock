---
title: SQL Studio
description: Execute queries against compatible database engines.
---

Open the database detail page and select the query editor. PostgreSQL, MySQL and MariaDB use the SQL query path. Redis accepts commands through the query endpoint; MongoDB and ClickHouse are unsupported by this implementation.

## Execute a query

Start with a bounded read:

```sql
SELECT id, display_name FROM accounts ORDER BY id LIMIT 25;
```

`POST /api/databases/:id/query` accepts `{ "query": "..." }` and returns query information inside the standard response envelope. Relational reads return columns and rows; command results vary, and Redis returns its command value in `result`.

```sh
curl "$CODEDOCK_URL/api/databases/$DATABASE_ID/query"   -H "Authorization: Bearer $CODEDOCK_TOKEN"   -H "Content-Type: application/json"   -d '{"query":"SELECT 1 AS healthy"}'
```

Queries run with the configured database credentials. This endpoint is not a read-only sandbox; writes and schema changes can execute. Back up data before destructive changes, and use migration tooling for changes that must be repeatable.

## Reachability

The query service currently connects to `localhost` at the database's configured port. Remote-only databases require an external client or separately configured reachability. A running container alone does not establish that SQL Studio can reach it.

## Redis commands

The Redis query path splits the command on whitespace. It is suitable for simple commands such as `PING`, but does not implement a complete Redis CLI quoting parser. Use a dedicated Redis client for binary values or complex arguments.

See [data browser](/databases/data-browser/) and [backups](/storage-and-backups/database-backups/).
