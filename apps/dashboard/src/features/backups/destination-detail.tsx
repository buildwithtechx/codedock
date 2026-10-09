import { Link } from '@tanstack/react-router';
import {
  ArrowLeft,
  BookOpen,
  CheckCircle2,
  Cloud,
  ExternalLink,
  HardDrive,
  Loader2,
  Lock,
  Pencil,
  Play,
  RefreshCw,
  Server,
  Star,
} from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '#/components/ui/empty';
import { Skeleton } from '#/components/ui/skeleton';
import { cn } from '#/lib/utils';
import {
  DestinationLoadError,
  DestinationVerificationBadge,
  kindLabel,
  toUnifiedS3,
  toUnifiedSftp,
} from './destination-display';
import { DestinationModal } from './destination-modal';
import { DestinationSnapshots } from './destination-snapshots';
import {
  useList,
  useListAllRecords,
  useListS3Destinations,
  useListSFTPDestinations,
  useTrigger,
  useVerifyS3Destination,
  useVerifySFTPDestination,
} from './hooks';

const DOCS_URL = 'https://docs.codedock.run';

export function DestinationDetail({ destinationId }: { destinationId: string }) {
  const s3Query = useListS3Destinations();
  const sftpQuery = useListSFTPDestinations();
  const configsQuery = useList();
  const recordsQuery = useListAllRecords(100);
  const trigger = useTrigger();
  const verifyS3 = useVerifyS3Destination();
  const verifySftp = useVerifySFTPDestination();

  const [verifying, setVerifying] = useState(false);
  const [editing, setEditing] = useState(false);
  const [triggeringId, setTriggeringId] = useState<string | null>(null);

  const s3 = useMemo(
    () => (s3Query.data?.data ?? []).find((destination) => destination.id === destinationId),
    [s3Query.data, destinationId]
  );
  const sftp = useMemo(
    () => (sftpQuery.data?.data ?? []).find((destination) => destination.id === destinationId),
    [sftpQuery.data, destinationId]
  );
  const unified = useMemo(
    () => (s3 ? toUnifiedS3(s3) : sftp ? toUnifiedSftp(sftp) : null),
    [s3, sftp]
  );
  const configs = useMemo(
    () =>
      (configsQuery.data?.data ?? []).filter(
        (config) =>
          config.s3DestinationId === destinationId || config.sftpDestinationId === destinationId
      ),
    [configsQuery.data, destinationId]
  );
  const snapshots = useMemo(
    () =>
      (recordsQuery.data?.data ?? []).filter(
        (record) =>
          record.s3DestinationId === destinationId || record.sftpDestinationId === destinationId
      ),
    [recordsQuery.data, destinationId]
  );
  const storedBytes = snapshots
    .filter((record) => record.status === 'completed')
    .reduce((total, record) => total + (record.fileSizeBytes || 0), 0);

  const loading = s3Query.isLoading || sftpQuery.isLoading;
  const refreshing =
    s3Query.isFetching ||
    sftpQuery.isFetching ||
    configsQuery.isFetching ||
    recordsQuery.isFetching;
  const loadError =
    s3Query.isError || sftpQuery.isError ? 'Could not load this destination.' : null;

  const reload = () => {
    void s3Query.refetch();
    void sftpQuery.refetch();
    void configsQuery.refetch();
    void recordsQuery.refetch();
  };

  const handleVerify = async () => {
    if (!unified || verifying) return;
    setVerifying(true);
    try {
      const res =
        unified.kind === 's3'
          ? await verifyS3.mutateAsync(unified.id)
          : await verifySftp.mutateAsync(unified.id);
      if (res.data?.ok) toast.success(`"${unified.name}" verified successfully`);
      else toast.error(res.data?.reason || 'Verification failed');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Verification failed');
    } finally {
      setVerifying(false);
      reload();
    }
  };

  const handleTrigger = async (configId: string, name: string) => {
    setTriggeringId(configId);
    try {
      await trigger.mutateAsync({ id: configId });
      toast.success(`Snapshot started for "${name}"`);
      void recordsQuery.refetch();
    } catch {
      toast.error('Failed to trigger snapshot');
    } finally {
      setTriggeringId(null);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!unified) {
    return (
      <Empty>
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <HardDrive />
          </EmptyMedia>
          <EmptyTitle>Destination not found</EmptyTitle>
          <EmptyDescription>
            This destination no longer exists or you don't have access to it.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <Button asChild variant="outline">
            <Link to="/backups">Back to Backups</Link>
          </Button>
        </EmptyContent>
      </Empty>
    );
  }

  const Icon = unified.kind === 's3' ? Cloud : Server;

  return (
    <div className="space-y-6">
      <Link
        to="/backups"
        className="inline-flex items-center gap-1.5 text-muted-foreground text-sm transition-colors hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        Backups
      </Link>

      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex min-w-0 items-center gap-3">
          <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
            <Icon className="h-5 w-5" />
          </div>
          <div className="min-w-0">
            <h1 className="truncate font-medium text-2xl" title={unified.name}>
              {unified.name}
            </h1>
            <div className="mt-1 flex flex-wrap items-center gap-2 text-muted-foreground text-xs">
              <span>{kindLabel(unified.kind)}</span>
              {unified.isDefault && (
                <span className="inline-flex items-center gap-1">
                  <Star className="h-3 w-3" />
                  Default
                </span>
              )}
              <DestinationVerificationBadge destination={unified} verifying={verifying} />
            </div>
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={reload}
            disabled={refreshing}
            title="Refresh"
          >
            <RefreshCw className={cn('h-4 w-4', refreshing && 'animate-spin')} />
          </Button>
          <Button variant="outline" onClick={handleVerify} disabled={verifying}>
            {verifying ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <CheckCircle2 className="h-4 w-4" />
            )}
            Verify connection
          </Button>
          {s3 && (
            <Button variant="outline" onClick={() => setEditing(true)}>
              <Pencil className="h-4 w-4" />
              Edit
            </Button>
          )}
        </div>
      </div>

      {loadError && <DestinationLoadError message={loadError} onRetry={reload} busy={refreshing} />}

      <div className="grid items-start gap-5 xl:grid-cols-[minmax(0,1fr)_300px]">
        <div className="min-w-0 space-y-5">
          <DestinationSnapshots
            snapshots={snapshots}
            storedBytes={storedBytes}
            isLoading={recordsQuery.isLoading}
            destinationName={unified.name}
            onChanged={() => void recordsQuery.refetch()}
          />

          <section aria-label="Used by" className="rounded-2xl border border-border/50 bg-card">
            <div className="px-5 py-4">
              <h2 className="font-medium text-base">Used by</h2>
              <p className="mt-1 text-muted-foreground text-sm">
                Backup policies writing snapshots to this destination.
              </p>
            </div>
            {configsQuery.isLoading ? (
              <div aria-busy="true" className="space-y-3 px-5 pb-5">
                <Skeleton className="h-12 w-full" />
              </div>
            ) : configs.length === 0 ? (
              <div className="px-5 pb-6">
                <p className="font-medium text-sm">No policies use this destination</p>
                <p className="mt-1 text-muted-foreground text-sm">
                  Point a policy at it from the database backup settings.
                </p>
                <Button asChild size="sm" variant="secondary" className="mt-3">
                  <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
                    <BookOpen className="size-3.5" />
                    Docs
                    <ExternalLink className="size-3 opacity-60" />
                  </a>
                </Button>
              </div>
            ) : (
              <ul className="divide-y divide-border/40 border-border/40 border-t">
                {configs.map((config) => (
                  <li
                    key={config.id}
                    className="flex flex-wrap items-center gap-3 px-5 py-4 last:rounded-b-2xl hover:bg-muted/20"
                  >
                    <div className="min-w-0 flex-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <p className="truncate font-medium text-sm">{config.name}</p>
                        {!config.backupEnabled && (
                          <span className="rounded-md bg-muted px-1.5 py-0.5 font-medium text-muted-foreground text-xs">
                            Paused
                          </span>
                        )}
                      </div>
                      <code className="mt-1.5 block font-mono text-muted-foreground text-xs">
                        {config.schedule || 'Manual only'}
                      </code>
                    </div>
                    <Button
                      variant="outline"
                      size="sm"
                      className="h-7 shrink-0 text-xs"
                      disabled={triggeringId !== null}
                      onClick={() => void handleTrigger(config.id, config.name)}
                    >
                      {triggeringId === config.id ? (
                        <Loader2 className="h-3 w-3 animate-spin" />
                      ) : (
                        <Play className="h-3 w-3" />
                      )}
                      Snapshot now
                    </Button>
                  </li>
                ))}
              </ul>
            )}
          </section>
        </div>

        <aside className="space-y-5 xl:sticky xl:top-6">
          <section
            aria-label="Destination details"
            className="rounded-2xl border border-border/50 bg-card p-5"
          >
            <h2 className="font-medium text-base">Destination details</h2>
            <div className="mt-4 space-y-4 text-sm">
              <div className="flex items-start gap-2.5 text-muted-foreground">
                <Icon className="mt-0.5 h-4 w-4 shrink-0" />
                <p className="min-w-0 break-words">{unified.detail}</p>
              </div>
              <div className="flex items-start gap-2.5 text-muted-foreground">
                <Lock className="mt-0.5 h-4 w-4 shrink-0" />
                <p>{unified.credential}</p>
              </div>
              {unified.lastVerifyError && (
                <p className="break-words text-destructive text-sm">{unified.lastVerifyError}</p>
              )}
              {unified.description && (
                <p className="text-muted-foreground text-sm">{unified.description}</p>
              )}
            </div>
          </section>
        </aside>
      </div>

      <DestinationModal
        open={editing}
        onOpenChange={setEditing}
        editing={s3 ?? null}
        onSaved={reload}
      />
    </div>
  );
}
