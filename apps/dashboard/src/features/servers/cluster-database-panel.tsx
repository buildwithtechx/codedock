import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { ClusterDatabaseCard } from './cluster-database-card';
import { ClusterDatabaseDialog } from './cluster-database-dialog';
import type { ClusterDatabase, DatabaseAction } from './cluster-database-types';
import type { Cluster } from './cluster-types';
export function ClusterDatabasePanel({
  projectId,
  cluster,
}: {
  projectId: string;
  cluster: Cluster;
}) {
  const records = useQuery({
    queryKey: ['cluster-databases', cluster.id],
    queryFn: () =>
      apiClient.get<BaseResponse<ClusterDatabase[]>>(
        `/projects/${projectId}/clusters/${cluster.id}/databases`
      ),
    refetchInterval: 5000,
  });
  const [selection, setSelection] = useState<{
    action: DatabaseAction;
    record?: ClusterDatabase;
  }>();
  return (
    <section className="space-y-3 border-t pt-3">
      <h4 className="font-medium">Cluster databases</h4>
      <p className="text-muted-foreground text-sm">
        Database instances provide engine replication, independently of application replicas.
        Replicated modes need three cluster nodes. PostgreSQL operators manage promotion and
        object-store recovery. Redis uses a fixed primary, persistent AOF and validated replicas;
        Sentinel failover is not enabled.
      </p>
      <div className="flex gap-2">
        <Button variant="outline" onClick={() => setSelection({ action: 'operators' })}>
          Prepare PostgreSQL operators
        </Button>
        <Button onClick={() => setSelection({ action: 'create' })}>Create database</Button>
      </div>
      {records.error && <p role="alert">{records.error.message}</p>}
      {records.data?.data.map((record) => (
        <ClusterDatabaseCard
          key={record.id}
          projectId={projectId}
          record={record}
          onAction={(action) => setSelection({ action, record })}
        />
      ))}
      {selection && (
        <ClusterDatabaseDialog
          projectId={projectId}
          cluster={cluster}
          {...selection}
          onClose={() => {
            setSelection(undefined);
            void records.refetch();
          }}
        />
      )}
    </section>
  );
}
