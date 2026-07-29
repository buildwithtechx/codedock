import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { ServiceRouteRules } from '#/features/services/components/service-route-rules';
import { useGetApp } from '#/hooks/use-apps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/route-rules')({
  component: ServiceRouteRulesRoute,
});

function ServiceRouteRulesRoute() {
  const { serviceId } = Route.useParams();
  const { data: appData, isLoading } = useGetApp(serviceId);

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;

  if (!app) {
    return <div>Service not found.</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="font-bold text-2xl">Route Rules</h1>
      <ServiceRouteRules serviceId={serviceId} />
    </div>
  );
}
