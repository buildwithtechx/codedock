---
title: Cloudflare R2 Destinations
description: Connect an R2 bucket through the S3-compatible backup destination workflow.
---

Create an R2 bucket and credentials with access to the required object operations. Use the S3 API endpoint for your account, not the public object URL.

## Destination settings

| Setting | Value |
| --- | --- |
| Endpoint | `https://ACCOUNT_ID.r2.cloudflarestorage.com` |
| Bucket | The existing R2 bucket name |
| Region | `auto` |
| Access key ID | R2 S3 access key ID |
| Secret access key | Matching R2 S3 secret |
| Path prefix | Optional object-key folder, applied to new uploads |

Use a dedicated bucket for backups. A path prefix is an object naming convention, not an access-isolation boundary; bucket policies and credentials control access.

New uploads use the configured prefix followed by the generated backup filename. Existing records keep their stored object keys for restore and retention.

Add the destination in backup settings and use the verify action. A successful check is useful, but a complete upload, download and restore test is the recovery check that matters.

## Enable uploads

Choose this destination in a backup configuration and enable S3 uploads. Decide whether to keep a local copy. Remote-only records require the destination credentials and bucket to remain available for download or restore.

Changing or deleting a destination can make existing remote records inaccessible. Preserve the bucket, object keys and authorized credentials as part of your recovery plan.

See [database backups](/storage-and-backups/database-backups/) and [restore and download](/storage-and-backups/restore-and-download/).
