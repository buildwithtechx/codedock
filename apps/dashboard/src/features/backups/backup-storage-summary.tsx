import {
  AlertCircle,
  CheckCircle2,
  Clock,
  Database,
  HardDrive,
  Loader2,
  Server,
} from 'lucide-react';
import { formatBytes } from './destination-display';
import type { BackupConfig, BackupRecord, S3Destination } from './interfaces';

type BackupStorageSummaryProps = {
  records: BackupRecord[];
  destinations: S3Destination[];
  configs: BackupConfig[];
  sftpCount?: number;
};

export function BackupStorageSummary({
  records,
  destinations,
  configs,
  sftpCount = 0,
}: BackupStorageSummaryProps) {
  const completed = records.filter((record) => record.status === 'completed');
  const totalStoredBytes = completed.reduce(
    (total, record) => total + (record.fileSizeBytes || 0),
    0
  );
  const runningCount = records.filter((record) => record.status === 'running').length;
  const failedCount = records.filter((record) => record.status === 'failed').length;
  const latest = [...completed]
    .filter((record) => record.completedAt)
    .sort((a, b) => new Date(b.completedAt).getTime() - new Date(a.completedAt).getTime())[0];

  return (
    <section
      aria-label="Storage summary"
      className="rounded-2xl border border-border/50 bg-card p-5"
    >
      <div className="mb-4 flex items-center gap-2">
        <HardDrive className="h-4 w-4 text-muted-foreground" />
        <h2 className="font-semibold text-foreground text-sm">Storage summary</h2>
      </div>
      <dl className="space-y-3">
        <div className="flex items-center justify-between gap-3 text-sm">
          <dt className="flex min-w-0 items-center gap-2 text-muted-foreground">
            <Database className="h-4 w-4 shrink-0" />
            Stored data
          </dt>
          <dd className="text-end font-medium">{formatBytes(totalStoredBytes)}</dd>
        </div>
        <div className="flex items-center justify-between gap-3 text-sm">
          <dt className="flex min-w-0 items-center gap-2 text-muted-foreground">
            <CheckCircle2 className="h-4 w-4 shrink-0" />
            Snapshots
          </dt>
          <dd className="text-end font-medium">{completed.length}</dd>
        </div>
        {runningCount > 0 && (
          <div className="flex items-center justify-between gap-3 text-sm">
            <dt className="flex min-w-0 items-center gap-2 text-muted-foreground">
              <Loader2 className="h-4 w-4 shrink-0 animate-spin" />
              Active
            </dt>
            <dd className="text-end font-medium text-sky-500">{runningCount}</dd>
          </div>
        )}
        {failedCount > 0 && (
          <div className="flex items-center justify-between gap-3 text-sm">
            <dt className="flex min-w-0 items-center gap-2 text-muted-foreground">
              <AlertCircle className="h-4 w-4 shrink-0" />
              Failed
            </dt>
            <dd className="text-end font-medium text-destructive">{failedCount}</dd>
          </div>
        )}
        <div className="flex items-center justify-between gap-3 text-sm">
          <dt className="flex min-w-0 items-center gap-2 text-muted-foreground">
            <Clock className="h-4 w-4 shrink-0" />
            Last snapshot
          </dt>
          <dd className="text-end font-medium">
            {latest?.completedAt
              ? new Date(latest.completedAt).toLocaleString(undefined, {
                  dateStyle: 'medium',
                  timeStyle: 'short',
                })
              : '—'}
          </dd>
        </div>
        <div className="flex items-center justify-between gap-3 text-sm">
          <dt className="flex min-w-0 items-center gap-2 text-muted-foreground">
            <Server className="h-4 w-4 shrink-0" />
            Destinations
          </dt>
          <dd className="text-end font-medium">{destinations.length + sftpCount}</dd>
        </div>
        <div className="flex items-center justify-between gap-3 text-sm">
          <dt className="flex min-w-0 items-center gap-2 text-muted-foreground">
            <CheckCircle2 className="h-4 w-4 shrink-0" />
            Active policies
          </dt>
          <dd className="text-end font-medium">
            {configs.filter((config) => config.backupEnabled).length}
          </dd>
        </div>
      </dl>
    </section>
  );
}
