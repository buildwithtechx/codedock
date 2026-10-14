---
title: Templates and Recipes
description: Choose a supported database workflow or configure a Docker or Compose recipe.
---

The [marketing catalogue](https://codedock.run/templates) links to these setup recipes. A catalogue entry does not mean its service is available as an automatic one-click installer.

The dashboard's one-click list is returned by `GET /api/one-click` from templates with the required metadata. Use the list returned by your instance as the source of availability. Do not assume the recipes below appear there.

## Database workflow

Create a database in a project, select an engine and version, configure resources and credentials, then start it. Copy connection values into the application's variables. Runtime support, browsing and backup support differ by engine; test the features you depend on.

### Postgres

Select `postgres`. Use the relational browser or SQL Studio for inspection. Configure a persistent data location and test a compatible backup and restore before relying on the service.

### MySQL

Select `mysql`. Configure a database name and credentials. Use the relational browser or SQL Studio after the service is reachable.

### Redis

Select `redis`. Set a password and choose persistence requirements appropriate to the application. The query endpoint accepts Redis commands; the generic relational table browser does not apply.

### MongoDB

Select `mongodb`. Use a MongoDB client for collection inspection. The current relational browser and SQL query endpoint do not support MongoDB collections.

## Container recipes

Use a versioned vendor image, configure required secrets and port settings, attach persistent volumes, and restrict public access. Follow the vendor's image documentation for the version you select. Review upgrades and backup requirements separately from application deployment.

### Supabase

Supabase is a multi-service stack. Obtain its Compose configuration and required secrets from its official distribution. The Codedock Compose importer does not preserve environment, volumes, dependencies or overrides; do not use it as a drop-in Supabase installation. Run the complete Compose stack directly, or configure and verify every resource explicitly.

### MinIO

Create a container service with a persistent data volume and explicit root credentials. Separate the S3 API from the administrative console. For backup uploads, register its S3 endpoint and bucket as a destination; this does not create a MinIO server or automatically inject credentials into applications.

See [MinIO destinations](/storage-and-backups/minio-storage/).

### N8N

Configure the n8n image, its database connection, encryption key, public URL and persistent data. Keep the encryption key available for recovery; recreating the container is not a substitute for a complete application backup.

### Grafana

Configure a Grafana image and persistent data. Supply authentication and data-source settings explicitly. Deploying Grafana does not automatically connect it to Codedock's metrics API or create a Prometheus instance.

### Ollama

Configure model storage and resource capacity. GPU drivers and runtime configuration are host prerequisites. See [AI workloads](/deployments/ai-workloads/).

### Qdrant

Configure persistent vector storage and authentication, then restrict access to the required application network. Use Qdrant's version-specific backup workflow.

## Application stacks

One-click Compose stacks from the verified catalogue. Each entry deploys its services, volumes and secrets together; review the stack's secrets and persistence before deploying.

### Convex

Self-hosted Convex reactive backend and dashboard. Configure the instance secret and a persistent data volume before deploying.

### Neon

Self-hosted Neon serverless Postgres with database branching, a web console and per-branch connection strings. Set the server secret and admin credentials.

### Cal.com

Self-hosted scheduling platform with calendar integrations and booking pages.

### Code-server

VS Code in the browser with a persistent workspace volume.

### Convex

Self-hosted Convex reactive backend and dashboard with a persistent data volume.

### Directus

Headless CMS and data platform over any SQL database.

### FreshRSS

Self-hosted RSS and Atom feed reader with a persistent database.

### Ghost

Professional publishing platform with themes and member subscriptions.

### Gitea

Self-hosted Git service with issues, pull requests and package registries.

### IT Tools

Handy developer utilities in a single page with no persistence requirements.

### Codedock-mail

Self-hosted mail control plane: maild API, Postgres metadata store and SnappyMail webmail. Point it at a Stalwart backend.

### Meilisearch

Fast instant-search engine with a REST API and persistent index storage.

### Metabase

Business intelligence dashboards over your databases with persistent application data.

### Nextcloud

Self-hosted file storage and collaboration suite with a database backend.

### NocoDB

Spreadsheet-style database interface over Postgres with persistent metadata.

### Plausible

Privacy-friendly web analytics with a ClickHouse backend.

### PocketBase

Backend with database, auth and file storage in a single binary with persistent data.

### Stirling PDF

PDF tools for merging, splitting, converting and OCR with no persistence requirements.

### Umami

Simple privacy-focused web analytics with a Postgres backend.

### Uptimekuma

Self-hosted uptime monitoring with persistent status data.

### Vaultwarden

Lightweight Bitwarden-compatible password manager. Back up its data volume; losing it loses every vault.

### WordPress

CMS and blogging platform with a MySQL backend and persistent uploads.

## Deployment checklist

1. Pin an image version and read its configuration requirements.
2. Generate the secrets the application requires.
3. Configure listening ports and persistent volumes.
4. Check logs and connectivity from the consuming application.
5. Back up data and test recovery before depending on the deployment.
