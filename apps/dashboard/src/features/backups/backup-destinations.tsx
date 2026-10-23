import { HardDrive, Loader2, Trash2 } from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Skeleton } from '@codedock/ui/skeleton';
import { BackupDestinationRow } from './backup-destination-row';
import type { UnifiedDestination } from './destination-display';
import {
  DestinationLoadError,
  toUnifiedS3,
  toUnifiedSftp,
} from './destination-display';
import { DestinationModal } from './destination-modal';
import {
  useDeleteS3Destination,
  useDeleteSFTPDestination,
  useSetDefaultS3Destination,
  useVerifyS3Destination,
  useVerifySFTPDestination,
} from '#/hooks/use-backups';
import type { BackupRecord, S3Destination, SFTPDestination } from '#/interfaces/backups';

export function BackupDestinations({
  s3,
  sftp,
  records,
  isLoading,
  loadError,
  onRetry,
  retrying,
  onChanged,
}: {
  s3: S3Destination[];
  sftp: SFTPDestination[];
  records: BackupRecord[];
  isLoading: boolean;
  loadError: string | null;
  onRetry: () => void;
  retrying: boolean;
  onChanged: () => void;
}) {
  const [verifyingIds, setVerifyingIds] = useState<Set<string>>(new Set());
  const [editing, setEditing] = useState<S3Destination | null>(null);
  const [deleting, setDeleting] = useState<UnifiedDestination | null>(null);
  const [deleteBusy, setDeleteBusy] = useState(false);

  const verifyS3 = useVerifyS3Destination();
  const verifySftp = useVerifySFTPDestination();
  const setDefault = useSetDefaultS3Destination();
  const deleteS3 = useDeleteS3Destination();
  const deleteSftp = useDeleteSFTPDestination();

  const s3ById = useMemo(
    () => new Map(s3.map((destination) => [destination.id, destination])),
    [s3]
  );
  const items = useMemo<UnifiedDestination[]>(
    () =>
      [...s3.map(toUnifiedS3), ...sftp.map(toUnifiedSftp)].sort((a, b) =>
        a.name.localeCompare(b.name)
      ),
    [s3, sftp]
  );
  const statsByDestination = useMemo(() => {
    const stats = new Map<
      string,
      { bytes: number; runs: number; active: number; failed: number }
    >();
    for (const record of records) {
      const id = record.s3DestinationId || record.sftpDestinationId;
      if (!id) continue;
      const entry = stats.get(id) ?? { bytes: 0, runs: 0, active: 0, failed: 0 };
      entry.runs += 1;
      if (record.status === 'completed') entry.bytes += record.fileSizeBytes || 0;
      if (record.status === 'running') entry.active += 1;
      if (record.status === 'failed') entry.failed += 1;
      stats.set(id, entry);
    }
    return stats;
  }, [records]);

  const handleVerify = async (row: UnifiedDestination) => {
    setVerifyingIds((previous) => new Set(previous).add(row.id));
    try {
      const res =
        row.kind === 's3'
          ? await verifyS3.mutateAsync(row.id)
          : await verifySftp.mutateAsync(row.id);
      if (res.data?.ok) toast.success(`"${row.name}" verified successfully`);
      else toast.error(res.data?.reason || 'Verification failed');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Verification failed');
    } finally {
      setVerifyingIds((previous) => {
        const next = new Set(previous);
        next.delete(row.id);
        return next;
      });
    }
  };

  const handleSetDefault = async (row: UnifiedDestination) => {
    if (row.kind !== 's3') return;
    try {
      await setDefault.mutateAsync(row.id);
      toast.success(`"${row.name}" set as default`);
      onChanged();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to set default');
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setDeleteBusy(true);
    try {
      if (deleting.kind === 's3') {
        await deleteS3.mutateAsync({ id: deleting.id, projectId: 'global' });
      } else {
        await deleteSftp.mutateAsync(deleting.id);
      }
      toast.success(`Destination "${deleting.name}" deleted`);
      setDeleting(null);
      onChanged();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to delete destination');
    } finally {
      setDeleteBusy(false);
    }
  };

  return (
    <section aria-label="Destinations" className="overflow-hidden rounded-2xl border border-border/50 bg-card shadow-xs">
      <div className="px-5 py-4">
        <h2 className="font-medium text-foreground text-sm">Destinations</h2>
        <p className="sr-only">
          S3-compatible buckets and SFTP servers holding off-server snapshots.
        </p>
      </div>

      {loadError && (
        <div className="px-5">
          <DestinationLoadError message={loadError} onRetry={onRetry} busy={retrying} />
        </div>
      )}

      {isLoading ? (
        <div aria-busy="true" className="space-y-3 px-5 pb-5">
          {[0, 1].map((row) => (
            <Skeleton key={row} className="h-20 w-full" />
          ))}
        </div>
      ) : items.length === 0 ? (
        <div className="flex items-center gap-4 px-5 pt-3 pb-6">
          <div className="flex size-11 shrink-0 items-center justify-center rounded-xl bg-background text-muted-foreground">
            <HardDrive className="size-5" />
          </div>
          <div>
            <p className="font-medium text-foreground text-sm">No backup destinations yet</p>
            <p className="mt-1 text-muted-foreground text-sm">
              Add an S3 or SFTP target to protect service volumes and cluster databases.
            </p>
          </div>
        </div>
      ) : (
        <ul className="divide-y divide-border/40">
          {items.map((row) => (
            <BackupDestinationRow
              key={`${row.kind}-${row.id}`}
              row={row}
              stats={statsByDestination.get(row.id)}
              verifying={verifyingIds.has(row.id)}
              s3Destination={row.kind === 's3' ? s3ById.get(row.id) : undefined}
              onVerify={(r) => void handleVerify(r)}
              onEdit={(s) => setEditing(s)}
              onSetDefault={(r) => void handleSetDefault(r)}
              onDelete={(r) => setDeleting(r)}
            />
          ))}
        </ul>
      )}

      <DestinationModal
        open={editing !== null}
        onOpenChange={(open) => !open && setEditing(null)}
        editing={editing}
        onSaved={onChanged}
      />

      <Dialog
        open={deleting !== null}
        onOpenChange={(open) => !open && !deleteBusy && setDeleting(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete "{deleting?.name}"?</DialogTitle>
            <DialogDescription>
              The destination is removed from Codedock. Snapshots already stored there are left
              untouched, but policies pointing at it stop working.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleting(null)} disabled={deleteBusy}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={confirmDelete} disabled={deleteBusy}>
              {deleteBusy ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Trash2 className="h-4 w-4" />
              )}
              Delete destination
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </section>
  );
}
