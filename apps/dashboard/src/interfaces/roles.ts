export type ServiceRoleType = 'api' | 'web' | 'worker' | 'cron' | 'database';

export interface ServiceRoleAssignment {
  serviceId: string;
  role: ServiceRoleType;
  primaryPort?: number;
  healthEndpoint?: string;
}

export interface RoleHandoverRecord {
  id: string;
  previousServiceId: string;
  targetServiceId: string;
  role: ServiceRoleType;
  initiatedAt: string;
  completedAt?: string;
  status: 'pending' | 'completed' | 'rolled_back';
}

export interface RuntimeRoleContainer {
  containerId: string | null;
  replica: number;
  status: 'running' | 'starting' | 'stopped' | 'failed' | 'restarting' | 'missing' | 'unknown';
  restartCount: number | null;
}

export interface RuntimeRole {
  name: string;
  kind: string;
  command: string;
  replicas: number;
  containers: RuntimeRoleContainer[];
}

export interface ProjectRuntimeRoles {
  deploymentId: string;
  roles: RuntimeRole[];
}

export interface RuntimeRoleConfiguration {
  supported: boolean;
  deploymentId?: string;
  updatedAt?: string;
  roles: Array<{
    name: string;
    kind: string;
    command: string;
    replicas: number;
  }>;
}
