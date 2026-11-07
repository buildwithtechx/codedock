import type { ApplySecurityInput } from '@codedock/contracts';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useCallback, useState } from 'react';
import { toast } from 'sonner';
import { securityService } from '#/services/security';

export const useSecurityAdvisories = () => {
  return useQuery({
    queryKey: ['security', 'advisories'],
    queryFn: async () => {
      const res = await securityService.getAdvisories();
      return res.data;
    },
    refetchInterval: 60_000,
  });
};

export const useProjectSecurity = (projectId?: string, initialData?: any) => {
  return useQuery({
    queryKey: ['security', 'project', projectId],
    queryFn: async () => {
      if (!projectId) throw new Error('Project ID required');
      const res = await securityService.getProjectSecurity(projectId);
      return res.data;
    },
    initialData,
    enabled: Boolean(projectId),
    refetchInterval: 30_000,
  });
};

export const useCheckProjectSecurity = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (projectId: string) => securityService.checkProject(projectId),
    onSuccess: async (response, projectId) => {
      toast.success(response.message || 'Security check completed');
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['security', 'project', projectId] }),
        queryClient.invalidateQueries({ queryKey: ['security', 'advisories'] }),
      ]);
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to check project security');
    },
  });
};

export const useSetProjectSecurityPolicy = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ projectId, autoRefresh }: { projectId: string; autoRefresh: boolean }) =>
      securityService.setProjectPolicy(projectId, autoRefresh),
    onSuccess: async (_, { projectId }) => {
      toast.success('Security policy updated');
      await queryClient.invalidateQueries({ queryKey: ['security', 'project', projectId] });
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to update security policy');
    },
  });
};

export const useApplySecurityRemediation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ projectId, payload }: { projectId: string; payload: ApplySecurityInput }) =>
      securityService.applyRemediation(projectId, payload),
    onSuccess: async (response, { projectId }) => {
      toast.success(response.message || 'Remediation initiated');
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['security', 'project', projectId] }),
        queryClient.invalidateQueries({ queryKey: ['security', 'advisories'] }),
        queryClient.invalidateQueries({ queryKey: ['deployments'] }),
      ]);
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to apply security remediation');
    },
  });
};

const DISMISSED_ADVISORIES_KEY = 'codedock_dismissed_advisories';

export function useDismissedAdvisories() {
  const [dismissed, setDismissed] = useState<string[]>(() => {
    if (typeof window === 'undefined') return [];
    try {
      const stored = localStorage.getItem(DISMISSED_ADVISORIES_KEY);
      return stored ? JSON.parse(stored) : [];
    } catch {
      return [];
    }
  });

  const dismiss = useCallback((advisoryId: string) => {
    setDismissed((prev) => {
      if (prev.includes(advisoryId)) return prev;
      const next = [...prev, advisoryId];
      try {
        localStorage.setItem(DISMISSED_ADVISORIES_KEY, JSON.stringify(next));
      } catch {}
      return next;
    });
  }, []);

  const reset = useCallback(() => {
    setDismissed([]);
    try {
      localStorage.removeItem(DISMISSED_ADVISORIES_KEY);
    } catch {}
  }, []);

  return { dismissed, dismiss, reset };
}
