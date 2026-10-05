---
title: Database Provisioning
description: Create supported database resources and verify runtime and operational capabilities.
---

Database resources belong to a project and can be associated with an environment. Their containers require Docker and persistent storage on the selected host.

## Engine identifiers

The current database model defines these engine IDs:

| Engine | Identifier | Conventional port |
| --- | --- | --- |
| PostgreSQL | `postgres` | 5432 |
| MySQL | `mysql` | 3306 |
| MariaDB | `mariadb` | 3306 |
| Redis | `redis` | 6379 |
| MongoDB | `mongodb` | 27017 |
| ClickHouse | `clickhouse` | 9000 |

Select a version compatible with the deployed image. An engine identifier is not a guarantee that every runtime template, browser or backup feature is available on your installation. Verify startup and the operational tools you require.

Other software can be deployed as application containers with explicit configuration; it is not automatically a built-in managed database engine.

## Create a database

1. Choose the intended project and environment.
2. Create a database resource with a name, engine and version.
3. Set its database name and required credentials.
4. Review port and persistent storage settings.
5. Start the resource and check its status and logs.

The API uses `POST /api/databases`. See [database API fields](/reference/api-databases/).

## Connect an application

Copy the connection values from the database resource into the application's variables. Use a hostname and port reachable from that application, and keep credentials in secret settings. Do not assume creation injects the correct connection string into every service.

A localhost address refers to the current container or process host. A Docker-internal hostname requires shared network access. Remote worker placement needs explicit cross-host connectivity.

## Operational support

The relational browser and SQL query paths currently support PostgreSQL, MySQL and MariaDB over a control-plane-local connection. Redis has a command query path. MongoDB and ClickHouse require appropriate external clients for unsupported browsing and querying operations.

Backups are configured separately and depend on compatible engine template metadata. Daily backups are not enabled automatically just by creating a database. Configure and test the workflow you need.

See [data browser](/databases/data-browser/), [SQL Studio](/databases/sql-studio/) and [backup configuration](/storage-and-backups/database-backups/).
