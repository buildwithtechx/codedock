import { toast } from 'sonner';
import { useGetRollbackCapacity, useUpdateRollbackSettings } from '#/hooks/use-projects';
import { MAX_ROLLBACK_WINDOW } from './rollback-retention-cards';

function clampWindow(val: number): number {
  if (Number.isNaN(val) || val < 0) return 0;
  if (val > MAX_ROLLBACK_WINDOW) return MAX_ROLLBACK_WINDOW;
  return Math.floor(val);
}

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
      const next = value === null ? null : clampWindow(value);
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
