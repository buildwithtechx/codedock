import {
  ArchiveRestore,
  CheckCircle2,
  ChevronDown,
  Download,
  Loader2,
  Trash2,
  XCircle,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import { getApiBaseUrl } from '#/lib/api-client';
import { cn } from '#/lib/utils';
import { BackupRecordProtection } from './backup-record-protection';
import { formatBytes, formatDuration } from './destination-display';
import type { BackupRecord } from './interfaces';

export function SnapshotStatusChip({ status }: { status: BackupRecord['status'] }) {
  const completed = status === 'completed';
  const failed = status === 'failed';
  return (
    <span
      className={cn(
        'inline-flex shrink-0 items-center gap-1 rounded-full px-2 py-0.5 font-medium text-[11px] capitalize',
        completed && 'bg-emerald-500/10 text-emerald-500',
        failed && 'bg-destructive/10 text-destructive',
        !completed && !failed && 'bg-sky-500/10 text-sky-500'
      )}
    >
      {completed ? (
        <CheckCircle2 className="h-3 w-3" />
      ) : failed ? (
        <XCircle className="h-3 w-3" />
      ) : (
        <Loader2 className="h-3 w-3 animate-spin" />
      )}
      {status}
    </span>
  );
}

export function SnapshotTableRow({
  record,
  name,
  destinationName,
  expanded,
  onToggle,
  onRestore,
  onDelete,
  deletePending,
}: {
  record: BackupRecord;
  name: string;
  destinationName: string;
  expanded: boolean;
  onToggle: () => void;
  onRestore: () => void;
  onDelete: () => void;
  deletePending: boolean;
}) {
  const completed = record.status === 'completed';
  const downloadHref = record.s3Url
    ? record.s3Url
    : record.filePath
      ? `${getApiBaseUrl()}/backups/${record.backupConfigId}/records/${record.id}/download`
      : null;
  return (
    <>
      <tr className="hover:bg-muted/20">
        <td className="px-5 py-3">
          <div className="flex items-center gap-2">
            <span className="truncate font-medium text-foreground text-sm" title={name}>
              {name}
            </span>
            <SnapshotStatusChip status={record.status} />
          </div>
          <p className="mt-0.5 truncate text-muted-foreground" title={destinationName}>
            {destinationName}
          </p>
        </td>
        <td className="px-3 py-3">
          <p className="truncate text-foreground" title={record.filePath || '—'}>
            {record.filePath ? record.filePath.split('/').pop() : '—'}
          </p>
          {record.verifiedAt && <p className="mt-0.5 text-muted-foreground">Verified snapshot</p>}
        </td>
        <td className="whitespace-nowrap px-3 py-3 text-muted-foreground">
          <time
            dateTime={record.startedAt}
            className="block truncate"
            title={record.startedAt ? new Date(record.startedAt).toLocaleString() : undefined}
          >
            {record.startedAt
              ? new Date(record.startedAt).toLocaleString(undefined, {
                  dateStyle: 'medium',
                  timeStyle: 'short',
                })
              : '—'}
          </time>
          <p className="mt-0.5 truncate">{formatDuration(record.startedAt, record.completedAt)}</p>
        </td>
        <td className="whitespace-nowrap px-3 py-3 text-right text-muted-foreground tabular-nums">
          {record.fileSizeBytes ? formatBytes(record.fileSizeBytes) : '—'}
        </td>
        <td className="px-5 py-3">
          <div className="flex items-center justify-end gap-1">
            {completed && (
              <Button variant="outline" size="sm" className="h-7 text-xs" onClick={onRestore}>
                <ArchiveRestore className="h-3.5 w-3.5" />
                Restore
              </Button>
            )}
            {downloadHref && (
              <Button variant="ghost" size="icon" className="h-7 w-7 shrink-0" asChild>
                <a href={downloadHref} target="_blank" rel="noreferrer" title="Download archive">
                  <Download className="h-3.5 w-3.5" />
                </a>
              </Button>
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 shrink-0"
              title="View details"
              aria-expanded={expanded}
              onClick={onToggle}
            >
              <ChevronDown
                className={cn('h-4 w-4 transition-transform', expanded && 'rotate-180')}
              />
            </Button>
          </div>
        </td>
      </tr>
      {expanded && (
        <tr className="bg-muted/20">
          <td colSpan={5} className="px-5 py-4">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <p className="font-mono text-muted-foreground text-xs">
                {record.s3Url || record.sftpUrl || record.filePath || record.id}
              </p>
              <div className="flex items-center gap-1.5">
                {completed && (
                  <BackupRecordProtection recordId={record.id} until={record.protectedUntil} />
                )}
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 text-destructive text-xs hover:bg-destructive/10 hover:text-destructive"
                  disabled={deletePending || (record.protectedUntil ?? 0) * 1000 > Date.now()}
                  onClick={onDelete}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                  Delete
                </Button>
              </div>
            </div>
            {record.logs ? (
              <pre className="mt-3 max-h-48 overflow-auto whitespace-pre-wrap rounded-lg bg-background/60 p-3 font-mono text-[11px] text-muted-foreground">
                {record.logs}
              </pre>
            ) : (
              <p className="mt-3 text-muted-foreground text-xs">No logs captured for this run.</p>
            )}
          </td>
        </tr>
      )}
    </>
  );
}
