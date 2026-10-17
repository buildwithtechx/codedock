import type { HttpClient } from '../http';
import type { ActionRun, ActionRunner, ActionWorkflow } from '../types';

export function createActionsOperations(http: HttpClient) {
  return {
    listWorkflows: (): Promise<ActionWorkflow[]> =>
      http.get<ActionWorkflow[]>('/actions/workflows'),
    createWorkflow: (input: {
      name: string;
      description?: string;
      trigger?: string;
      yaml?: string;
    }): Promise<ActionWorkflow> => http.post<ActionWorkflow>('/actions/workflows', input),
    listRuns: (): Promise<ActionRun[]> => http.get<ActionRun[]>('/actions/runs'),
    getRun: (id: string): Promise<ActionRun> =>
      http.get<ActionRun>(`/actions/runs/${encodeURIComponent(id)}`),
    retryRun: (id: string): Promise<ActionRun> =>
      http.post<ActionRun>(`/actions/runs/${encodeURIComponent(id)}/retry`),
    cancelRun: (id: string): Promise<ActionRun> =>
      http.post<ActionRun>(`/actions/runs/${encodeURIComponent(id)}/cancel`),
    listRunners: (): Promise<ActionRunner[]> => http.get<ActionRunner[]>('/actions/runners'),
    createRunner: (input: {
      name: string;
      mode?: string;
      image?: string;
      labels?: string[];
      enabled: boolean;
    }): Promise<ActionRunner> => http.post<ActionRunner>('/actions/runners', input),
  };
}
