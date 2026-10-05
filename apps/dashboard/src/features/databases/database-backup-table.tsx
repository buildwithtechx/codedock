import {
  ArchiveRestore,
  CheckCircle2,
  DatabaseBackup,
  Download,
  HardDrive,
  Loader2,
  Play,
  Trash2,
  XCircle,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Card, CardContent } from '#/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import type { BackupRecord } from '#/features/backups/interfaces';
import { formatBytes } from '#/lib/utils';

type DatabaseBackupTableProps = {
  records: BackupRecord[];
  onRestore: (record: BackupRecord) => void;
  onDelete: (configId: string, recordId: string) => void;
  onBackupNow: () => void;
  restorePending: boolean;
  deletePending: boolean;
  backupPending: boolean;
};

export function DatabaseBackupTable({
  records,
  onRestore,
  onDelete,
  onBackupNow,
  restorePending,
  deletePending,
  backupPending,
}: DatabaseBackupTableProps) {
  return (
    <Card className="border-border/60 bg-card">
      <CardContent className="p-0">
        {records.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-muted/60 text-muted-foreground">
              <DatabaseBackup className="h-6 w-6" />
            </div>
            <h3 className="mt-3 font-semibold text-foreground text-sm">No snapshots recorded</h3>
            <p className="mt-1 max-w-sm text-xs leading-5">
              Run an immediate snapshot or configure automated backups to protect your database.
            </p>
            <Button
              size="sm"
              className="mt-4 gap-1.5"
              onClick={onBackupNow}
              disabled={backupPending}
            >
              <Play className="h-4 w-4 fill-current" />
              Take first snapshot
            </Button>
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-36">Status</TableHead>
                <TableHead>Storage</TableHead>
                <TableHead>Size</TableHead>
                <TableHead>Timestamp</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {records.map((record) => {
                const isCompleted = record.status === 'completed';
                const isFailed = record.status === 'failed';

                return (
                  <TableRow key={record.id} className="hover:bg-muted/30">
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
                    <TableCell>
                      <div className="flex items-center gap-1.5 text-xs">
                        {record.s3Url ? (
                          <span className="flex items-center gap-1 text-muted-foreground">
                            <HardDrive className="h-3.5 w-3.5 text-primary" /> S3 Bucket
                          </span>
                        ) : (
                          <span className="flex items-center gap-1 text-muted-foreground">
                            <HardDrive className="h-3.5 w-3.5" /> Server Disk
                          </span>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="font-mono text-xs">
                      {formatBytes(record.fileSizeBytes)}
                    </TableCell>
                    <TableCell className="text-muted-foreground text-xs">
                      {record.startedAt ? new Date(record.startedAt).toLocaleString() : '—'}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        <Button
                          variant="outline"
                          size="sm"
                          className="h-7 text-xs"
                          onClick={() => onRestore(record)}
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
                              href={`${import.meta.env.VITE_API_URL || ''}/backups/${record.backupConfigId}/records/${record.id}/download`}
                              target="_blank"
                              rel="noreferrer"
                              title="Download archive"
                            >
                              <Download className="h-3.5 w-3.5 text-muted-foreground hover:text-foreground" />
                            </a>
                          </Button>
                        ) : null}

                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-7 w-7 text-destructive hover:bg-destructive/10 hover:text-destructive"
                          onClick={() => onDelete(record.backupConfigId, record.id)}
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
        )}
      </CardContent>
    </Card>
  );
}
