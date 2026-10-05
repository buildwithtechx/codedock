---
title: Restore and Download
description: Recover from a completed local or S3-compatible backup record.
---

Inspect completed backup records before restoring. Confirm the target, engine, file size and logs. Keep the credentials for any destination containing remote-only backups.

## Download

`GET /api/backups/:configId/records/:recordId/download` streams the record. The server uses an available local file or downloads the stored object through its S3 destination. Missing files, inaccessible buckets or changed credentials prevent recovery.

## Restore

`POST /api/backups/:recordId/restore` restores a specific record. The same route also accepts a configuration ID and selects a completed record; use an explicit record ID when the recovery point matters.

The restore operation requires administrative access to the target project. Engine restore support depends on the installed template metadata and compatible database image. Test recovery before an incident.

1. Verify the database or volume you are about to overwrite.
2. Make a fresh recovery copy where possible.
3. Pause application writes and plan for interruption.
4. Restore the selected completed record.
5. Inspect service logs, data and application connectivity.

Restore overwrites target data. A successful API response should be followed by application-level verification. Do not assume an archive from one engine or image version is portable to another.

## Delete and retain

Delete records only after checking retention needs. Keep an independent off-server copy for important recovery points and verify it periodically.

See [backup configuration](/storage-and-backups/database-backups/).
