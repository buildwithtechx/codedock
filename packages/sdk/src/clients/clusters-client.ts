import type { ClusterScaling, ComputeCluster } from '@codedock/contracts';
import type { HttpClient } from '../http';

export function createClustersOperations(http: HttpClient) {
  return {
    list: (projectId: string): Promise<ComputeCluster[]> =>
      http.get<ComputeCluster[]>(`/projects/${encodeURIComponent(projectId)}/clusters`),
    getScaling: (clusterId: string): Promise<ClusterScaling> =>
      http.get<ClusterScaling>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling`),
    setupScaling: (clusterId: string): Promise<ClusterScaling> =>
      http.post<ClusterScaling>(
        `/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/setup`
      ),
    retryScaling: (clusterId: string): Promise<ClusterScaling> =>
      http.post<ClusterScaling>(
        `/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/retry`
      ),
    removeScaling: (clusterId: string): Promise<ClusterScaling> =>
      http.delete<ClusterScaling>(
        `/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling`
      ),
  };
}
