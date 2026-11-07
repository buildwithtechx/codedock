import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@codedock/ui/card';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Switch } from '@codedock/ui/switch';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { toast } from 'sonner';
import { apiClient } from '#/lib/api-client';

interface RuntimeModeCardProps {
  serviceId: string;
  initialData?: {
    runtimeMode?: 'web' | 'worker' | 'static';
    internalPort?: number;
    spaFallback?: boolean;
  };
}

export function RuntimeModeCard({ serviceId, initialData }: RuntimeModeCardProps) {
  const queryClient = useQueryClient();
  const [runtimeMode, setRuntimeMode] = useState<'web' | 'worker' | 'static'>(
    initialData?.runtimeMode || 'web'
  );
  const [internalPort, setInternalPort] = useState(initialData?.internalPort || 3000);
  const [spaFallback, setSpaFallback] = useState(initialData?.spaFallback !== false);

  const updateMutation = useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const appRes = await apiClient.get<any>(`/apps/${serviceId}`);
      const app = appRes.data?.data || appRes.data;
      return apiClient.put(`/apps/${serviceId}`, { ...app, ...data });
    },
    onSuccess: () => {
      toast.success('Runtime mode updated');
      queryClient.invalidateQueries({ queryKey: ['service', serviceId] });
    },
  });

  const handleSave = () => {
    updateMutation.mutate({
      runtimeMode,
      internalPort: runtimeMode === 'web' ? internalPort : undefined,
      spaFallback: runtimeMode === 'static' ? spaFallback : undefined,
    });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Runtime Mode</CardTitle>
        <CardDescription>Choose how this service behaves inside the cluster.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <button
            type="button"
            className={`w-full cursor-pointer rounded-lg border p-4 text-left transition-colors ${
              runtimeMode === 'web'
                ? 'border-zinc-900 bg-zinc-50 dark:border-zinc-100 dark:bg-zinc-900'
                : 'border-zinc-200 hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700'
            }`}
            onClick={() => setRuntimeMode('web')}
          >
            <div className="mb-2 flex items-center justify-between">
              <h3 className="font-semibold">Web Service</h3>
              {runtimeMode === 'web' && <Badge variant="default">Active</Badge>}
            </div>
            <p className="text-sm text-zinc-500">
              Exposes a public HTTP route and binds to an internal port. Monitored via /healthz.
            </p>
          </button>

          <button
            type="button"
            className={`w-full cursor-pointer rounded-lg border p-4 text-left transition-colors ${
              runtimeMode === 'worker'
                ? 'border-zinc-900 bg-zinc-50 dark:border-zinc-100 dark:bg-zinc-900'
                : 'border-zinc-200 hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700'
            }`}
            onClick={() => setRuntimeMode('worker')}
          >
            <div className="mb-2 flex items-center justify-between">
              <h3 className="font-semibold">Background Worker</h3>
              {runtimeMode === 'worker' && <Badge variant="default">Active</Badge>}
            </div>
            <p className="text-sm text-zinc-500">
              No exposed public route. Runs as a background process and monitored via uptime checks.
            </p>
          </button>

          <button
            type="button"
            className={`w-full cursor-pointer rounded-lg border p-4 text-left transition-colors ${
              runtimeMode === 'static'
                ? 'border-zinc-900 bg-zinc-50 dark:border-zinc-100 dark:bg-zinc-900'
                : 'border-zinc-200 hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700'
            }`}
            onClick={() => setRuntimeMode('static')}
          >
            <div className="mb-2 flex items-center justify-between">
              <h3 className="font-semibold">Static Site</h3>
              {runtimeMode === 'static' && <Badge variant="default">Active</Badge>}
            </div>
            <p className="text-sm text-zinc-500">
              Serves pre-built HTML/CSS/JS assets from your build output directory via high-speed
              HTTP.
            </p>
          </button>
        </div>

        {runtimeMode === 'web' && (
          <div className="space-y-2 rounded-md bg-zinc-50 p-4 dark:bg-zinc-900/50">
            <Label htmlFor="internalPort">Internal Port</Label>
            <Input
              id="internalPort"
              type="number"
              value={internalPort}
              onChange={(e) => setInternalPort(parseInt(e.target.value, 10))}
              className="max-w-50"
            />
            <p className="text-xs text-zinc-500">
              The port your application server listens on (e.g., 3000, 8080).
            </p>
          </div>
        )}

        {runtimeMode === 'static' && (
          <div className="flex items-center justify-between rounded-md border bg-muted/20 p-4">
            <div className="space-y-0.5 pr-4">
              <Label htmlFor="spaFallback" className="font-medium text-sm">
                SPA fallback
              </Label>
              <p className="text-muted-foreground text-xs">
                Static sites: serve index.html for missing paths. Turn off for a real 404, using
                404.html when available. Redeploy container-hosted sites to apply changes.
              </p>
            </div>
            <Switch id="spaFallback" checked={spaFallback} onCheckedChange={setSpaFallback} />
          </div>
        )}

        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={updateMutation.isPending}>
            {updateMutation.isPending ? 'Saving...' : 'Save Mode'}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
