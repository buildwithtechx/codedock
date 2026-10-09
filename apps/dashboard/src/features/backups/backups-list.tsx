import { useSearch } from '@tanstack/react-router';
import { Plus, RefreshCw } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import {
  useList,
  useListAllRecords,
  useListS3Destinations,
  useListSFTPDestinations,
} from '#/features/backups';
import { BackupDestinationHistory } from './backup-destination-history';
import { BackupDestinations } from './backup-destinations';
import { BackupOperationHistory } from './backup-operation-history';
import { BackupPolicies } from './backup-policies';
import { BackupStorageSummary } from './backup-storage-summary';
import { DestinationModal } from './destination-modal';

export function BackupsList() {
  const search = useSearch({ strict: false }) as { add?: string } | undefined;
  const [isModalOpen, setIsModalOpen] = useState(search?.add === 'true');

  useEffect(() => {
    if (search?.add === 'true') setIsModalOpen(true);
  }, [search?.add]);

  const configsQuery = useList();
  const recordsQuery = useListAllRecords(50);
  const s3Query = useListS3Destinations();
  const sftpQuery = useListSFTPDestinations();

  const configs = configsQuery.data?.data ?? [];
  const records = recordsQuery.data?.data ?? [];
  const s3 = s3Query.data?.data ?? [];
  const sftp = sftpQuery.data?.data ?? [];
  const loading = configsQuery.isLoading || recordsQuery.isLoading || s3Query.isLoading;
  const refreshing =
    configsQuery.isFetching ||
    recordsQuery.isFetching ||
    s3Query.isFetching ||
    sftpQuery.isFetching;

  const destinationNames = useMemo(() => {
    const names = new Map<string, string>();
    for (const destination of s3) names.set(destination.id, destination.name);
    for (const destination of sftp) names.set(destination.id, destination.name);
    return names;
  }, [s3, sftp]);

  const refreshAll = async () => {
    await Promise.all([
      configsQuery.refetch(),
      recordsQuery.refetch(),
      s3Query.refetch(),
      sftpQuery.refetch(),
    ]);
  };

  const handleRefresh = async () => {
    await refreshAll();
    toast.success('Backups refreshed');
  };

  return (
    <div className="space-y-6">
      <BackupOperationHistory />

      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="font-medium text-2xl">Backups</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Snapshots, off-server destinations, and automated policies.
          </p>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={handleRefresh}
            disabled={refreshing}
            aria-label="Refresh backups"
            title="Refresh backups"
          >
            <RefreshCw className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
          </Button>
          <Button onClick={() => setIsModalOpen(true)}>
            <Plus className="h-4 w-4" />
            Add destination
          </Button>
        </div>
      </div>

      {loading ? (
        <div aria-busy="true" className="grid gap-5 xl:grid-cols-[minmax(0,1fr)_300px]">
          <Skeleton className="h-72 w-full rounded-2xl" />
          <Skeleton className="h-52 w-full rounded-2xl" />
        </div>
      ) : (
        <div className="grid items-start gap-5 xl:grid-cols-[minmax(0,1fr)_300px]">
          <div className="min-w-0 space-y-5">
            <BackupDestinationHistory
              records={records}
              configs={configs}
              destinations={destinationNames}
              isLoading={recordsQuery.isLoading}
              loadError={recordsQuery.isError ? 'Could not load snapshot history.' : null}
              onRetry={() => void recordsQuery.refetch()}
              retrying={recordsQuery.isFetching}
              onChanged={() => void refreshAll()}
            />
            <BackupDestinations
              s3={s3}
              sftp={sftp}
              records={records}
              isLoading={s3Query.isLoading || sftpQuery.isLoading}
              loadError={
                s3Query.isError || sftpQuery.isError ? 'Could not load destinations.' : null
              }
              onRetry={() => {
                void s3Query.refetch();
                void sftpQuery.refetch();
              }}
              retrying={s3Query.isFetching || sftpQuery.isFetching}
              onChanged={() => void refreshAll()}
              onAddDestination={() => setIsModalOpen(true)}
            />
            <BackupPolicies configs={configs} isLoading={configsQuery.isLoading} />
          </div>
          <aside className="space-y-5 xl:sticky xl:top-6">
            <BackupStorageSummary
              records={records}
              destinations={s3}
              configs={configs}
              sftpCount={sftp.length}
            />
          </aside>
        </div>
      )}

      <DestinationModal open={isModalOpen} onOpenChange={setIsModalOpen} onSaved={refreshAll} />
    </div>
  );
}
