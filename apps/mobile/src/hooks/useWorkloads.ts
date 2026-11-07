import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CodedockApp } from '@/interfaces';
import { MobileApi } from '@/services/api';
import { useAppStore } from '@/store/useAppStore';

export function useAppsQuery() {
  const connection = useAppStore((state) => state.connection);

  return useQuery<CodedockApp[], Error>({
    queryKey: ['apps', connection?.url],
    queryFn: async () => {
      if (!connection?.url || !connection?.token) {
        return [];
      }
      return MobileApi.listApps(connection.url, connection.token);
    },
    enabled: Boolean(connection?.url && connection?.token),
    refetchInterval: 10000,
  });
}

export function useRedeployMutation() {
  const connection = useAppStore((state) => state.connection);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (appId: string) => {
      if (!connection?.url || !connection?.token) {
        throw new Error('Not connected');
      }
      return MobileApi.redeployApp(connection.url, connection.token, appId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apps'] });
      queryClient.invalidateQueries({ queryKey: ['deployments'] });
    },
  });
}

export function useRestartMutation() {
  const connection = useAppStore((state) => state.connection);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (appId: string) => {
      if (!connection?.url || !connection?.token) {
        throw new Error('Not connected');
      }
      return MobileApi.restartApp(connection.url, connection.token, appId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apps'] });
    },
  });
}

export function useStopMutation() {
  const connection = useAppStore((state) => state.connection);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (appId: string) => {
      if (!connection?.url || !connection?.token) {
        throw new Error('Not connected');
      }
      return MobileApi.stopApp(connection.url, connection.token, appId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apps'] });
    },
  });
}
