import { useMutation, useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import type { ClusterDatabase, DatabaseAction } from './cluster-database-types';
export function ClusterDatabaseCard({
  projectId,
  record,
  onAction,
}: {
  projectId: string;
  record: ClusterDatabase;
  onAction: (action: DatabaseAction) => void;
}) {
  const [showBackups, setShowBackups] = useState(false);
  const endpoint = `/projects/${projectId}/clusters/${record.clusterId}/databases/${record.id}`;
  const credentials = useMutation({
    mutationFn: () =>
      apiClient.get<
        BaseResponse<{ username: string; password: string; host: string; port: number }>
      >(`${endpoint}/credentials`),
  });
  const backups = useQuery({
    queryKey: ['cluster-database-backups', record.id],
    queryFn: () =>
      apiClient.get<
        BaseResponse<
          { name: string; phase: string; startedAt: string; stoppedAt: string; error: string }[]
        >
      >(`${endpoint}/backups`),
    enabled: showBackups,
    refetchInterval: showBackups ? 5000 : false,
  });
  return (
    <article className="space-y-2 rounded border p-3">
      <div className="flex justify-between">
        <strong>{record.spec.name}</strong>
        <span>{record.status}</span>
      </div>
      <p className="text-sm">
        {record.spec.engine} /{' '}
        {record.spec.instances === 1 ? 'Standalone' : '3 replicated instances'} /{' '}
        {record.spec.storageGiB} GiB per instance
      </p>
      {record.error && (
        <p role="alert" className="text-destructive">
          {record.error}
        </p>
      )}
      <div className="flex flex-wrap gap-2">
        <Button variant="outline" onClick={() => onAction('recover')}>
          Review recovery
        </Button>
        <Button
          variant="outline"
          disabled={record.status !== 'READY' || credentials.isPending}
          onClick={() => credentials.mutate()}
        >
          Show connection
        </Button>
        {record.spec.engine === 'postgres' && record.spec.s3DestinationId && (
          <>
            <Button variant="outline" onClick={() => onAction('backup')}>
              Review backup
            </Button>
            <Button variant="outline" onClick={() => onAction('restore')}>
              Restore into new database
            </Button>
            <Button variant="ghost" onClick={() => setShowBackups(!showBackups)}>
              Backup history
            </Button>
          </>
        )}
      </div>
      {credentials.error && <p role="alert">{credentials.error.message}</p>}
      {credentials.data && (
        <div className="space-y-1 rounded bg-muted p-3">
          <p className="text-sm">Private cluster DNS; use from workloads in this cluster.</p>
          <pre className="overflow-auto whitespace-pre-wrap text-xs">
            {JSON.stringify(credentials.data.data, null, 2)}
          </pre>
          <Button variant="ghost" onClick={() => credentials.reset()}>
            Hide credentials
          </Button>
        </div>
      )}
      {showBackups && (
        <div>
          {backups.error && <p role="alert">{backups.error.message}</p>}
          {backups.data?.data.map((backup) => (
            <p key={backup.name} className="text-xs">
              {backup.name} / {backup.phase} / {backup.stoppedAt ?? backup.startedAt} {backup.error}
            </p>
          ))}
        </div>
      )}
    </article>
  );
}
