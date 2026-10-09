import { Link } from '@tanstack/react-router';
import {
  BookOpen,
  ExternalLink,
  HardDrive,
  Loader2,
  MoreVertical,
  Pencil,
  RefreshCw,
  Star,
  Trash2,
} from 'lucide-react';
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Skeleton } from '#/components/ui/skeleton';
import type { UnifiedDestination } from './destination-display';
import {
  DestinationKindIcon,
  DestinationLoadError,
  DestinationVerificationBadge,
  formatBytes,
  kindLabel,
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
} from './hooks';
import type { BackupRecord, S3Destination, SFTPDestination } from './interfaces';

const DOCS_URL = 'https://docs.codedock.run';

export function BackupDestinations({
  s3,
  sftp,
  records,
  isLoading,
  loadError,
  onRetry,
  retrying,
  onChanged,
  onAddDestination,
}: {
  s3: S3Destination[];
  sftp: SFTPDestination[];
  records: BackupRecord[];
  isLoading: boolean;
  loadError: string | null;
  onRetry: () => void;
  retrying: boolean;
  onChanged: () => void;
  onAddDestination: () => void;
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
      onChanged();
    }
  };

  const handleSetDefault = async (row: UnifiedDestination) => {
    try {
      await setDefault.mutateAsync(row.id);
      toast.success(`"${row.name}" is now the default destination`);
      onChanged();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to set default destination');
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
    <section aria-label="Destinations" className="rounded-2xl border border-border/50 bg-card">
      <div className="px-5 py-4">
        <h2 className="font-medium text-base">Destinations</h2>
        <p className="mt-1 text-muted-foreground text-sm">
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
        <div className="flex items-center gap-4 px-5 pb-6">
          <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
            <HardDrive className="h-5 w-5" />
          </div>
          <div className="min-w-0 flex-1">
            <p className="font-medium text-sm">No storage destinations</p>
            <p className="mt-1 text-muted-foreground text-sm">
              Connect R2, S3, MinIO, or an SFTP server to store snapshots off-server.
            </p>
          </div>
          <div className="flex shrink-0 items-center gap-2">
            <Button size="sm" onClick={onAddDestination}>
              Connect first destination
            </Button>
            <Button asChild size="sm" variant="secondary">
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
          {items.map((row) => {
            const stats = statsByDestination.get(row.id);
            const verifying = verifyingIds.has(row.id);
            const s3Destination = row.kind === 's3' ? s3ById.get(row.id) : undefined;
            return (
              <li
                key={`${row.kind}-${row.id}`}
                className="group relative flex items-start gap-3 px-5 py-4 transition-colors last:rounded-b-2xl hover:bg-muted/20"
              >
                <Link
                  to="/backups/$backupId"
                  params={{ backupId: row.id }}
                  aria-label={row.name}
                  className="absolute inset-0 z-0 rounded-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                />
                <div className="mt-0.5 flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
                  <DestinationKindIcon kind={row.kind} className="h-5 w-5" />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="truncate font-medium text-foreground text-sm">{row.name}</p>
                    {row.isDefault && (
                      <span className="inline-flex items-center gap-1 text-muted-foreground text-xs">
                        <Star className="h-3 w-3" />
                        Default
                      </span>
                    )}
                    <DestinationVerificationBadge destination={row} verifying={verifying} />
                  </div>
                  <p className="mt-1 truncate text-muted-foreground text-xs" title={row.detail}>
                    {kindLabel(row.kind)} · {row.detail}
                  </p>
                  <div className="mt-1.5 flex flex-wrap items-center gap-x-2 gap-y-1 text-muted-foreground text-xs">
                    {stats && stats.runs > 0 ? (
                      <>
                        <span>
                          {formatBytes(stats.bytes)} stored · {stats.runs} runs
                        </span>
                        {stats.active > 0 && (
                          <span className="text-sky-500">· {stats.active} active</span>
                        )}
                        {stats.failed > 0 && (
                          <span className="text-destructive">· {stats.failed} failed</span>
                        )}
                      </>
                    ) : (
                      <span>No runs yet</span>
                    )}
                  </div>
                  {row.lastVerifyError && (
                    <p
                      className="mt-1 truncate text-destructive text-xs"
                      title={row.lastVerifyError}
                    >
                      {row.lastVerifyError}
                    </p>
                  )}
                </div>
                <div className="relative z-10 flex shrink-0 items-center gap-1">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => void handleVerify(row)}
                    disabled={verifying}
                    title="Verify connection"
                    aria-label={`Verify connection: ${row.name}`}
                  >
                    {verifying ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <RefreshCw className="h-4 w-4" />
                    )}
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" aria-label={`Actions: ${row.name}`}>
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      {s3Destination && (
                        <DropdownMenuItem onClick={() => setEditing(s3Destination)}>
                          <Pencil className="h-4 w-4" />
                          Edit
                        </DropdownMenuItem>
                      )}
                      {row.kind === 's3' && !row.isDefault && (
                        <DropdownMenuItem onClick={() => void handleSetDefault(row)}>
                          <Star className="h-4 w-4" />
                          Set as default
                        </DropdownMenuItem>
                      )}
                      {(s3Destination || (row.kind === 's3' && !row.isDefault)) && (
                        <DropdownMenuSeparator />
                      )}
                      <DropdownMenuItem variant="destructive" onClick={() => setDeleting(row)}>
                        <Trash2 className="h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </li>
            );
          })}
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
