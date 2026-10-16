import { apiClient } from '#/lib/api-client';

export interface ClusterNode {
  serverId: string;
  serverName?: string;
  role: 'server' | 'agent';
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
  controls: number;
  nodes: ClusterNode[];
  revision: number;
  status: string;
  scaling?: ClusterScaling;
  error?: string;
  updatedAt?: string;
}

export const clustersService = {
  list: async (projectId: string): Promise<ComputeCluster[]> => {
    const res = await apiClient.get<{ data: ComputeCluster[] }>(`/projects/${encodeURIComponent(projectId)}/clusters`);
    return res.data ?? [];
  },

  getScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.get<{ data: ClusterScaling }>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling`);
    return res.data;
  },

  setupScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.post<{ data: ClusterScaling }>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/setup`, {});
    return res.data;
  },

  retryScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.post<{ data: ClusterScaling }>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/retry`, {});
    return res.data;
  },

  removeScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.post<{ data: ClusterScaling }>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/remove`, {});
    return res.data;
  },
};
