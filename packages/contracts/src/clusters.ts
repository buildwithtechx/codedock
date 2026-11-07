import type { ClusterNode } from './servers';

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
  controls: number;
  nodes: ClusterNode[];
  revision: number;
  status: string;
  scaling?: ClusterScaling;
  error?: string;
  updatedAt?: string;
}
