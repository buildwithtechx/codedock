import { useSearch } from '@tanstack/react-router';
import { Clock, HardDrive, History, Plus, RefreshCw } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import {
  useDeleteRecord,
  useList,
  useListAllRecords,
  useListS3Destinations,
  useRestore,
} from '#/features/backups';
import { BackupDestinationHistory } from './backup-destination-history';
import { BackupDestinations } from './backup-destinations';
import { BackupPolicies } from './backup-policies';
import { BackupStorageSummary } from './backup-storage-summary';
import { CreateS3DestinationDialog } from './create-s3-destination-dialog';

export function BackupsList() {
  const search = useSearch({ strict: false }) as { tab?: string; add?: string } | undefined;
  const [activeTab, setActiveTab] = useState(search?.tab || 'history');
  const [isDestinationDialogOpen, setIsDestinationDialogOpen] = useState(search?.add === 'true');

  useEffect(() => {
    if (search?.tab) {
      setActiveTab(search.tab);
    }
    if (search?.add === 'true') {
      setIsDestinationDialogOpen(true);
    }
  }, [search?.tab, search?.add]);

  const { data: configsData, isLoading: isLoadingConfigs, refetch: refetchConfigs } = useList();
  const configs = configsData?.data || [];

  const {
    data: recordsData,
    isLoading: isLoadingRecords,
    refetch: refetchRecords,
  } = useListAllRecords(50);
  const records = recordsData?.data || [];

  const { data: s3Data, refetch: refetchS3 } = useListS3Destinations();
  const destinations = s3Data?.data || [];

  const restoreMutation = useRestore();
  const deleteRecordMutation = useDeleteRecord();

  const handleRefreshAll = async () => {
    await Promise.all([refetchConfigs(), refetchRecords(), refetchS3()]);
    toast.success('Backup records and storage updated');
  };

  const handleRestoreRecord = async (recordId: string) => {
    const record = records.find((r) => r.id === recordId);
    if (!record) return;
    await restoreMutation.mutateAsync({ id: record.backupConfigId });
  };

  const handleDeleteRecord = async (configId: string, recordId: string) => {
    if (!window.confirm('Are you sure you want to delete this snapshot?')) return;
    try {
      await deleteRecordMutation.mutateAsync({ id: configId, recordId });
      toast.success('Snapshot deleted');
      void refetchRecords();
    } catch {
      toast.error('Failed to delete snapshot');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <PageHeader
          title="Backups"
          description="Manage S3 storage destinations, automated snapshot schedules, and disaster recovery."
        />
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefreshAll}
            className="gap-1.5"
            title="Refresh backups and storage"
          >
            <RefreshCw className="h-4 w-4" />
            Refresh
          </Button>
          <CreateS3DestinationDialog
            isOpen={isDestinationDialogOpen}
            setIsOpen={setIsDestinationDialogOpen}
            trigger={
              <Button size="sm" className="gap-1.5">
                <Plus className="h-4 w-4" />
                Add destination
              </Button>
            }
          />
        </div>
      </div>

      <div className="grid items-start gap-6 xl:grid-cols-[minmax(0,1fr)_320px]">
        <div className="min-w-0 space-y-6">
          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="history" className="gap-1.5 text-xs sm:text-sm">
                <History className="h-4 w-4" />
                <span>Snapshots</span>
                <span className="ml-1 rounded-full bg-muted px-1.5 py-0.2 font-mono text-[10px]">
                  {records.length}
                </span>
              </TabsTrigger>
              <TabsTrigger value="destinations" className="gap-1.5 text-xs sm:text-sm">
                <HardDrive className="h-4 w-4" />
                <span>Destinations</span>
                <span className="ml-1 rounded-full bg-muted px-1.5 py-0.2 font-mono text-[10px]">
                  {destinations.length}
                </span>
              </TabsTrigger>
              <TabsTrigger value="policies" className="gap-1.5 text-xs sm:text-sm">
                <Clock className="h-4 w-4" />
                <span>Policies</span>
                <span className="ml-1 rounded-full bg-muted px-1.5 py-0.2 font-mono text-[10px]">
                  {configs.length}
                </span>
              </TabsTrigger>
            </TabsList>

            <TabsContent value="history" className="mt-4">
              <BackupDestinationHistory
                records={records}
                isLoading={isLoadingRecords}
                onRestore={handleRestoreRecord}
                onDeleteRecord={handleDeleteRecord}
                restorePending={restoreMutation.isPending}
                deletePending={deleteRecordMutation.isPending}
              />
            </TabsContent>

            <TabsContent value="destinations" className="mt-4">
              <BackupDestinations onAddDestination={() => setIsDestinationDialogOpen(true)} />
            </TabsContent>

            <TabsContent value="policies" className="mt-4">
              <BackupPolicies configs={configs} isLoading={isLoadingConfigs} />
            </TabsContent>
          </Tabs>
        </div>

        <aside className="space-y-6 xl:sticky xl:top-6">
          <BackupStorageSummary records={records} destinations={destinations} configs={configs} />
        </aside>
      </div>
    </div>
  );
}
