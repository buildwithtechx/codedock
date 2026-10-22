/** Wire contracts for Docker Swarm cluster connectivity and inspection. */

export interface ConnectSwarmInput {
  name: string;
  requestId: string;
}

export interface SwarmInspectionInput {
  details?: boolean;
}

export type SwarmState = 'inactive' | 'active' | 'pending' | 'locked' | 'error' | 'unknown';
export type SwarmNodeRole = 'manager' | 'worker' | 'unknown';
export type SwarmServiceMode = 'replicated' | 'global' | 'job' | 'unknown';

export interface SwarmNodeSummary {
  id: string;
  name: string;
  role: SwarmNodeRole;
  availability: string;
  state: string;
  leader: boolean;
}

export interface SwarmPortMapping {
  published: number;
  target: number;
  protocol: string;
  mode: string;
}

export interface SwarmVolumeMount {
  type: string;
  target: string;
  readOnly: boolean;
}

export interface SwarmServiceSummary {
  id: string;
  name: string;
  image: string;
  mode: SwarmServiceMode;
  desiredTasks: number | null;
  runningTasks: number | null;
  updateState: string | null;
  ports: SwarmPortMapping[];
  mounts: SwarmVolumeMount[];
}

export interface SwarmInspection {
  state: SwarmState;
  role: 'manager' | 'worker' | null;
  nodeId: string | null;
  clusterId: string | null;
  nodeCount: number | null;
  managerCount: number | null;
  observedAt: string;
  inventory: {
    nodes: SwarmNodeSummary[];
    stacks: Array<{ name: string | null; services: SwarmServiceSummary[] }>;
  } | null;
}
