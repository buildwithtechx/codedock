import { useMutation, useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import type { WorkloadObservation } from './runtime-types';

export function RuntimeStatus({ serviceId }: { serviceId: string }) {
  const [command, setCommand] = useState('');
  const [showLogs, setShowLogs] = useState(false);
  const status = useQuery({
    queryKey: ['runtime-observation', serviceId],
    queryFn: () =>
      apiClient.get<BaseResponse<WorkloadObservation>>(`/apps/${serviceId}/runtime/observe`),
    refetchInterval: 5000,
  });
  const logs = useQuery({
    queryKey: ['runtime-logs', serviceId],
    queryFn: () => apiClient.get<BaseResponse<string>>(`/apps/${serviceId}/runtime/logs`),
    enabled: showLogs,
    refetchInterval: 5000,
  });
  const exec = useMutation({
    mutationFn: () => {
      const args: unknown = JSON.parse(command);
      if (!Array.isArray(args) || args.some((arg) => typeof arg !== 'string'))
        throw new Error('Enter a JSON array of command arguments');
      return apiClient.post<BaseResponse<string>>(`/apps/${serviceId}/runtime/exec`, {
        command: args,
      });
    },
  });
  const observed = status.data?.data;
  const native = observed?.kind === 'bare';
  return (
    <div className="space-y-3 border-t pt-4">
      <h3 className="font-medium">Observed workload</h3>
      {status.error && <p role="alert">{status.error.message}</p>}
      {observed && (
        <>
          <p>
            {observed.status} · {observed.available}/{observed.desired} replicas available
          </p>
          {observed.metricsError && <p role="status">{observed.metricsError}</p>}
          <ul className="space-y-2">
            {observed.pods.map((pod) => (
              <li key={pod.name} className="text-sm">
                {pod.name} · {pod.node} · {pod.phase} · {pod.restarts} restarts
                {observed.metricsAvailable && (
                  <>
                    {' '}
                    · {pod.cpu.toFixed(3)} CPU cores · {(pod.memoryBytes / 1024 / 1024).toFixed(1)}{' '}
                    MiB
                  </>
                )}
              </li>
            ))}
          </ul>
        </>
      )}
      <Button variant="outline" onClick={() => setShowLogs(!showLogs)}>
        {showLogs ? 'Hide logs' : native ? 'Show service logs' : 'Show pod logs'}
      </Button>
      {logs.error && <p role="alert">{logs.error.message}</p>}
      {showLogs && (
        <pre className="max-h-64 overflow-auto whitespace-pre-wrap text-xs">{logs.data?.data}</pre>
      )}
      <Input
        aria-label="Command arguments as JSON"
        placeholder={'["ls", "-la", "/app"]'}
        value={command}
        onChange={(event) => setCommand(event.target.value)}
      />
      <Button variant="outline" disabled={exec.isPending || !command} onClick={() => exec.mutate()}>
        {native ? 'Run as service user' : 'Run in a ready pod'}
      </Button>
      {exec.error && <p role="alert">{exec.error.message}</p>}
      {exec.data && (
        <pre className="max-h-64 overflow-auto whitespace-pre-wrap text-xs">{exec.data.data}</pre>
      )}
    </div>
  );
}
