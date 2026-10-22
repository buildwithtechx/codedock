import { useCallback, useState } from 'react';
import { toast } from 'sonner';
import type { RuntimeImageCheckResult } from '@codedock/contracts';
import { apiClient } from '#/lib/api-client';

export function useRuntimeUpdate(projectId: string, serviceId?: string) {
  const [checking, setChecking] = useState(false);
  const [updating, setUpdating] = useState(false);
  const [checkResult, setCheckResult] = useState<RuntimeImageCheckResult | null>(null);

  const checkUpdates = useCallback(async () => {
    if (checking) return;
    setChecking(true);
    try {
      const endpoint = serviceId
        ? `/apps/${serviceId}/runtime-updates`
        : `/projects/${projectId}/runtime-updates`;
      const res: any = await apiClient.get(endpoint);
      const data: RuntimeImageCheckResult = res.data?.data || res.data;
      setCheckResult(data);

      const images = data?.images || [];
      const hasUpdates = images.some((img) => img.status === 'update_available');
      if (hasUpdates) {
        toast.info('Newer base images with security patches are available.');
      } else if (images.length > 0 && images.every((img) => img.status === 'current' || img.status === 'pinned')) {
        toast.success('All container runtime base images are up to date.');
      } else {
        toast.info('No newer base images detected.');
      }
    } catch (err: any) {
      toast.error(err?.response?.data?.message || err?.message || 'Failed to check runtime image updates');
    } finally {
      setChecking(false);
    }
  }, [projectId, serviceId, checking]);

  const update = useCallback(async () => {
    if (updating) return;
    setUpdating(true);
    try {
      const endpoint = serviceId
        ? `/apps/${serviceId}/runtime-updates`
        : `/projects/${projectId}/runtime-updates`;
      await apiClient.post(endpoint, {});
      toast.success('Redeploying with latest upstream container images...');
    } catch (err: any) {
      toast.error(err?.response?.data?.message || err?.message || 'Failed to trigger runtime image redeploy');
    } finally {
      setUpdating(false);
    }
  }, [projectId, serviceId, updating]);

  return {
    checkUpdates,
    checking,
    update,
    updating,
    checkResult,
  };
}
