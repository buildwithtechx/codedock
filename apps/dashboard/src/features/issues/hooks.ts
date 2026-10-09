import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { useOrganizationStore } from '#/stores/organization-store';
import { issuesService } from './api';
import type { ActIssueInput } from './interfaces';

export const useActiveOrganizationId = () =>
  useOrganizationStore((state) => state.activeOrganizationId);

export const useListAttentionIssues = () => {
  const activeOrganizationId = useActiveOrganizationId();

  return useQuery({
    queryKey: ['issues', 'list', activeOrganizationId],
    queryFn: () => issuesService.list(activeOrganizationId as string),
    enabled: Boolean(activeOrganizationId),
    refetchInterval: 15_000,
  });
};

export const useEvaluateIssues = () => {
  const queryClient = useQueryClient();
  const activeOrganizationId = useActiveOrganizationId();

  return useMutation({
    mutationFn: () => issuesService.evaluate(activeOrganizationId as string),
    onSuccess: async (response) => {
      toast.success(response.message || 'Rescan finished');
      await queryClient.invalidateQueries({
        queryKey: ['issues', 'list', activeOrganizationId],
      });
    },
  });
};

export const useAcknowledgeIssue = () => {
  const queryClient = useQueryClient();
  const activeOrganizationId = useActiveOrganizationId();

  return useMutation({
    mutationFn: (issueId: string) =>
      issuesService.acknowledge(activeOrganizationId as string, issueId),
    onSuccess: async (response) => {
      toast.success(response.message || 'Issue acknowledged');
      await queryClient.invalidateQueries({
        queryKey: ['issues', 'list', activeOrganizationId],
      });
    },
  });
};

export const useResolveIssue = () => {
  const queryClient = useQueryClient();
  const activeOrganizationId = useActiveOrganizationId();

  return useMutation({
    mutationFn: (issueId: string) => issuesService.resolve(activeOrganizationId as string, issueId),
    onSuccess: async (response) => {
      toast.success(response.message || 'Issue resolved');
      await queryClient.invalidateQueries({
        queryKey: ['issues', 'list', activeOrganizationId],
      });
    },
  });
};

export const useActOnIssue = () => {
  const queryClient = useQueryClient();
  const activeOrganizationId = useActiveOrganizationId();

  return useMutation({
    mutationFn: (input: ActIssueInput) =>
      issuesService.act(activeOrganizationId as string, input.issueId, {
        action: input.action,
      }),
    onSuccess: async (response) => {
      toast.success(response.data?.outcome || response.message || 'Remediation executed');
      await queryClient.invalidateQueries({
        queryKey: ['issues', 'list', activeOrganizationId],
      });
    },
  });
};
