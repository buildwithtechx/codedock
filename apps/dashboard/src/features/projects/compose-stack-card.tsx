import { useMutation } from '@tanstack/react-query';
import { Button } from '#/components/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { activeStackStates, type ComposeStack, stackResults } from './compose-stack-types';

export function ComposeStackCard({
  stack,
  onRefresh,
  onEdit,
}: {
  stack: ComposeStack;
  onRefresh: () => void;
  onEdit: (stack: ComposeStack, content: string) => void;
}) {
  const operation = useMutation({
    mutationFn: (action: 'deploy' | 'cancel') =>
      apiClient.post(`/projects/${stack.projectId}/stacks/${stack.id}/${action}`),
    onSuccess: onRefresh,
  });
  const load = useMutation({
    mutationFn: () =>
      apiClient.get<BaseResponse<{ stack: ComposeStack; content: string }>>(
        `/projects/${stack.projectId}/stacks/${stack.id}/config`
      ),
    onSuccess: (response) => onEdit(response.data.stack, response.data.content),
  });
  return (
    <article className="space-y-3 rounded-xl border p-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium">{stack.name}</h3>
        <span className="text-sm">{stack.status}</span>
      </div>
      <p className="text-muted-foreground text-xs">
        Revision {stack.revision}. Named volumes remain attached across redeployments.
      </p>
      {stack.error && (
        <p role="alert" className="text-destructive text-sm">
          {stack.error}
        </p>
      )}
      <ul className="space-y-1 text-sm">
        {stackResults(stack).map((service) => (
          <li key={service.containerId}>
            {service.name}: {service.state}; health {service.health}; exit {service.exitCode}
          </li>
        ))}
      </ul>
      <Button
        variant={activeStackStates.has(stack.status) ? 'destructive' : 'outline'}
        disabled={
          operation.isPending ||
          (operation.isSuccess &&
            operation.variables === 'cancel' &&
            activeStackStates.has(stack.status))
        }
        onClick={() => operation.mutate(activeStackStates.has(stack.status) ? 'cancel' : 'deploy')}
      >
        {activeStackStates.has(stack.status) ? 'Cancel deployment' : 'Deploy saved configuration'}
      </Button>
      <Button
        variant="outline"
        disabled={load.isPending || activeStackStates.has(stack.status)}
        onClick={() => load.mutate()}
      >
        Edit saved configuration
      </Button>
      {load.isError && (
        <p role="alert" className="text-destructive text-sm">
          {load.error.message}
        </p>
      )}
      {operation.isError && (
        <p role="alert" className="text-destructive text-sm">
          {operation.error.message}
        </p>
      )}
    </article>
  );
}
