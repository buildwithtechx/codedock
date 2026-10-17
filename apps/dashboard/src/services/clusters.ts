import type { ClusterNode, ClusterScaling, ComputeCluster } from '#/interfaces/clusters';
import { apiClient } from '#/lib/api-client';

export const clustersService = {
  list: async (projectId: string): Promise<ComputeCluster[]> => {
    const res = await apiClient.get<{ data: ComputeCluster[] }>(
      `/projects/${encodeURIComponent(projectId)}/clusters`
    );
    return res.data ?? [];
  },

  getScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.get<{ data: ClusterScaling }>(
      `/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling`
    );
    return res.data;
  },

  setupScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.post<{ data: ClusterScaling }>(
      `/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/setup`,
      {}
    );
    return res.data;
  },

  retryScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.post<{ data: ClusterScaling }>(
      `/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/retry`,
      {}
    );
    return res.data;
  },

  removeScaling: async (clusterId: string): Promise<ClusterScaling> => {
    const res = await apiClient.post<{ data: ClusterScaling }>(
      `/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/remove`,
      {}
    );
    return res.data;
  },
};

export type { ComputeCluster, ClusterNode, ClusterScaling };
