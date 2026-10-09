import { BookOpen, ChevronLeft, ChevronRight, DatabaseBackup, ExternalLink } from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Skeleton } from '#/components/ui/skeleton';
import { DestinationLoadError } from './destination-display';
import { useDeleteRecord } from './hooks';
import type { BackupConfig, BackupRecord } from './interfaces';
import { RestoreWizard } from './restore-wizard';
import { SnapshotTableRow } from './snapshot-table-row';

const DOCS_URL = 'https://docs.codedock.run';

const pageSize = 10;

export function BackupDestinationHistory({
  records,
  configs,
  destinations,
  isLoading,
  loadError,
  onRetry,
  retrying,
  onChanged,
}: {
  records: BackupRecord[];
  configs: BackupConfig[];
  destinations: Map<string, string>;
  isLoading: boolean;
  loadError: string | null;
  onRetry: () => void;
  retrying: boolean;
  onChanged: () => void;
}) {
  const [page, setPage] = useState(0);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [restoreRecord, setRestoreRecord] = useState<BackupRecord | null>(null);
  const [pendingDelete, setPendingDelete] = useState<BackupRecord | null>(null);
  const deleteRecord = useDeleteRecord();

  const configById = useMemo(
    () => new Map(configs.map((config) => [config.id, config])),
    [configs]
  );
  const pageCount = Math.max(1, Math.ceil(records.length / pageSize));
  const safePage = Math.min(page, pageCount - 1);
  const visible = records.slice(safePage * pageSize, safePage * pageSize + pageSize);

  const sourceName = (record: BackupRecord): string => {
    const config = configById.get(record.backupConfigId);
    if (config) return config.name;
    if (record.databaseId) return `Database ${record.databaseId.slice(0, 8)}`;
    return 'System backup';
  };

  const destinationName = (record: BackupRecord): string => {
    const id = record.s3DestinationId || record.sftpDestinationId;
    if (!id) return 'Local only';
    return destinations.get(id) ?? 'Off-server';
  };

  const confirmDelete = async () => {
    if (!pendingDelete) return;
    try {
      await deleteRecord.mutateAsync({
        id: pendingDelete.backupConfigId,
        recordId: pendingDelete.id,
      });
      toast.success('Snapshot deleted');
      setPendingDelete(null);
      onChanged();
    } catch {
      toast.error('Failed to delete snapshot');
    }
  };

  return (
    <section
      aria-label="Recent snapshots"
      className="min-w-0 overflow-hidden rounded-2xl border border-border/50 bg-card"
    >
      <div className="flex items-start justify-between gap-4 px-5 py-4">
        <div>
          <h2 className="font-medium text-base">Recent snapshots</h2>
          <p className="mt-1 text-muted-foreground text-sm">
            Every automated and manual run across all destinations.
          </p>
        </div>
        {records.length > 0 && (
          <span className="mt-1 shrink-0 text-muted-foreground text-xs">
            {records.length} total
          </span>
        )}
      </div>

      {loadError && (
        <div className="px-5">
          <DestinationLoadError message={loadError} onRetry={onRetry} busy={retrying} />
        </div>
      )}

      {isLoading ? (
        <div aria-busy="true" className="space-y-3 px-5 pb-5">
          {[0, 1, 2].map((row) => (
            <Skeleton key={row} className="h-12 w-full" />
          ))}
        </div>
      ) : records.length === 0 ? (
        <div className="flex items-center gap-4 px-5 py-8">
          <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
            <DatabaseBackup className="h-5 w-5" />
          </div>
          <div>
            <p className="font-medium text-sm">No backup runs recorded yet</p>
            <p className="mt-1 text-muted-foreground text-sm">
              Automated and manual snapshots will appear here with restore actions.
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
        <div className="relative overflow-x-auto">
          <table className="w-full min-w-[760px] table-fixed text-left text-xs">
            <colgroup>
              <col className="w-[30%]" />
              <col />
              <col className="w-40" />
              <col className="w-20" />
              <col className="w-44" />
            </colgroup>
            <thead className="border-border/40 border-y bg-muted/20 text-muted-foreground">
              <tr>
                <th scope="col" className="px-5 py-3 font-medium">
                  Backup
                </th>
                <th scope="col" className="px-3 py-3 font-medium">
                  Contents
                </th>
                <th scope="col" className="px-3 py-3 font-medium">
                  Started
                </th>
                <th scope="col" className="px-3 py-3 text-right font-medium">
                  Size
                </th>
                <th scope="col" className="px-5 py-3">
                  <span className="sr-only">Actions</span>
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/40">
              {visible.map((record) => (
                <SnapshotTableRow
                  key={record.id}
                  record={record}
                  name={sourceName(record)}
                  destinationName={destinationName(record)}
                  expanded={expandedId === record.id}
                  onToggle={() => setExpandedId(expandedId === record.id ? null : record.id)}
                  onRestore={() => setRestoreRecord(record)}
                  onDelete={() => setPendingDelete(record)}
                  deletePending={deleteRecord.isPending}
                />
              ))}
            </tbody>
          </table>
        </div>
      )}

      {pageCount > 1 && (
        <nav
          aria-label="Snapshot pages"
          className="flex items-center justify-between gap-3 border-border/40 border-t px-5 py-3 text-muted-foreground text-xs"
        >
          <span>
            Page {safePage + 1} of {pageCount}
          </span>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              disabled={safePage === 0}
              onClick={() => {
                setExpandedId(null);
                setPage(safePage - 1);
              }}
            >
              <ChevronLeft className="h-4 w-4" />
              Previous
            </Button>
            <Button
              variant="ghost"
              size="sm"
              disabled={safePage >= pageCount - 1}
              onClick={() => {
                setExpandedId(null);
                setPage(safePage + 1);
              }}
            >
              Next
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        </nav>
      )}

      {restoreRecord && (
        <RestoreWizard
          record={restoreRecord}
          sourceName={sourceName(restoreRecord)}
          onClose={() => setRestoreRecord(null)}
          onDone={onChanged}
        />
      )}

      <Dialog
        open={pendingDelete !== null}
        onOpenChange={(open) => !open && setPendingDelete(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete this snapshot?</DialogTitle>
            <DialogDescription>
              This permanently removes the snapshot from its destination. This cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPendingDelete(null)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={confirmDelete} disabled={deleteRecord.isPending}>
              {deleteRecord.isPending ? 'Deleting…' : 'Delete snapshot'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </section>
  );
}
