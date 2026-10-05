---
title: MinIO Destinations
description: Connect an existing MinIO bucket as an S3-compatible backup destination.
---

MinIO can be a backup destination. Registering a destination does not create a MinIO instance, bucket or application credentials.

## Prepare MinIO

Use an existing MinIO deployment or configure an image-based service with a persistent volume. Set explicit credentials, restrict console access and create the backup bucket. Keep the storage on a separate host if you need recovery from loss of the application server.

The S3 endpoint must be reachable from the control plane. A Docker service name is usable only when the control plane can resolve it and reach that network.

## Register the destination

In backup destination settings, provide:

| Setting | Value |
| --- | --- |
| Name | A recognizable label |
| Provider | MinIO / S3-compatible |
| Endpoint | Your reachable S3 API endpoint, including scheme |
| Bucket | The existing backup bucket |
| Region | The region configured for your instance |
| Access key and secret | Credentials authorized for backup objects |
| Path prefix | Optional organization prefix |

Verify the destination, select it in a backup configuration, enable S3 uploads and trigger a test. Download and restore a record before relying on it.

Application S3 credentials are configured separately in service variables. Codedock does not automatically inject this backup destination into all services.

See [database backups](/storage-and-backups/database-backups/).
