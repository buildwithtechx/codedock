import { Loader2, Sparkles } from 'lucide-react';
import { toast } from 'sonner';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Switch } from '#/components/ui/switch';
import { useGetApp, useRedeployApp, useUpdateApp } from '#/hooks/use-apps';

export function OptimizeCard({ serviceId }: { serviceId: string }) {
  const { data: appData, isLoading } = useGetApp(serviceId);
  const updateApp = useUpdateApp();
  const redeployApp = useRedeployApp();

  if (isLoading) {
    return (
      <Card>
        <CardContent className="flex justify-center p-6">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    );
  }

  const app = appData?.data;
  if (!app) return null;

  const handleToggle = async (checked: boolean) => {
    try {
      await updateApp.mutateAsync({
        appId: serviceId,
        payload: { ...app, optimize: checked },
      });
      toast.success('Cache warming updated. Redeploying service to apply changes...');
      await redeployApp.mutateAsync({ appId: serviceId });
      toast.success('Service redeployment started successfully');
    } catch (_error) {
      toast.error('Failed to update cache warming');
    }
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-violet-500" />
          <CardTitle>Framework Cache Warming</CardTitle>
        </div>
        <CardDescription>
          Run the framework's cache-warming command before the server starts, when the stack defines
          one (Laravel runs <code>php artisan optimize</code>).
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between rounded-lg border p-4">
          <div className="space-y-0.5">
            <h4 className="font-medium text-sm">Warm caches on deploy</h4>
            <p className="text-muted-foreground text-sm">
              The warm step runs before startup; a failure prevents the server from starting.
            </p>
          </div>
          <Switch
            checked={!!app.optimize}
            onCheckedChange={handleToggle}
            disabled={updateApp.isPending || redeployApp.isPending}
          />
        </div>
      </CardContent>
    </Card>
  );
}
