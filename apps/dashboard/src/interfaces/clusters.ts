export type ClusterStatus = 'provisioning' | 'healthy' | 'degraded' | 'scaling';

export interface ClusterNode {
  serverId: string;
  serverName?: string;
  role?: 'server' | 'agent';
  privateIp: string;
  interface?: string;
  fingerprint?: string;
}

export interface ClusterScaling {
  status: 'not_enabled' | 'setting_up' | 'ready' | 'failed' | 'interrupted' | 'removing';
  k3sVersion: string;
  controlNodes: number;
  workerNodes: number;
  podCidr?: string;
  serviceCidr?: string;
  clusterDomain?: string;
  error?: string;
  updatedAt: string;
}

export interface ComputeCluster {
  id: string;
  projectId?: string;
  name: string;
  version: string;
  networkId?: string;
  controls?: number;
  nodes: ClusterNode[];
  revision: number;
  status: string;
  scaling?: ClusterScaling;
  error?: string;
  updatedAt?: string;
}

export type Cluster = ComputeCluster;

export interface ClusterOperation {
  id: string;
  status: string;
  phase: string;
  effects: string;
  logs: string;
  error: string;
  expiresAt: number;
}

export interface ClusterReview {
  operation: ClusterOperation;
  confirmation: string;
}

export interface ScaleDatabaseSpec {
  id: string;
  instances: number;
}

export interface ScaleDatabaseReviewInput {
  action: 'scale';
  spec: ScaleDatabaseSpec;
}

export interface ApplyDatabaseScaleInput {
  confirmation: string;
}
