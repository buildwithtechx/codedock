export type WorkflowStatus = 'queued' | 'running' | 'ready' | 'failed' | 'cancelled';

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
  status: WorkflowStatus;
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
