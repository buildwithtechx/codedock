import { useQuery } from '@tanstack/react-query';
import type { CodedockDeployment } from '@/interfaces';
import { MobileApi } from '@/services/api';
import { useAppStore } from '@/store/useAppStore';

export function useDeploymentsQuery(limit: number = 25) {
  const connection = useAppStore((state) => state.connection);

  return useQuery<CodedockDeployment[], Error>({
    queryKey: ['deployments', connection?.url, limit],
    queryFn: async () => {
      if (!connection?.url || !connection?.token) {
        return [];
      }
      return MobileApi.listDeployments(connection.url, connection.token, limit);
    },
    enabled: Boolean(connection?.url && connection?.token),
    refetchInterval: 10000,
  });
}

export function useDeploymentLogsQuery(deploymentId: string | null) {
  const connection = useAppStore((state) => state.connection);

  return useQuery<string, Error>({
    queryKey: ['deploymentLogs', connection?.url, deploymentId],
    queryFn: async () => {
      if (!connection?.url || !connection?.token || !deploymentId) {
        return '';
      }
      return MobileApi.getDeploymentLogs(connection.url, connection.token, deploymentId);
    },
    enabled: Boolean(connection?.url && connection?.token && deploymentId),
  });
}
