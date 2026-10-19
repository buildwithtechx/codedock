import { useCallback, useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import type { RollbackCapacityUI } from '#/interfaces/projects';
import { projectsApi } from '#/services/projects';
import { MAX_ROLLBACK_WINDOW } from './rollback-retention-cards';

function clampWindow(val: number): number {
  if (Number.isNaN(val) || val < 0) return 0;
  if (val > MAX_ROLLBACK_WINDOW) return MAX_ROLLBACK_WINDOW;
  return Math.floor(val);
}

export function useRollbackSettings(projectId?: string | null, enabled = true) {
  const identity = useRef(projectId);
  identity.current = projectId;
  const mounted = useRef(false);
  const pending = useRef(false);
  const [state, setState] = useState<{
    projectId?: string | null;
    capacity: RollbackCapacityUI | null;
    error: string | null;
    loading: boolean;
  }>({ projectId, capacity: null, error: null, loading: true });
  const [saving, setSaving] = useState<'window' | 'strategy' | null>(null);

  const reload = useCallback(async () => {
    if (!mounted.current || !projectId || identity.current !== projectId) return;
    setState((prev) => ({ ...prev, loading: true, error: null }));
    try {
      const response = await projectsApi.getRollbackCapacity(projectId);
      if (mounted.current && identity.current === projectId) {
        setState({
          projectId,
          capacity: response.data,
          error: null,
          loading: false,
        });
      }
    } catch (err: any) {
      if (mounted.current && identity.current === projectId) {
        setState({
          projectId,
          capacity: null,
          loading: false,
          error: err?.message || 'Could not load rollback settings.',
        });
      }
    }
  }, [projectId]);

  useEffect(() => {
    mounted.current = true;
    if (enabled && projectId) {
      void reload();
    }
    return () => {
      mounted.current = false;
    };
  }, [enabled, projectId, reload]);

  const capacity = state.projectId === projectId ? state.capacity : null;

  const save = async (
    kind: 'window' | 'strategy',
    patch: { defaultRollbackStrategy?: 'snapshot' | 'git'; rollbackWindow?: number | null }
  ) => {
    if (!projectId || !capacity || pending.current) return;
    pending.current = true;
    setSaving(kind);
    try {
      const response = await projectsApi.updateRollbackSettings(projectId, patch);
      if (mounted.current && identity.current === projectId) {
        setState({
          projectId,
          capacity: response.data,
          error: null,
          loading: false,
        });
        toast.success(
          kind === 'strategy' ? 'Rollback strategy updated' : 'Rollback history limit updated'
        );
      }
    } catch (err: any) {
      if (identity.current === projectId) {
        toast.error(
          err?.message ||
            (kind === 'window'
              ? 'Failed to update rollback history'
              : 'Failed to update rollback strategy')
        );
      }
    } finally {
      pending.current = false;
      if (mounted.current) setSaving(null);
    }
  };

  return {
    capacity,
    error: state.projectId === projectId ? state.error : null,
    loading: state.projectId !== projectId || state.loading,
    savingWindow: saving === 'window',
    togglingStrategy: saving === 'strategy',
    reload,
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
