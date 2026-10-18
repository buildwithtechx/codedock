import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AuditListParams, AuditSettings } from '#/interfaces/audit';
import { auditLogsService } from '#/services/audit-logs';

export const useAuditLogs = (params: AuditListParams = {}) => {
  return useQuery({
    queryKey: ['auditLogs', params],
    queryFn: () => auditLogsService.list(params),
  });
};

export const useAuditFacets = (params: { category?: string; q?: string } = {}) => {
  return useQuery({
    queryKey: ['auditLogs', 'facets', params],
    queryFn: () => auditLogsService.facets(params),
  });
};

export const useAuditSettings = (organizationId?: string) => {
  return useQuery({
    queryKey: ['auditLogs', 'settings', organizationId].filter(Boolean),
    queryFn: () => auditLogsService.getSettings(organizationId),
  });
};

export const useUpdateAuditSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (settings: Partial<AuditSettings>) => auditLogsService.updateSettings(settings),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['auditLogs', 'settings'] });
      await queryClient.invalidateQueries({ queryKey: ['auditLogs'] });
    },
  });
};
