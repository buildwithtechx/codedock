---
title: Databases API
description: Current provisioning, query, data and backup routes.
---

## Create a database

`POST /api/databases` accepts `projectId`, optional `environmentId`, `name`, `engine`, `version`, `databaseName`, `username`, `password` and optional port and volume settings.

The model constants include `postgres`, `mysql`, `mariadb`, `redis`, `mongodb` and `clickhouse`. The provisioning UI also offers template IDs listed in [database provisioning](/databases/provisioning/). Provisioning support does not imply browser, query or backup support for every engine. Supply `port: 3306` explicitly for MariaDB and `port: 9000` for ClickHouse; the current omitted-port fallback is 5432 for engines outside the specific PostgreSQL, MySQL, Redis and MongoDB cases.

```json
{"projectId":"PROJECT_ID","name":"primary","engine":"postgres","version":"16","databaseName":"app","username":"app","password":"YOUR_GENERATED_PASSWORD"}
```

See the [database DTO](https://github.com/buildwithtechx/codedock/blob/main/pkg/types/database.go) for all creation fields.

## Manage a database

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/databases` | List databases |
| GET | `/api/databases/:id` | Inspect |
| PUT | `/api/databases/:id` | Update supported settings |
| DELETE | `/api/databases/:id` | Delete |
| POST | `/api/databases/:id/start` | Start |
| POST | `/api/databases/:id/stop` | Stop |
| POST | `/api/databases/:id/restart` | Restart |
| POST | `/api/databases/:id/credentials/reveal` | Reveal connection credentials |
| POST | `/api/databases/:id/import` | Import data |

These routes require the `database:manage` scope and project authorization. Do not record revealed credentials in logs.

## Queries and data

`POST /api/databases/:id/query` accepts `{ "query": "SELECT 1" }`. SQL runs for PostgreSQL, MySQL and MariaDB; Redis uses its command path. MongoDB and ClickHouse are unsupported by this query service. Connections currently use the control plane's localhost and configured database port.

| Method | Path | Payload / parameters |
| --- | --- | --- |
| GET | `/api/databases/:id/schemas` | Relational schema list |
| GET | `/api/databases/:id/data/:table` | `limit` and `offset` query parameters |
| POST | `/api/databases/:id/data/:table` | Column-value JSON object |
| PUT | `/api/databases/:id/data/:table` | `{ "keys": {...}, "data": {...} }` |
| DELETE | `/api/databases/:id/data/:table` | Key-value JSON object |

The relational browser supports PostgreSQL, MySQL and MariaDB. See [data browser](/databases/data-browser/) and [SQL Studio](/databases/sql-studio/).

## Backups

List database backup records with `GET /api/databases/:id/backups`; trigger with `POST /api/databases/:id/backups`. Backup configuration, destination and record operations use the separate `/backups` and `/s3-destinations` routes. Engine support depends on compatible template metadata. See [backup configuration](/storage-and-backups/database-backups/).
