import { DatabaseBackup, Loader2, Play, Settings2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  useCreate,
  useDeleteRecord,
  useList,
  useListAllRecords,
  useListS3Destinations,
  useRestore,
  useTriggerDatabaseBackup,
  useUpdate,
} from '#/features/backups';
import type { BackupRecord } from '#/features/backups/interfaces';
import type { Database } from '#/features/databases/interfaces';
import { DatabaseBackupScheduleForm } from './database-backup-schedule-form';
import { DatabaseBackupTable } from './database-backup-table';
import { DatabaseRestoreDialog } from './database-restore-dialog';

export function BackupManager({ database }: { database: Database }) {
  const { data: configsData, isLoading: isLoadingConfigs, refetch: refetchConfigs } = useList();
  const {
    data: recordsData,
    isLoading: isLoadingRecords,
    refetch: refetchRecords,
  } = useListAllRecords(100);
  const { data: s3Data } = useListS3Destinations();

  const triggerDatabaseBackup = useTriggerDatabaseBackup();
  const createConfig = useCreate();
  const updateConfig = useUpdate();
  const restoreMutation = useRestore();
  const deleteRecordMutation = useDeleteRecord();

  const configs = configsData?.data || [];
  const allRecords = recordsData?.data || [];
  const destinations = s3Data?.data || [];

  const existingConfig = configs.find((c) => c.databaseId === database.id);
  const databaseRecords = allRecords.filter(
    (r) =>
      r.databaseId === database.id || (existingConfig && r.backupConfigId === existingConfig.id)
  );

  const [isConfiguring, setIsConfiguring] = useState(false);
  const [schedule, setSchedule] = useState(existingConfig?.schedule || '0 3 * * *');
  const [retentionDays, setRetentionDays] = useState(String(existingConfig?.retentionDays ?? 7));
  const [s3DestinationId, setS3DestinationId] = useState(
    existingConfig?.s3DestinationId || 'local'
  );
  const [selectedRecordForRestore, setSelectedRecordForRestore] = useState<BackupRecord | null>(
    null
  );

  useEffect(() => {
    if (existingConfig) {
      setSchedule(existingConfig.schedule || '0 3 * * *');
      setRetentionDays(String(existingConfig.retentionDays ?? 7));
      setS3DestinationId(existingConfig.s3DestinationId || 'local');
    }
  }, [existingConfig]);

  const handleBackupNow = async () => {
    try {
      await triggerDatabaseBackup.mutateAsync(database.id);
      toast.success('Backup triggered successfully');
      void refetchRecords();
      void refetchConfigs();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to trigger backup');
    }
  };

  const handleSaveSchedule = async (e: React.FormEvent) => {
    e.preventDefault();
    const isS3 = s3DestinationId !== 'local' && Boolean(s3DestinationId);
    try {
      if (existingConfig) {
        await updateConfig.mutateAsync({
          id: existingConfig.id,
          payload: {
            ...existingConfig,
            schedule,
            retentionDays: parseInt(retentionDays, 10) || 7,
            s3Enabled: isS3,
            s3DestinationId: isS3 ? s3DestinationId : undefined,
            disableLocal: isS3,
          },
        });
        toast.success('Backup schedule updated');
      } else {
        await createConfig.mutateAsync({
          payload: {
            projectId: database.projectId || 'global',
            databaseId: database.id,
            name: `${database.name}-backup`,
            description: `Automated backups for ${database.name}`,
            dbUser: database.username || 'postgres',
            backupEnabled: true,
            s3Enabled: isS3,
            disableLocal: isS3,
            s3DestinationId: isS3 ? s3DestinationId : undefined,
            schedule,
            timezone: 'UTC',
            timeout: 3600,
            retentionDays: parseInt(retentionDays, 10) || 7,
            maxBackups: 0,
            maxStorageGb: 0,
          },
        });
        toast.success('Backup schedule configured');
      }
      setIsConfiguring(false);
      void refetchConfigs();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to save schedule');
    }
  };

  const handleConfirmRestore = async (recordId: string) => {
    try {
      await restoreMutation.mutateAsync({ id: recordId });
      toast.success('Database restore triggered successfully');
      setSelectedRecordForRestore(null);
      void refetchRecords();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to restore database');
    }
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

  if (isLoadingConfigs || isLoadingRecords) {
    return (
      <div className="flex h-48 items-center justify-center rounded-2xl border border-border/60 bg-card">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="flex items-center gap-2 font-semibold text-foreground text-lg tracking-tight">
            <DatabaseBackup className="h-5 w-5 text-primary" />
            Backups & Disaster Recovery
          </h2>
          <p className="mt-0.5 text-muted-foreground text-xs">
            Manage scheduled snapshots and one-click data restoration for this database.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsConfiguring(!isConfiguring)}
            className="gap-1.5"
          >
            <Settings2 className="h-4 w-4" />
            {isConfiguring ? 'Close settings' : 'Configure schedule'}
          </Button>

          <Button
            size="sm"
            onClick={handleBackupNow}
            disabled={triggerDatabaseBackup.isPending}
            className="gap-1.5"
          >
            {triggerDatabaseBackup.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Play className="h-4 w-4 fill-current" />
            )}
            Back up now
          </Button>
        </div>
      </div>

      {isConfiguring && (
        <DatabaseBackupScheduleForm
          s3DestinationId={s3DestinationId}
          setS3DestinationId={setS3DestinationId}
          schedule={schedule}
          setSchedule={setSchedule}
          retentionDays={retentionDays}
          setRetentionDays={setRetentionDays}
          destinations={destinations}
          onSave={handleSaveSchedule}
          onCancel={() => setIsConfiguring(false)}
          isPending={createConfig.isPending || updateConfig.isPending}
        />
      )}

      <DatabaseBackupTable
        records={databaseRecords}
        onRestore={(record) => setSelectedRecordForRestore(record)}
        onDelete={handleDeleteRecord}
        onBackupNow={handleBackupNow}
        restorePending={restoreMutation.isPending}
        deletePending={deleteRecordMutation.isPending}
        backupPending={triggerDatabaseBackup.isPending}
      />

      <DatabaseRestoreDialog
        record={selectedRecordForRestore}
        databaseName={database.name}
        onClose={() => setSelectedRecordForRestore(null)}
        onConfirm={handleConfirmRestore}
        isPending={restoreMutation.isPending}
      />
    </div>
  );
}
