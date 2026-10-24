import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CodedockDatabase } from '@/interfaces';
import { MobileApi } from '@/services/api';
import { useAppStore } from '@/store/useAppStore';

export function useDatabasesQuery() {
  const connection = useAppStore((state) => state.connection);

  return useQuery<CodedockDatabase[], Error>({
    queryKey: ['databases', connection?.url],
    queryFn: async () => {
      if (!connection?.url || !connection?.token) {
        return [];
      }
      return MobileApi.listDatabases(connection.url, connection.token);
    },
    enabled: Boolean(connection?.url && connection?.token),
    refetchInterval: 15000,
  });
}

export function useRestartDatabaseMutation() {
  const connection = useAppStore((state) => state.connection);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (dbId: string) => {
      if (!connection?.url || !connection?.token) {
        throw new Error('Not connected');
      }
      return MobileApi.restartDatabase(connection.url, connection.token, dbId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['databases'] });
    },
  });
}

export function useStartDatabaseMutation() {
  const connection = useAppStore((state) => state.connection);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (dbId: string) => {
      if (!connection?.url || !connection?.token) {
        throw new Error('Not connected');
      }
      return MobileApi.startDatabase(connection.url, connection.token, dbId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['databases'] });
    },
  });
}

export function useStopDatabaseMutation() {
  const connection = useAppStore((state) => state.connection);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (dbId: string) => {
      if (!connection?.url || !connection?.token) {
        throw new Error('Not connected');
      }
      return MobileApi.stopDatabase(connection.url, connection.token, dbId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['databases'] });
    },
  });
}

export function useRevealCredentialsMutation() {
  const connection = useAppStore((state) => state.connection);

  return useMutation({
    mutationFn: async (dbId: string) => {
      if (!connection?.url || !connection?.token) {
        throw new Error('Not connected');
      }
      return MobileApi.revealDatabaseCredentials(connection.url, connection.token, dbId);
    },
  });
}
