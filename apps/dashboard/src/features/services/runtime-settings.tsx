import { useMutation, useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { Cluster } from '#/features/servers/cluster-types';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { RuntimeStatus } from './runtime-status';
import type { RuntimeReview, RuntimeTarget, ServiceRuntime } from './runtime-types';
import { RuntimeVolumeFields } from './runtime-volume-fields';

export function RuntimeSettings({
  serviceId,
  projectId,
}: {
  serviceId: string;
  projectId: string;
}) {
  const runtime = useQuery({
    queryKey: ['service-runtime', serviceId],
    queryFn: () => apiClient.get<BaseResponse<ServiceRuntime>>(`/apps/${serviceId}/runtime`),
    refetchInterval: 3000,
  });
  const clusters = useQuery({
    queryKey: ['runtime-clusters', projectId],
    queryFn: () =>
      apiClient.get<BaseResponse<Cluster[]>>(`/projects/${projectId}/runtime-clusters`),
  });
  const [draft, setDraft] = useState<RuntimeTarget>();
  const registries = useQuery({
    queryKey: ['runtime-registries', projectId],
    queryFn: () =>
      apiClient.get<BaseResponse<{ id: string; registryUrl: string }[]>>(
        `/projects/${projectId}/registries`
      ),
  });
  const [review, setReview] = useState<RuntimeReview>();
  const target = draft ?? runtime.data?.data.target ?? { kind: 'docker', nodeIds: [], volumes: [] };
  const update = (changes: Partial<RuntimeTarget>) => {
    setDraft({ ...target, ...changes });
    setReview(undefined);
  };
  const prepare = useMutation({
    mutationFn: () =>
      apiClient.post<BaseResponse<RuntimeReview>>(`/apps/${serviceId}/runtime/review`, {
        target,
        revision: runtime.data?.data.revision,
      }),
    onSuccess: (result) => setReview(result.data),
  });
  const operation = useQuery({
    queryKey: ['runtime-operation', serviceId, review?.operation.id],
    queryFn: () =>
      apiClient.get<BaseResponse<{ status: string; error: string }>>(
        `/apps/${serviceId}/runtime/operation?operationId=${review?.operation.id}`
      ),
    enabled: !!review,
    refetchInterval: 1500,
  });
  const current = operation.data?.data;
  const apply = useMutation({
    mutationFn: () =>
      apiClient.post(`/apps/${serviceId}/runtime/apply`, {
        operationId: review?.operation.id,
        confirmation: review?.confirmation,
      }),
    onSuccess: () => {
      void runtime.refetch();
      void operation.refetch();
    },
  });
  const error =
    runtime.error ??
    clusters.error ??
    registries.error ??
    prepare.error ??
    apply.error ??
    operation.error;
  const selected = clusters.data?.data.find((cluster) => cluster.id === target.clusterId);
  return (
    <section className="space-y-4 rounded-lg border p-5">
      <h2 className="font-semibold">Deployment destination</h2>
      <p className="text-muted-foreground text-sm">
        Docker replicas run on one Docker host. Kubernetes replicas are scheduled across the
        selected cluster nodes. Neither setting creates database replication; configure database
        replication separately.
      </p>
      <p className="text-muted-foreground text-sm">
        Stop an existing workload before changing its destination. Persistent data stays at its
        current destination until you migrate it.
      </p>
      {error && (
        <p role="alert" className="text-destructive">
          {error.message}
        </p>
      )}
      {runtime.data?.data.error && <p role="alert">{runtime.data.data.error}</p>}
      <Label htmlFor="runtime-kind">Runtime</Label>
      <select
        id="runtime-kind"
        className="w-full rounded border bg-background p-2"
        value={target.kind}
        onChange={(event) =>
          update({
            kind: event.target.value as RuntimeTarget['kind'],
            clusterId: undefined,
            nodeIds: [],
            volumes: [],
          })
        }
      >
        <option value="docker">Docker destination</option>
        <option value="kubernetes">Kubernetes cluster</option>
      </select>
      {target.kind === 'kubernetes' && (
        <>
          <Label htmlFor="runtime-cluster">Ready cluster</Label>
          <select
            id="runtime-cluster"
            className="w-full rounded border bg-background p-2"
            value={target.clusterId ?? ''}
            onChange={(event) => update({ clusterId: event.target.value, nodeIds: [] })}
          >
            <option value="">Select a cluster</option>
            {clusters.data?.data
              .filter((cluster) => cluster.status === 'READY')
              .map((cluster) => (
                <option key={cluster.id} value={cluster.id}>
                  {cluster.name}
                </option>
              ))}
          </select>
          <Label htmlFor="runtime-image">Image repository for Git builds</Label>
          <select
            aria-label="Project image registry"
            className="w-full rounded border bg-background p-2"
            value={target.registryId ?? ''}
            onChange={(event) => update({ registryId: event.target.value })}
          >
            <option value="">Use application's existing registry</option>
            {registries.data?.data.map((registry) => (
              <option key={registry.id} value={registry.id}>
                {registry.registryUrl}
              </option>
            ))}
          </select>
          <Input
            id="runtime-image"
            placeholder="registry.example.com/team/application"
            value={target.imageRepository ?? ''}
            onChange={(event) => update({ imageRepository: event.target.value })}
          />
          <p className="text-muted-foreground text-sm">
            Git builds publish to the selected registry before deployment. Image deployments use
            their existing image reference.
          </p>
          <fieldset className="space-y-2">
            <legend className="font-medium text-sm">
              Placement (leave empty to allow every cluster node)
            </legend>
            {selected?.nodes.map((node) => (
              <label className="flex gap-2 text-sm" key={node.serverId}>
                <input
                  type="checkbox"
                  checked={(target.nodeIds ?? []).includes(node.serverId)}
                  onChange={(event) =>
                    update({
                      nodeIds: event.target.checked
                        ? [...(target.nodeIds ?? []), node.serverId]
                        : (target.nodeIds ?? []).filter((id) => id !== node.serverId),
                    })
                  }
                />
                {node.privateIp}
              </label>
            ))}
          </fieldset>
          <RuntimeVolumeFields
            volumes={target.volumes ?? []}
            onChange={(volumes) => update({ volumes })}
          />
        </>
      )}
      <Button
        disabled={!runtime.data || prepare.isPending || apply.isPending}
        onClick={() => prepare.mutate()}
      >
        Review destination
      </Button>
      {review && (
        <div className="space-y-3 rounded border p-3">
          <p>{review.operation.effects}</p>
          {current && <p role="status">{current.status}</p>}
          {current?.error && <p role="alert">{current.error}</p>}
          <Button
            disabled={
              apply.isPending ||
              (current?.status !== undefined && current.status !== 'REVIEWED') ||
              Date.now() / 1000 >= review.operation.expiresAt
            }
            onClick={() => apply.mutate()}
          >
            Apply reviewed destination
          </Button>
        </div>
      )}
      {runtime.data?.data.target.kind === 'kubernetes' && <RuntimeStatus serviceId={serviceId} />}
    </section>
  );
}
