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

## Deployment checklist

1. Pin an image version and read its configuration requirements.
2. Generate the secrets the application requires.
3. Configure listening ports and persistent volumes.
4. Check logs and connectivity from the consuming application.
5. Back up data and test recovery before depending on the deployment.
