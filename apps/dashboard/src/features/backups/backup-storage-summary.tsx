import {
  AlertCircle,
  CheckCircle2,
  Clock,
  Database,
  HardDrive,
  Loader2,
  Server,
} from 'lucide-react';
import type { BackupConfig, BackupRecord, S3Destination } from './interfaces';

type BackupStorageSummaryProps = {
  records: BackupRecord[];
  destinations: S3Destination[];
  configs: BackupConfig[];
};

export function BackupStorageSummary({
  records,
  destinations,
  configs,
}: BackupStorageSummaryProps) {
  const totalStoredBytes = records
    .filter((r) => r.status === 'completed')
    .reduce((acc, r) => acc + (r.fileSizeBytes || 0), 0);

  const completedCount = records.filter((r) => r.status === 'completed').length;
  const runningCount = records.filter((r) => r.status === 'running').length;
  const failedCount = records.filter((r) => r.status === 'failed').length;

  const latestRecord = [...records]
    .filter((r) => r.completedAt)
    .sort((a, b) => new Date(b.completedAt).getTime() - new Date(a.completedAt).getTime())[0];

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / k ** i).toFixed(1)} ${sizes[i]}`;
  };

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-sm">
      <div className="mb-4 flex items-center gap-2">
        <HardDrive className="h-4 w-4 text-primary" />
        <h2 className="font-semibold text-foreground text-sm">Storage summary</h2>
      </div>

      <dl className="space-y-3">
        <div className="flex items-center justify-between text-sm">
          <dt className="flex items-center gap-2 text-muted-foreground">
            <Database className="h-4 w-4 shrink-0" />
            <span>Total stored data</span>
          </dt>
          <dd className="font-semibold text-foreground">{formatBytes(totalStoredBytes)}</dd>
        </div>

        <div className="flex items-center justify-between text-sm">
          <dt className="flex items-center gap-2 text-muted-foreground">
            <CheckCircle2 className="h-4 w-4 shrink-0 text-emerald-500" />
            <span>Completed snapshots</span>
          </dt>
          <dd className="font-medium text-foreground">{completedCount}</dd>
        </div>

        {runningCount > 0 && (
          <div className="flex items-center justify-between text-sm">
            <dt className="flex items-center gap-2 text-muted-foreground">
              <Loader2 className="h-4 w-4 shrink-0 animate-spin text-sky-500" />
              <span>Active in progress</span>
            </dt>
            <dd className="font-medium text-sky-500">{runningCount}</dd>
          </div>
        )}

        {failedCount > 0 && (
          <div className="flex items-center justify-between text-sm">
            <dt className="flex items-center gap-2 text-muted-foreground">
              <AlertCircle className="h-4 w-4 shrink-0 text-rose-500" />
              <span>Failed snapshots</span>
            </dt>
            <dd className="font-medium text-rose-500">{failedCount}</dd>
          </div>
        )}

        <div className="flex items-center justify-between text-sm">
          <dt className="flex items-center gap-2 text-muted-foreground">
            <Server className="h-4 w-4 shrink-0" />
            <span>Connected destinations</span>
          </dt>
          <dd className="font-medium text-foreground">{destinations.length}</dd>
        </div>

        <div className="flex items-center justify-between text-sm">
          <dt className="flex items-center gap-2 text-muted-foreground">
            <Clock className="h-4 w-4 shrink-0" />
            <span>Active policies</span>
          </dt>
          <dd className="font-medium text-foreground">
            {configs.filter((c) => c.backupEnabled).length}
          </dd>
        </div>

        <div className="flex items-center justify-between border-border/40 border-t pt-3 text-xs">
          <dt className="text-muted-foreground">Last successful backup</dt>
          <dd className="font-medium text-foreground">
            {latestRecord?.completedAt
              ? new Date(latestRecord.completedAt).toLocaleString()
              : 'None yet'}
          </dd>
        </div>
      </dl>
    </section>
  );
}
