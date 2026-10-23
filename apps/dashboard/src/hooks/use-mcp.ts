import { useQuery } from '@tanstack/react-query';
import { mcpService } from '#/services/mcp';

export const useMcpCatalog = () => {
  return useQuery({
    queryKey: ['mcp', 'catalog'],
    queryFn: () => mcpService.getCatalog(),
  });
};
