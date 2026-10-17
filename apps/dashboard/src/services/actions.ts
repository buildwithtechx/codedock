import { apiClient } from '#/lib/api-client';

export interface ActionWorkflow {
  id: string;
  projectId?: string;
  name: string;
  description?: string;
  trigger: string;
  yaml?: string;
  enabled: boolean;
  owner?: string;
  repo?: string;
  path?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ActionStep {
  id: string;
  name: string;
  status: string;
  durationSeconds?: number;
}

export interface ActionJob {
  id: string;
  name: string;
  status: string;
  steps: ActionStep[];
  runnerId?: string;
  durationSeconds?: number;
}

export interface ActionRun {
  id: string;
  workflowId: string;
  name: string;
  number: number;
  status: 'queued' | 'running' | 'ready' | 'failed' | 'cancelled';
  triggeredBy: string;
  runnerId?: string;
  logs?: string;
  jobs?: ActionJob[];
  attempt?: number;
  owner?: string;
  repo?: string;
  ref?: string;
  revision?: string;
  startedAt?: string;
  finishedAt?: string;
  durationSeconds: number;
  createdAt: string;
  updatedAt: string;
}

export interface ActionCapabilities {
  os: 'linux' | 'macos';
  architecture: 'x64' | 'arm64';
  docker: boolean;
  dockerArchitecture?: 'x64' | 'arm64';
  dockerPlatforms?: ('linux/amd64' | 'linux/arm64')[];
  git: boolean;
  node: boolean;
  distribution?: string | null;
  version?: string | null;
}

export interface ActionRunnerConfig {
  mode: 'container' | 'native';
  labels: string[];
  image: string | null;
  maxParallel: number;
  cpu: number;
  memoryMb: number;
  allowDockerSocket: boolean;
  cloudDiskGb?: number;
}

export interface ActionRunner {
  id: string;
  name: string;
  serverId?: string | null;
  cloudPoolId?: string | null;
  kind?: 'server' | 'cloud';
  mode: 'container' | 'native';
  image?: string;
  status: 'online' | 'offline' | 'busy';
  labels: string[];
  enabled: boolean;
  error?: string | null;
  config?: ActionRunnerConfig;
  capabilities?: ActionCapabilities | null;
  checkedAt?: string;
  lastHeartbeatAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateWorkflowInput {
  projectId?: string;
  name: string;
  description?: string;
  trigger?: string;
  yaml?: string;
  path?: string;
  ref?: string;
  owner?: string;
  repo?: string;
  yamlContent?: string;
}

export interface CreateRunnerInput {
  name: string;
  serverId?: string;
  cloudPoolId?: string;
  kind?: string;
  mode?: string;
  image?: string;
  labels?: string[];
  enabled: boolean;
  config?: ActionRunnerConfig;
}

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
