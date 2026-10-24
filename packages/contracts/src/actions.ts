export interface ActionWorkflow {
  id: string;
  projectId?: string;
  name: string;
  description?: string;
  trigger: string;
  yaml?: string;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
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
  durationSeconds: number;
  createdAt: string;
  updatedAt: string;
}

export interface ActionRunner {
  id: string;
  name: string;
  mode: 'container' | 'native';
  image?: string;
  status: 'online' | 'offline' | 'busy';
  labels: string[];
  enabled: boolean;
  lastHeartbeatAt?: string;
  createdAt: string;
  updatedAt: string;
}
