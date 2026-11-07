export interface CodedockConnection {
  url: string;
  token: string;
  userName?: string;
  userEmail?: string;
  version?: string;
  instanceName?: string;
}

export interface CodedockMeta {
  name: string;
  version: string;
  capabilities: {
    sseLogs: boolean;
    push: boolean;
    actions: boolean;
    apps: boolean;
    [key: string]: boolean;
  };
}

export interface UserProfile {
  id: string;
  email: string;
  name?: string;
  role?: string;
}

export interface CodedockProject {
  id: string;
  name: string;
  description?: string;
  environments?: CodedockEnvironment[];
}

export interface CodedockEnvironment {
  id: string;
  name: string;
  projectId: string;
}

export interface CodedockApp {
  id: string;
  name: string;
  status?: 'running' | 'stopped' | 'failed' | 'building' | 'deploying' | string;
  projectId?: string;
  projectName?: string;
  environmentId?: string;
  environmentName?: string;
  domain?: string;
  gitRepo?: string;
  branch?: string;
  replicas?: number;
  cpuLimit?: string;
  memoryLimit?: string;
  updatedAt?: string;
}

export interface CodedockDatabase {
  id: string;
  name: string;
  type: 'postgres' | 'mysql' | 'redis' | 'mongodb' | 'mariadb' | string;
  status: 'running' | 'stopped' | 'failed' | string;
  projectId?: string;
  port?: number;
  internalHost?: string;
  externalPort?: number;
}

export interface CodedockServer {
  id: string;
  name: string;
  ipAddress: string;
  status: 'online' | 'offline' | 'unreachable' | string;
  sshPort?: number;
  dockerVersion?: string;
  osInfo?: string;
  isDefault?: boolean;
}

export interface SystemStats {
  cpu?: { percent: number; cores: number };
  memory?: { totalMb: number; usedMb: number; percent: number };
  disk?: { totalGb: number; usedGb: number; percent: number };
  uptimeSeconds?: number;
}

export interface CodedockDeployment {
  id: string;
  serviceId?: string;
  serviceName?: string;
  projectName?: string;
  status: 'pending' | 'running' | 'succeeded' | 'failed' | 'cancelled' | string;
  commitHash?: string;
  commitMessage?: string;
  branch?: string;
  trigger?: string;
  startedAt?: string;
  finishedAt?: string;
  durationSeconds?: number;
}

export interface ServiceVariable {
  id: string;
  key: string;
  value: string;
  isSecret?: boolean;
}
