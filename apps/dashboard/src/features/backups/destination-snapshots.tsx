import { BookOpen, DatabaseBackup, ExternalLink } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { cn } from '#/lib/utils';
import { formatBytes } from './destination-display';
import type { BackupRecord } from './interfaces';
import { RestoreWizard } from './restore-wizard';
import { SnapshotStatusChip } from './snapshot-table-row';

const DOCS_URL = 'https://docs.codedock.run';

export function DestinationSnapshots({
  snapshots,
  storedBytes,
  isLoading,
  destinationName,
  onChanged,
}: {
  snapshots: BackupRecord[];
  storedBytes: number;
  isLoading: boolean;
  destinationName: string;
  onChanged: () => void;
}) {
  const [restoreRecord, setRestoreRecord] = useState<BackupRecord | null>(null);

  return (
    <section aria-label="Snapshots" className="rounded-2xl border border-border/50 bg-card">
      <div className="px-5 py-4">
        <h2 className="font-medium text-base">Snapshots on this destination</h2>
        <p className="mt-1 text-muted-foreground text-sm">
          {snapshots.length} runs stored {formatBytes(storedBytes)} here.
        </p>
      </div>
      {isLoading ? (
        <div aria-busy="true" className="space-y-3 px-5 pb-5">
          {[0, 1, 2].map((row) => (
            <Skeleton key={row} className="h-12 w-full" />
          ))}
        </div>
      ) : snapshots.length === 0 ? (
        <div className="flex items-center gap-4 px-5 pb-6">
          <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
            <DatabaseBackup className="h-5 w-5" />
          </div>
          <div>
            <p className="font-medium text-sm">No snapshots stored yet</p>
            <p className="mt-1 text-muted-foreground text-sm">
              Runs from policies using this destination will appear here.
            </p>
            <Button asChild size="sm" variant="secondary" className="mt-3">
              <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
                <BookOpen className="size-3.5" />
                Docs
                <ExternalLink className="size-3 opacity-60" />
              </a>
            </Button>
          </div>
        </div>
      ) : (
        <ul className="divide-y divide-border/40 border-border/40 border-t">
          {snapshots.slice(0, 20).map((record) => (
            <li
              key={record.id}
              className={cn(
                'flex flex-wrap items-center gap-3 px-5 py-3 last:rounded-b-2xl hover:bg-muted/20'
              )}
            >
              <SnapshotStatusChip status={record.status} />
              <span className="min-w-0 flex-1 truncate text-muted-foreground text-xs">
                {record.startedAt ? new Date(record.startedAt).toLocaleString() : 'Unknown'} ·{' '}
                {record.fileSizeBytes ? formatBytes(record.fileSizeBytes) : '—'}
              </span>
              {record.status === 'completed' && (
                <Button
                  variant="outline"
                  size="sm"
                  className="h-7 text-xs"
                  onClick={() => setRestoreRecord(record)}
                >
                  Restore
                </Button>
              )}
            </li>
          ))}
        </ul>
      )}

      {restoreRecord && (
        <RestoreWizard
          record={restoreRecord}
          sourceName={destinationName}
          onClose={() => setRestoreRecord(null)}
          onDone={onChanged}
        />
      )}
    </section>
  );
}
