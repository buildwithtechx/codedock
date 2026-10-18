export interface CodedockConnection {
  url: string;
  token: string;
  userName?: string;
  userEmail?: string;
  version?: string;
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

export interface CodedockApp {
  id: string;
  name: string;
  status?: 'running' | 'stopped' | 'failed' | 'building' | 'deploying' | string;
  projectId?: string;
  projectName?: string;
  domain?: string;
  gitRepo?: string;
  branch?: string;
  updatedAt?: string;
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
  durationSeconds?: number;
}

export interface CodedockDatabase {
  id: string;
  projectId?: string;
  name: string;
  engine: 'postgres' | 'redis' | 'mysql' | 'mongodb' | 'mariadb' | string;
  version?: string;
  username?: string;
  databaseName?: string;
  internalHost?: string;
  internalDns?: string;
  externalDns?: string;
  port?: number;
  externalPort?: number;
  status: 'running' | 'stopped' | 'failed' | 'starting' | string;
  containerId?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface CodedockServer {
  id: string;
  name: string;
  organizationId?: string;
  ipAddress?: string;
  isLocal?: boolean;
  sshHost?: string;
  sshPort?: number;
  sshUser?: string;
  isControlPlane?: boolean;
  status: 'online' | 'offline' | 'unreachable' | 'provisioning' | string;
  provider?: string;
  region?: string;
  serverType?: string;
  lastSeenAt?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface SystemStats {
  cpuUsage?: number;
  memoryUsage?: number;
  diskUsage?: number;
}
