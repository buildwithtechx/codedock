import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { Node } from '@xyflow/react';
import { Button } from '#/components/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';

export function CanvasResourcePanel({
  node,
  onClose,
  environmentId,
}: {
  node: Node;
  onClose: () => void;
  environmentId: string;
}) {
  const client = useQueryClient();
  const serviceId = typeof node.data.serviceId === 'string' ? node.data.serviceId : '';
  const databaseId = typeof node.data.databaseId === 'string' ? node.data.databaseId : '';
  const path = serviceId ? `/apps/${serviceId}` : databaseId ? `/databases/${databaseId}` : '';
  const details = useQuery({
    queryKey: ['canvas-resource', node.id],
    queryFn: () => apiClient.get<BaseResponse<Record<string, unknown>>>(path),
    enabled: !!path,
  });
  const metrics = useQuery({
    queryKey: ['canvas-observation', environmentId, node.id],
    queryFn: () =>
      apiClient.get<BaseResponse<{ logs: string; metrics: Record<string, unknown> }>>(
        `/environments/${environmentId}/canvas/observe?node=${encodeURIComponent(node.id)}`
      ),
    enabled: !!path || !!node.data.stackId,
    refetchInterval: 5000,
  });
  const action = useMutation({
    mutationFn: (operation: string) => apiClient.post(`${path}/${operation}`),
    onSuccess: async () => {
      await details.refetch();
      await client.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
  return (
    <aside className="w-80 shrink-0 space-y-4 overflow-y-auto border-l bg-card p-4">
      <div className="flex items-center justify-between">
        <h2 className="font-medium">{String(node.data.name ?? node.data.label ?? 'Resource')}</h2>
        <Button size="sm" variant="ghost" onClick={onClose}>
          Close
        </Button>
      </div>
      <p aria-live="polite" className="text-sm">
        Runtime: {String(node.data.status ?? 'unobserved')}
      </p>
      {serviceId && (
        <nav className="flex flex-wrap gap-3 text-primary text-sm">
          <a href={`/services/${serviceId}`}>Details and logs</a>
          <a href={`/services/${serviceId}/metrics`}>Metrics</a>
          <a href={`/services/${serviceId}/variables`}>Variables</a>
          <a href={`/services/${serviceId}/domains`}>Routing</a>
        </nav>
      )}
      {databaseId && (
        <a className="text-primary text-sm" href={`/services/${databaseId}`}>
          Database details
        </a>
      )}
      {!!node.data.stackId && (
        <a
          className="text-primary text-sm"
          href={`/projects/${String(node.data.projectId)}/compose`}
        >
          Stack configuration and operations
        </a>
      )}
      {!!path && (
        <div className="flex flex-wrap gap-2">
          {(serviceId ? ['redeploy', 'restart', 'stop'] : ['start', 'restart', 'stop']).map(
            (operation) => (
              <Button
                key={operation}
                size="sm"
                variant="outline"
                disabled={action.isPending}
                onClick={() => action.mutate(operation)}
              >
                {operation}
              </Button>
            )
          )}
        </div>
      )}
      {(details.isError || metrics.isError || action.isError) && (
        <p role="alert" className="text-destructive text-sm">
          {action.error?.message ?? 'Runtime details unavailable. Retry or open the resource page.'}
        </p>
      )}
      {details.data && (
        <dl className="space-y-2 text-sm">
          {['name', 'imageRef', 'repositoryUrl', 'internalPort', 'domain', 'engine', 'version']
            .filter((key) => details.data.data[key])
            .map((key) => (
              <div key={key}>
                <dt className="text-muted-foreground">{key}</dt>
                <dd className="break-all">{String(details.data.data[key])}</dd>
              </div>
            ))}
        </dl>
      )}
      {metrics.data && (
        <pre className="overflow-auto rounded bg-muted p-2 text-xs">
          {JSON.stringify(metrics.data.data.metrics, null, 2)}
        </pre>
      )}
      {metrics.data && (
        <pre className="max-h-64 overflow-auto rounded bg-muted p-2 text-xs">
          {metrics.data.data.logs || 'No log entries yet.'}
        </pre>
      )}
    </aside>
  );
}
