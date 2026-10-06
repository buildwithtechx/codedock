import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from '#/components/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';

type Operation = {
  id: string;
  kind: string;
  status: string;
  phase: string;
  logs: string;
  error: string;
};
export function BackupOperationHistory() {
  const client = useQueryClient();
  const query = useQuery({
    queryKey: ['backup-operations'],
    queryFn: () => apiClient.get<BaseResponse<Operation[]>>('/operations'),
    refetchInterval: 3000,
  });
  const cancel = useMutation({
    mutationFn: (id: string) => apiClient.post(`/operations/${id}/cancel`),
    onSettled: () => client.invalidateQueries({ queryKey: ['backup-operations'] }),
  });
  return (
    <section className="space-y-3">
      <h2 className="font-semibold">Backup and restore runs</h2>
      {query.error && <p role="alert">{query.error.message}</p>}
      {cancel.error && <p role="alert">{cancel.error.message}</p>}
      {query.data?.data
        .filter((operation) => operation.kind === 'backup' || operation.kind === 'restore')
        .map((operation) => (
          <details key={operation.id} className="rounded border p-3">
            <summary>
              {operation.kind} · {operation.status} · {operation.phase}
            </summary>
            {operation.error && (
              <p role="alert" className="text-destructive">
                {operation.error}
              </p>
            )}
            <pre className="mt-2 max-h-48 overflow-auto whitespace-pre-wrap text-xs">
              {operation.logs}
            </pre>
            {operation.status === 'RUNNING' && (
              <Button
                size="sm"
                variant="destructive"
                disabled={cancel.isPending}
                onClick={() => cancel.mutate(operation.id)}
              >
                Interrupt run
              </Button>
            )}
          </details>
        ))}
    </section>
  );
}
