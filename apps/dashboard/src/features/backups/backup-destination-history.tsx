import {
  ArchiveRestore,
  CheckCircle2,
  Download,
  History,
  Loader2,
  Trash2,
  XCircle,
} from 'lucide-react';
import { useState } from 'react';
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { getApiBaseUrl } from '#/lib/api-client';
import type { BackupRecord } from './interfaces';

type BackupDestinationHistoryProps = {
  records: BackupRecord[];
  isLoading: boolean;
  onRestore: (recordId: string) => Promise<'confirmation' | 'completed'>;
  onDeleteRecord: (configId: string, recordId: string) => Promise<void>;
  restorePending: boolean;
  deletePending: boolean;
};

export function BackupDestinationHistory({
  records,
  isLoading,
  onRestore,
  onDeleteRecord,
  restorePending,
  deletePending,
}: BackupDestinationHistoryProps) {
  const [selectedRecordForRestore, setSelectedRecordForRestore] = useState<BackupRecord | null>(
    null
  );

  const formatBytes = (bytes: number): string => {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / k ** i).toFixed(1)} ${sizes[i]}`;
  };

  const calculateDuration = (start?: string, end?: string): string => {
    if (!start || !end) return '—';
    const diffMs = new Date(end).getTime() - new Date(start).getTime();
    if (diffMs <= 0) return '< 1s';
    const seconds = Math.floor(diffMs / 1000);
    if (seconds < 60) return `${seconds}s`;
    const mins = Math.floor(seconds / 60);
    const remainingSecs = seconds % 60;
    return `${mins}m ${remainingSecs}s`;
  };

  const handleConfirmRestore = async () => {
    if (!selectedRecordForRestore) return;
    try {
      const result = await onRestore(selectedRecordForRestore.id);
      setSelectedRecordForRestore(null);
      if (result === 'completed') toast.success('Database restore completed');
    } catch {
      toast.error('Failed to trigger restore');
    }
  };

  return (
    <>
      <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-sm">
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <History className="h-4 w-4 text-primary" />
            <h2 className="font-semibold text-foreground text-sm">Recent snapshots & runs</h2>
          </div>
          <span className="text-muted-foreground text-xs">{records.length} total records</span>
        </div>

        {isLoading ? (
          <div className="flex min-h-48 items-center justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : records.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
            <History className="mb-3 h-8 w-8 opacity-20" />
            <p className="font-medium text-sm">No backup runs recorded yet</p>
            <p className="mt-1 text-xs">
              Automated and manual snapshots will show status, size, and restore actions here.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Status</TableHead>
                  <TableHead>Target</TableHead>
                  <TableHead>Size</TableHead>
                  <TableHead>Duration</TableHead>
                  <TableHead>Timestamp</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {records.map((record) => {
                  const isCompleted = record.status === 'completed';
                  const isFailed = record.status === 'failed';

                  return (
                    <TableRow key={record.id} className="hover:bg-muted/40">
                      <TableCell>
                        <span
                          className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 font-medium text-xs ${
                            isCompleted
                              ? 'bg-emerald-500/10 text-emerald-500'
                              : isFailed
                                ? 'bg-rose-500/10 text-rose-500'
                                : 'bg-sky-500/10 text-sky-500'
                          }`}
                        >
                          {isCompleted ? (
                            <CheckCircle2 className="h-3.5 w-3.5" />
                          ) : isFailed ? (
                            <XCircle className="h-3.5 w-3.5" />
                          ) : (
                            <Loader2 className="h-3.5 w-3.5 animate-spin" />
                          )}
                          <span className="capitalize">{record.status}</span>
                        </span>
                      </TableCell>
                      <TableCell className="font-medium text-xs">
                        {record.databaseId ? `db-${record.databaseId.slice(0, 8)}` : 'System state'}
                      </TableCell>
                      <TableCell className="font-mono text-xs">
                        {formatBytes(record.fileSizeBytes)}
                      </TableCell>
                      <TableCell className="text-muted-foreground text-xs">
                        {calculateDuration(record.startedAt, record.completedAt)}
                      </TableCell>
                      <TableCell className="whitespace-nowrap text-muted-foreground text-xs">
                        {record.startedAt ? new Date(record.startedAt).toLocaleString() : 'Unknown'}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-1.5">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-7 text-xs"
                            onClick={() => setSelectedRecordForRestore(record)}
                            disabled={!isCompleted || restorePending}
                            title="Restore this snapshot"
                          >
                            <ArchiveRestore className="mr-1 h-3.5 w-3.5" />
                            Restore
                          </Button>

                          {record.s3Url ? (
                            <Button variant="ghost" size="icon" className="h-7 w-7" asChild>
                              <a
                                href={record.s3Url}
                                target="_blank"
                                rel="noreferrer"
                                title="Download S3 archive"
                              >
                                <Download className="h-3.5 w-3.5 text-muted-foreground hover:text-foreground" />
                              </a>
                            </Button>
                          ) : record.filePath ? (
                            <Button variant="ghost" size="icon" className="h-7 w-7" asChild>
                              <a
                                href={`${getApiBaseUrl()}/backups/${record.backupConfigId}/records/${record.id}/download`}
                                target="_blank"
                                rel="noreferrer"
                                title="Download local archive"
                              >
                                <Download className="h-3.5 w-3.5 text-muted-foreground hover:text-foreground" />
                              </a>
                            </Button>
                          ) : null}

                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-7 w-7 text-destructive hover:bg-destructive/10 hover:text-destructive"
                            onClick={() => onDeleteRecord(record.backupConfigId, record.id)}
                            disabled={deletePending}
                            title="Delete snapshot"
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </div>
        )}
      </section>

      <Dialog
        open={Boolean(selectedRecordForRestore)}
        onOpenChange={(open) => !open && setSelectedRecordForRestore(null)}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <ArchiveRestore className="h-5 w-5 text-primary" />
              Restore snapshot confirmation
            </DialogTitle>
            <DialogDescription>
              Are you sure you want to restore this backup? This operation will replace current
              database tables and storage with the snapshot data taken on{' '}
              {selectedRecordForRestore?.startedAt
                ? new Date(selectedRecordForRestore.startedAt).toLocaleString()
                : 'selected date'}
              .
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2 sm:gap-0">
            <Button variant="outline" onClick={() => setSelectedRecordForRestore(null)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleConfirmRestore} disabled={restorePending}>
              {restorePending ? 'Restoring...' : 'Confirm restore'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
