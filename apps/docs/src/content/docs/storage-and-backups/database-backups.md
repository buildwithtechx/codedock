---
title: Database Backups
description: Configure backup schedules, S3 destinations and recovery checks.
---

A backup configuration selects a database or service volume, a schedule and retention settings. Availability of database dumps depends on the installed engine template's backup command metadata. Do not assume every provisioned engine has a working backup and restore path.

## Configure a destination

An administrator can add an S3-compatible destination with a name, provider, endpoint, bucket, region, optional path prefix and access credentials. Use a pre-created bucket and verify access from the control plane. R2 and MinIO use this same destination workflow.

Destination creation is separate from provisioning an object-storage server. The credentials need permissions for the upload, download and deletion operations used by your backup lifecycle.

## Create a backup configuration

Select the target database or volume. Set `backupEnabled`, a cron `schedule`, `timezone`, timeout and retention settings. To upload off-server, enable `s3Enabled` and choose `s3DestinationId`. Keep local copies unless you intentionally enable `disableLocal` and have verified the remote recovery path.

| Field | Purpose |
| --- | --- |
| `schedule` | Cron expression, for example `0 2 * * *` |
| `timezone` | Time zone for the schedule, for example `UTC` |
| `timeout` | Execution timeout |
| `retentionDays` | Retention window |
| `maxBackups` | Maximum retained record count |
| `maxStorageGb` | Storage limit in the create configuration payload |

The scheduler uses the configured cron expression. Daily, weekly and monthly presets are not separate fixed retention guarantees.

## Test before scheduling

1. Create a configuration with the intended target and destination.
2. Trigger a manual backup.
3. Wait for a completed record and inspect its logs and file size.
4. Download the record and test restoration in a safe environment.
5. Enable the schedule and check that later runs complete.

A template without compatible dump metadata returns an unsupported-engine error. Volume archives are filesystem copies and do not guarantee application-level consistency for a live database.

## Recovery scope

Database backups and volume archives do not replace a backup of Codedock's SQLite state and vault key. Preserve the control-plane data directory and generated self-hosted settings securely; encrypted credentials depend on those keys.

See [restore and download](/storage-and-backups/restore-and-download/), [R2](/storage-and-backups/r2-storage/) and [MinIO](/storage-and-backups/minio-storage/).
