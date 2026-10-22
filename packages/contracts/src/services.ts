export type RuntimeMode = 'web' | 'worker' | 'static';
export type BuildEngine =
  | 'nixpacks'
  | 'dockerfile'
  | 'buildpacks'
  | 'static'
  | 'railpack'
  | 'serverless';
export type ServiceStatus = 'created' | 'building' | 'running' | 'stopped' | 'error';
export type JobStatus = 'active' | 'inactive' | 'running' | 'completed' | 'failed' | 'error';

export interface EnvExampleVariableSuggestion {
  key: string;
  value: string;
  label: string;
  sourcePath: string;
}

export interface AppService {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  repositoryUrl: string;
  imageRef?: string;
  branch: string;
  rootDirectory: string;
  runtimeMode: RuntimeMode;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: BuildEngine;
  internalPort: number;
  domain: string;
  staticOutput: string;
  healthCheckPath: string;
  containerId: string;
  status: ServiceStatus;
  createdAt: string;
  updatedAt: string;
  cpuLimit?: number;
  memoryLimit?: number;
  swapDisabled?: boolean;
  enablePRPreviews: boolean;
  maintenanceMode: boolean;
  optimize: boolean;
  stack?: string;
  stopGracePeriodSeconds?: number;
  icon?: string;
  deployToken?: string;
  appId?: string;
  replicas?: number;
}

export interface Variable {
  id: string;
  serviceId: string;
  projectId: string;
  environmentId: string;
  key: string;
  value: string;
  isSecret: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Job {
  id: string;
  projectId?: string;
  serviceId: string;
  name: string;
  schedule: string;
  command: string;
  status: JobStatus;
  lastRunAt?: string;
  lastOutput: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAppServiceRequest {
  id?: string;
  projectId: string;
  name: string;
  repositoryUrl: string;
  imageRef?: string;
  branch: string;
  rootDirectory: string;
  runtimeMode: RuntimeMode;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: BuildEngine;
  internalPort: number;
  domain: string;
  staticOutput: string;
  healthCheckPath: string;
  cpuLimit?: number;
  memoryLimit?: number;
  stack?: string;
}

export interface UpdateAppServiceRequest {
  name: string;
  repositoryUrl: string;
  imageRef?: string;
  branch: string;
  rootDirectory: string;
  runtimeMode: RuntimeMode;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: BuildEngine;
  internalPort: number;
  domain: string;
  staticOutput: string;
  healthCheckPath: string;
  containerId: string;
  status: ServiceStatus;
  cpuLimit?: number;
  memoryLimit?: number;
  icon?: string;
  deployToken?: string;
  enablePRPreviews?: boolean;
  maintenanceMode?: boolean;
  optimize?: boolean;
  stack?: string;
  stopGracePeriodSeconds?: number;
}

export interface CreateServiceVarRequest {
  key: string;
  value: string;
  isSecret: boolean;
  expectedTopologyRevision?: string;
}

export interface UpdateServiceVarRequest {
  key: string;
  value: string;
  isSecret: boolean;
}

export interface CreateJobRequest {
  serviceId: string;
  name: string;
  schedule: string;
  command: string;
}

export interface UpdateJobRequest {
  name?: string;
  schedule?: string;
  command?: string;
  status?: JobStatus;
}
