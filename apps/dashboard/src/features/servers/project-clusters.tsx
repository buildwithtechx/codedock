import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { ClusterReviewDialog } from './cluster-review-dialog';
import type { Cluster } from './cluster-types';

export function ProjectClusters({ projectId }: { projectId: string }) {
  const clusters = useQuery({
    queryKey: ['clusters', projectId],
    queryFn: () => apiClient.get<BaseResponse<Cluster[]>>(`/projects/${projectId}/clusters`),
    refetchInterval: 5000,
  });
  const [selected, setSelected] = useState<{
    cluster?: Cluster;
    action: 'install' | 'join' | 'upgrade' | 'remove' | 'recover';
  }>();
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-semibold">Compute clusters</h3>
        <Button onClick={() => setSelected({ action: 'install' })}>Prepare cluster</Button>
      </div>
      <p className="text-muted-foreground text-sm">
        Instance administrators can manage K3s infrastructure here. Select a ready cluster from an
        application's build settings to deploy it with cluster placement and persistent storage.
      </p>
      {clusters.isLoading && <p role="status">Loading clusters…</p>}
      {clusters.error && <p role="alert">{clusters.error.message}</p>}
      {clusters.data?.data.map((cluster) => (
        <div key={cluster.id} className="space-y-3 rounded-md border p-4">
          <div className="flex justify-between">
            <span>
              {cluster.name} · {cluster.version}
            </span>
            <span>{cluster.status}</span>
          </div>
          <p className="text-sm">{cluster.nodes.length} nodes · one control plane</p>
          {cluster.error && (
            <p role="alert" className="text-destructive">
              {cluster.error}
            </p>
          )}
          <div className="flex flex-wrap gap-2">
            {(['install', 'join', 'upgrade', 'remove', 'recover'] as const)
              .filter((action) =>
                action === 'recover'
                  ? ['FAILED', 'INTERRUPTED', 'RECOVERY_REQUIRED'].includes(cluster.status)
                  : action !== 'install' || cluster.status !== 'READY'
              )
              .map((action) => (
                <Button
                  key={action}
                  variant={action === 'remove' ? 'destructive' : 'outline'}
                  disabled={cluster.status === 'APPLYING' || cluster.status === 'REMOVED'}
                  onClick={() => setSelected({ cluster, action })}
                >
                  Review {action}
                </Button>
              ))}
          </div>
        </div>
      ))}
      {selected && (
        <ClusterReviewDialog
          projectId={projectId}
          {...selected}
          onClose={() => {
            setSelected(undefined);
            void clusters.refetch();
          }}
        />
      )}
    </div>
  );
}
