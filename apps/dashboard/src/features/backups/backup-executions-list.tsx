import { format } from 'date-fns';
import { ArchiveRestore, Check, Download, History, Trash2 } from 'lucide-react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Section } from '#/components/ui/section';
import type { BackupRecord } from '#/interfaces/backup';

type Props = {
  records: BackupRecord[];
  configId: string;
  isLoadingRecords: boolean;
  handleRestore: (recordId: string) => void;
  restorePending: boolean;
  onDeleteRecord: (configId: string, recordId: string) => void;
  deletePending: boolean;
};

export function BackupExecutionsList({
  records,
  configId,
  isLoadingRecords,
  handleRestore,
  restorePending,
  onDeleteRecord,
  deletePending,
}: Props) {
  return (
    <Section icon={<History className="h-4 w-4" />} title={`Executions (${records.length})`}>
      <div className="flex flex-col gap-4 py-4">
        <div className="flex flex-col gap-4">
          {isLoadingRecords ? (
            <div className="py-8 text-center text-muted-foreground">Loading executions...</div>
          ) : records.length === 0 ? (
            <div className="py-8 text-center text-muted-foreground">No executions yet.</div>
          ) : (
            records.map((record) => (
              <div
                key={record.id}
                className="border/50 flex flex-col gap-3 rounded-lg border bg-background/50 p-4"
              >
                <Badge
                  variant="outline"
                  className={
                    record.status === 'completed'
                      ? 'w-fit border-green-500/20 bg-green-500/10 text-green-500'
                      : record.status === 'failed'
                        ? 'w-fit border-red-500/20 bg-red-500/10 text-red-500'
                        : 'w-fit border-yellow-500/20 bg-yellow-500/10 text-yellow-500'
                  }
                >
                  {record.status === 'completed'
                    ? 'Success'
                    : record.status === 'failed'
                      ? 'Failed'
                      : 'Running'}
                </Badge>

                <div className="text-muted-foreground text-sm leading-relaxed">
                  {record.startedAt
                    ? format(new Date(record.startedAt), 'MMM d, HH:mm')
                    : 'Unknown time'}{' '}
                  • Database: codedock • Size: {(record.fileSizeBytes / 1024 / 1024).toFixed(2)} MB
                  <br />
                  Location: {record.filePath}
                </div>

                <div className="flex items-center gap-2 text-sm">
                  <span className="text-muted-foreground">Backup Availability:</span>
                  <Badge
                    variant="outline"
                    className="gap-1 border-green-500/20 bg-green-500/10 text-green-500"
                  >
                    <Check className="h-3 w-3" /> Local Storage
                  </Badge>
                </div>

                <div className="mt-2 flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleRestore(record.id)}
                    disabled={restorePending || record.status !== 'completed' || !record.filePath}
                  >
                    <ArchiveRestore className="mr-2 h-4 w-4" />
                    Restore
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    asChild
                    disabled={!record.s3Url && !record.filePath}
                  >
                    {record.s3Url ? (
                      <a href={record.s3Url} target="_blank" rel="noreferrer">
                        <Download className="mr-2 h-4 w-4" />
                        Download S3
                      </a>
                    ) : record.filePath ? (
                      <a
                        href={`${import.meta.env.VITE_API_URL}/backups/${configId}/records/${record.id}/download`}
                        target="_blank"
                        rel="noreferrer"
                      >
                        <Download className="mr-2 h-4 w-4" />
                        Download Local
                      </a>
                    ) : (
                      <span>
                        <Download className="mr-2 h-4 w-4" />
                        Download
                      </span>
                    )}
                  </Button>
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => onDeleteRecord(configId, record.id)}
                    disabled={deletePending}
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    {deletePending ? 'Deleting...' : 'Delete'}
                  </Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </Section>
  );
}
