export interface CliConfig {
  serverUrl?: string;
  token?: string;
  email?: string;
}

export interface CliContext {
  config: CliConfig;
  serverUrl: string;
  token?: string;
  json: boolean;
}

export interface ApiResponse<T = unknown> {
  status?: string;
  message?: string;
  data?: T;
}

export interface UserProfile {
  id: string;
  email: string;
  name: string;
  role: string;
  planType?: string;
}

export interface ServerRecord {
  id: string;
  name: string;
  ipAddress: string;
  port?: number;
  user?: string;
  isLocal: boolean;
  status: string;
  dockerVersion?: string;
  createdAt: string;
}

export interface ProjectRecord {
  id: string;
  name: string;
  description?: string;
  organizationId?: string;
  serverId?: string;
  createdAt: string;
}

export interface ServiceRecord {
  id: string;
  name: string;
  projectId: string;
  status: string;
  domain?: string;
  repositoryUrl?: string;
  branch?: string;
  createdAt: string;
  updatedAt: string;
}

export interface DeploymentRecord {
  id: string;
  serviceId?: string;
  status: string;
  branch?: string;
  commitHash?: string;
  createdAt: string;
}

export interface DatabaseRecord {
  id: string;
  name: string;
  engine: string;
  status: string;
  serverId?: string;
  createdAt: string;
}
