---
title: Restore and Download
description: Recover from a completed local or S3-compatible backup record.
---

Inspect completed backup records before restoring. Confirm the target, engine, file size and logs. Keep the credentials for any destination containing remote-only backups.

## Download

`GET /api/backups/:configId/records/:recordId/download` streams the record. The handler serves the recorded local file when present, or streams the stored S3 object when the local archive is absent. Remote-only records use their recorded destination and object key. The same project administrator authorization applies to both paths.

## Restore

`POST /api/backups/:recordId/restore` restores a specific record. The same route also accepts a configuration ID and selects a completed record; use an explicit record ID when the recovery point matters.

The restore operation requires administrative access to the target project and a database-backed record. Database restore and named-volume restore use separate workflows. Database restore can read a local file or fetch its S3 object with the configured destination credentials. Engine restore support depends on the installed template metadata and compatible database image. Test recovery before an incident.

1. Verify the database you are about to overwrite.
2. Make a fresh recovery copy where possible.
3. Pause application writes and plan for interruption.
4. Restore the selected completed record.
5. Inspect service logs, data and application connectivity.

Restore overwrites target data. A successful API response should be followed by application-level verification. Do not assume an archive from one engine or image version is portable to another.

## Named volume restore

The backup history opens a volume restore dialog for volume configurations. The target must be an existing named Docker volume mounted by its registered local service or database. Host bind mounts, unowned volumes and remote targets are rejected. Stop every container using the volume before validation and restoration.

- `GET /api/backup-records/:recordId/volume-restore` validates the target and reports its name and timeout.
- `POST /api/backup-records/:recordId/volume-restore` requires `{ "volumeName": "exact-target-name", "confirmOverwrite": true }` and restores the archive.
- `DELETE /api/backup-records/:recordId/volume-restore` interrupts the active restore after checking the same project administrator permissions.

Files present in the archive overwrite matching files; unrelated target files remain. The configured timeout is applied, capped at 24 hours. Disconnecting the request, reaching the timeout or choosing **Interrupt restore** stops and removes the helper container. Interruption cannot undo files already written. Keep consumers stopped, verify the contents and restore again if needed. Only one restore may run against a volume at a time.

## Delete and retain

Delete records only after checking retention needs. Keep an independent off-server copy for important recovery points and verify it periodically.

See [backup configuration](/storage-and-backups/database-backups/).
