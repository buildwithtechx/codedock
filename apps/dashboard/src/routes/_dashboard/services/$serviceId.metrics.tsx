import { createFileRoute } from '@tanstack/react-router';
import { RuntimeMetrics } from '#/features/services/runtime-metrics';

export const Route = createFileRoute('/_dashboard/services/$serviceId/metrics')({
  component: ServiceMetricsRoute,
});

function ServiceMetricsRoute() {
  const { serviceId } = Route.useParams();
  return <RuntimeMetrics serviceId={serviceId} />;
}
