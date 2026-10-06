import { useMutation, useQuery } from '@tanstack/react-query';
import { Button } from '#/components/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';

const finished = new Set(['READY', 'FAILED', 'CANCELLED']);
const phases = [
  'PREPARING',
  'CLONING',
  'PULLING',
  'BUILDING',
  'STARTING',
  'READINESS',
  'ROUTING',
  'READY',
];

export function DeploymentProgress({ deploymentId }: { deploymentId: string }) {
  const progress = useQuery({
    queryKey: ['deployment-progress', deploymentId],
    queryFn: () =>
      apiClient.get<BaseResponse<{ status: string; buildLogs: string }>>(
        `/deployments/${deploymentId}/logs`
      ),
    refetchInterval: (query) => (finished.has(query.state.data?.data.status ?? '') ? false : 1500),
  });
  const cancel = useMutation({
    mutationFn: () => apiClient.post(`/deployments/${deploymentId}/cancel`),
    onSuccess: () => progress.refetch(),
  });
  const status = progress.data?.data.status ?? 'PREPARING';
  return (
    <div className="space-y-4">
      <p aria-live="polite">Deployment: {status}</p>
      <div className="flex flex-wrap gap-2">
        {phases.map((phase) => (
          <span
            key={phase}
            className={`rounded-md border px-2 py-1 text-xs ${status === phase ? 'bg-primary text-primary-foreground' : ''}`}
          >
            {phase.toLowerCase()}
          </span>
        ))}
      </div>
      {progress.isError && (
        <p role="alert" className="text-destructive">
          Progress unavailable.{' '}
          <Button variant="outline" onClick={() => progress.refetch()}>
            Retry
          </Button>
        </p>
      )}
      <pre className="max-h-80 overflow-auto rounded-md bg-muted p-3 text-xs">
        {progress.data?.data.buildLogs ?? 'Waiting for deployment…'}
      </pre>
      {!finished.has(status) && (
        <Button
          variant="destructive"
          disabled={cancel.isPending || cancel.isSuccess}
          onClick={() => cancel.mutate()}
        >
          {cancel.isSuccess ? 'Cancellation requested' : 'Cancel deployment'}
        </Button>
      )}
      {cancel.isError && (
        <p role="alert" className="text-destructive">
          {cancel.error.message}
        </p>
      )}
    </div>
  );
}
