---
title: Data Browser
description: Inspect relational tables and make small, explicit data changes.
---

Open a running database's detail page and its data browser. The current backend implements schema and table access for PostgreSQL, MySQL and MariaDB. PostgreSQL schema discovery lists tables in the `public` schema.

MongoDB collection browsing, a Redis key editor and ClickHouse table browsing are not implemented by these generic endpoints. Use an engine-specific client where the browser is unsupported.

## Read records

Choose a table and use pagination to inspect a bounded selection of rows. The API accepts `limit` and `offset`; omitted or non-positive limits default to 100. Generic filter expressions are not accepted by this endpoint.

```sh
curl "$CODEDOCK_URL/api/databases/$DATABASE_ID/data/accounts?limit=25&offset=0"   -H "Authorization: Bearer $CODEDOCK_TOKEN"
```

## Edit records

Inserts accept a JSON object of column values. Updates accept `keys` to identify the existing row and `data` for changed values. Deletes accept a key-value object. Use stable unique keys and verify the affected record before making another change. Do not assume schema discovery reliably identifies primary keys.

```json
{"keys":{"id":42},"data":{"display_name":"Updated name"}}
```

Writes use the configured database credentials and can affect live application data. Make a recovery copy before migration or bulk changes. Use your migration tooling for repeatable schema changes.

## Connection requirements

The current browser connects from the control plane to `localhost` using the database's configured port. A database running only on a remote worker or without a reachable local port cannot be queried through this path. Use a client with the appropriate secure network connection instead.

Check database state, port, credentials and host reachability if schema loading fails. See [SQL Studio](/databases/sql-studio/) and [database API](/reference/api-databases/).
