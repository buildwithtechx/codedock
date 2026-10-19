import { clampRollbackWindow } from '@codedock/core';
import { toast } from 'sonner';
import { useGetRollbackCapacity, useUpdateRollbackSettings } from '#/hooks/use-projects';

export function useRollbackSettings(projectId?: string | null, enabled = true) {
  const {
    data: capacityResponse,
    isLoading,
    error: queryError,
    refetch,
  } = useGetRollbackCapacity(projectId as string, Boolean(projectId && enabled));

  const updateMutation = useUpdateRollbackSettings();

  const capacity = capacityResponse?.data ?? null;

  const save = async (
    kind: 'window' | 'strategy',
    patch: { defaultRollbackStrategy?: 'snapshot' | 'git'; rollbackWindow?: number | null }
  ) => {
    if (!projectId || !capacity || updateMutation.isPending) return;
    try {
      await updateMutation.mutateAsync({ projectId, payload: patch });
      toast.success(
        kind === 'strategy' ? 'Rollback strategy updated' : 'Rollback history limit updated'
      );
    } catch (err: any) {
      toast.error(
        err?.message ||
          (kind === 'window'
            ? 'Failed to update rollback history'
            : 'Failed to update rollback strategy')
      );
    }
  };

  const isPending = updateMutation.isPending;
  const pendingPayload = updateMutation.variables?.payload;

  return {
    capacity,
    error: (queryError as any)?.message ?? null,
    loading: Boolean(projectId) && isLoading,
    savingWindow: isPending && pendingPayload?.rollbackWindow !== undefined,
    togglingStrategy: isPending && pendingPayload?.defaultRollbackStrategy !== undefined,
    reload: refetch,
    changeWindow: async (value: number | null) => {
      const next = value === null ? null : clampRollbackWindow(value, capacity?.maxWindow);
      if (
        capacity &&
        (next === null
          ? capacity.explicit !== null
          : next !== capacity.window || capacity.explicit === null)
      ) {
        await save('window', { rollbackWindow: next });
      }
    },
    toggleStrategy: () =>
      save('strategy', {
        defaultRollbackStrategy: capacity?.strategy === 'snapshot' ? 'git' : 'snapshot',
      }),
  };
}

export default useRollbackSettings;
