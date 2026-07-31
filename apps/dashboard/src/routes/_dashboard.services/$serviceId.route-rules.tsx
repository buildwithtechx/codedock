import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { ServiceRouteRules } from '#/features/services/components/service-route-rules';
import { useGetApp } from '#/hooks/use-apps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/route-rules')({
  component: ServiceRouteRulesRoute,
});

function ServiceRouteRulesRoute() {
  const { serviceId } = Route.useParams();
  const { data: appData, isLoading, isError, error, refetch } = useGetApp(serviceId);

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 p-12">
        <div className="font-medium text-destructive">Failed to load service</div>
        <div className="text-muted-foreground text-sm">{error?.message}</div>
        <Button onClick={() => refetch()} variant="outline">
          Retry
        </Button>
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
