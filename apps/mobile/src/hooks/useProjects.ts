import { useQuery } from '@tanstack/react-query';
import type { CodedockProject } from '@/interfaces';
import { MobileApi } from '@/services/api';
import { useAppStore } from '@/store/useAppStore';

export function useProjectsQuery() {
  const connection = useAppStore((state) => state.connection);

  return useQuery<CodedockProject[], Error>({
    queryKey: ['projects', connection?.url],
    queryFn: async () => {
      if (!connection?.url || !connection?.token) {
        return [];
      }
      return MobileApi.listProjects(connection.url, connection.token);
    },
    enabled: Boolean(connection?.url && connection?.token),
    refetchInterval: 15000,
  });
}
