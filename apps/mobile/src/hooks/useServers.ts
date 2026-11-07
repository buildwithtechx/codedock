import { useQuery } from '@tanstack/react-query';
import type { CodedockServer } from '@/interfaces';
import { MobileApi } from '@/services/api';
import { useAppStore } from '@/store/useAppStore';

export function useServersQuery() {
  const connection = useAppStore((state) => state.connection);

  return useQuery<CodedockServer[], Error>({
    queryKey: ['servers', connection?.url],
    queryFn: async () => {
      if (!connection?.url || !connection?.token) {
        return [];
      }
      return MobileApi.listServers(connection.url, connection.token);
    },
    enabled: Boolean(connection?.url && connection?.token),
    refetchInterval: 15000,
  });
}
