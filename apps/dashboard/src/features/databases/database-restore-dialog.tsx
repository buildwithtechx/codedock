import { AlertTriangle, ArchiveRestore, Loader2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import type { BackupRecord } from '#/features/backups/interfaces';

type DatabaseRestoreDialogProps = {
  record: BackupRecord | null;
  databaseName: string;
  onClose: () => void;
  onConfirm: (recordId: string) => Promise<void>;
  isPending: boolean;
};

export function DatabaseRestoreDialog({
  record,
  databaseName,
  onClose,
  onConfirm,
  isPending,
}: DatabaseRestoreDialogProps) {
  if (!record) return null;

  return (
    <Dialog open={Boolean(record)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ArchiveRestore className="h-5 w-5 text-primary" />
            Restore Database Snapshot
          </DialogTitle>
          <DialogDescription>
            You are about to restore snapshot taken on{' '}
            <span className="font-semibold text-foreground">
              {record.startedAt ? new Date(record.startedAt).toLocaleString() : 'selected date'}
            </span>{' '}
            into database <span className="font-semibold text-foreground">{databaseName}</span>.
          </DialogDescription>
        </DialogHeader>

        <div className="rounded-xl border border-amber-500/20 bg-amber-500/10 p-3.5 text-amber-600 text-xs dark:text-amber-400">
          <div className="flex items-start gap-2.5">
            <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0 text-amber-500" />
            <div className="space-y-1">
              <p className="font-semibold">Warning: This operation will overwrite existing data</p>
              <p className="text-muted-foreground">
                All existing tables and rows will be replaced with the contents of this snapshot.
                Any data created after this snapshot was taken will be lost.
              </p>
            </div>
          </div>
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button variant="outline" onClick={onClose} disabled={isPending}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={() => onConfirm(record.id)}
            disabled={isPending}
            className="gap-2"
          >
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {isPending ? 'Restoring snapshot...' : 'Confirm Restore'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
