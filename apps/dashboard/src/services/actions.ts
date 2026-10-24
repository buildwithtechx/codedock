import type {
  ActionCapabilities,
  ActionRun,
  ActionRunner,
  ActionWorkflow,
  CreateRunnerInput,
  CreateWorkflowInput,
} from '#/interfaces/actions';
import { apiClient } from '#/lib/api-client';

export const actionsService = {
  listWorkflows: async (): Promise<ActionWorkflow[]> => {
    const res = await apiClient.get<{ data: ActionWorkflow[] }>('/actions/workflows');
    return res.data ?? [];
  },

  createWorkflow: async (input: CreateWorkflowInput): Promise<ActionWorkflow> => {
    const res = await apiClient.post<{ data: ActionWorkflow }>('/actions/workflows', input);
    return res.data;
  },

  triggerWorkflow: async (id: string): Promise<ActionRun> => {
    const res = await apiClient.post<{ data: ActionRun }>(
      `/actions/workflows/${encodeURIComponent(id)}/trigger`,
      {}
    );
    return res.data;
  },

  listRuns: async (): Promise<ActionRun[]> => {
    const res = await apiClient.get<{ data: ActionRun[] }>('/actions/runs');
    return res.data ?? [];
  },

  getRun: async (id: string): Promise<ActionRun> => {
    const res = await apiClient.get<{ data: ActionRun }>(`/actions/runs/${encodeURIComponent(id)}`);
    return res.data;
  },

  retryRun: async (id: string): Promise<ActionRun> => {
    const res = await apiClient.post<{ data: ActionRun }>(
      `/actions/runs/${encodeURIComponent(id)}/retry`,
      {}
    );
    return res.data;
  },

  cancelRun: async (id: string): Promise<ActionRun> => {
    const res = await apiClient.post<{ data: ActionRun }>(
      `/actions/runs/${encodeURIComponent(id)}/cancel`,
      {}
    );
    return res.data;
  },

  listRunners: async (): Promise<ActionRunner[]> => {
    const res = await apiClient.get<{ data: ActionRunner[] }>('/actions/runners');
    return res.data ?? [];
  },

  inspectDestination: async (serverId: string): Promise<ActionCapabilities> => {
    const res = await apiClient.post<{ data: ActionCapabilities }>('/actions/runners/inspect', {
      serverId,
    });
    return res.data;
  },

  createRunner: async (input: CreateRunnerInput): Promise<ActionRunner> => {
    const res = await apiClient.post<{ data: ActionRunner }>('/actions/runners', input);
    return res.data;
  },

  deleteWorkflow: async (id: string): Promise<void> => {
    await apiClient.delete(`/actions/workflows/${encodeURIComponent(id)}`);
  },

  deleteRunner: async (id: string): Promise<void> => {
    await apiClient.delete(`/actions/runners/${encodeURIComponent(id)}`);
  },

  deleteRun: async (id: string): Promise<void> => {
    await apiClient.delete(`/actions/runs/${encodeURIComponent(id)}`);
  },
};
