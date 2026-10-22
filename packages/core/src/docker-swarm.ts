/** Docker Swarm multi-node orchestration and task ownership utilities. */

export interface SwarmNodeSummary {
  id: string;
  name: string;
  role: 'manager' | 'worker' | 'unknown';
  availability: string;
  state: string;
  leader: boolean;
}

export interface SwarmServiceSummary {
  id: string;
  name: string;
  image: string;
  mode: 'replicated' | 'global' | 'job' | 'unknown';
  desiredTasks: number | null;
  runningTasks: number | null;
  updateState: string | null;
  ports: Array<{ published: number; target: number; protocol: string; mode: string }>;
  mounts: Array<{ type: string; target: string; readOnly: boolean }>;
}

export interface SwarmInspection {
  state: 'inactive' | 'active' | 'pending' | 'locked' | 'error' | 'unknown';
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

/** Resolves Swarm task owner label if the container is managed by Docker Swarm. */
export function swarmTaskOwner(labels: Record<string, string> | null | undefined): string | null {
  if (!labels) return null;
  const service = labels['com.docker.swarm.service.name'] || labels['com.docker.swarm.service.id'];
  const task = labels['com.docker.swarm.task.id'];
  const stack = labels['com.docker.stack.namespace'];
  return service || task || stack || null;
}

/** Asserts that a container is not owned by Docker Swarm scheduler. */
export function assertStandaloneContainer(labels: Record<string, string> | null | undefined): void {
  const owner = swarmTaskOwner(labels);
  if (owner) {
    throw new Error(
      `This container belongs to Docker Swarm (${owner}). Manage its service through a Swarm manager; it cannot be adopted or stopped as a standalone container.`
    );
  }
}

/** Asserts that a configuration does not request features unsupported by Docker Swarm. */
export function assertSwarmSupportedSettings(config: { swapDisabled?: boolean }): void {
  if (config.swapDisabled) {
    throw new Error(
      'Docker Swarm cannot enforce per-service swap settings. Remove Disable swap or use a standalone Docker target.'
    );
  }
}
