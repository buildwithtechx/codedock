import { useQuery } from '@tanstack/react-query';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { RuntimeStatus } from './runtime-status';
import type { ServiceRuntime } from './runtime-types';
import { ServiceMetricsPage } from './service-metrics';
export function RuntimeMetrics({ serviceId }: { serviceId: string }) {
  const runtime = useQuery({
    queryKey: ['service-runtime', serviceId],
    queryFn: () => apiClient.get<BaseResponse<ServiceRuntime>>(`/apps/${serviceId}/runtime`),
  });
  if (runtime.error) return <p role="alert">{runtime.error.message}</p>;
  if (!runtime.data) return <p role="status">Loading deployment destination?</p>;
  return runtime.data.data.target.kind !== 'docker' ? (
    <RuntimeStatus serviceId={serviceId} />
  ) : (
    <ServiceMetricsPage serviceId={serviceId} />
  );
}
