import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@codedock/ui/card';
import { Checkbox } from '@codedock/ui/checkbox';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import type { AppService } from '#/interfaces/services';
import { useUpdateApp } from '#/hooks/use-apps';

interface ResourceLimitsProps {
  app: AppService;
}

export function ResourceLimitsCard({ app }: ResourceLimitsProps) {
  const { mutateAsync: updateApp, isPending } = useUpdateApp();
  const [cpuLimit, setCpuLimit] = useState(app.cpuLimit?.toString() || '');
  const [memoryLimit, setMemoryLimit] = useState(app.memoryLimit?.toString() || '');
  const [swapDisabled, setSwapDisabled] = useState(Boolean(app.swapDisabled));

  // Reset local state if app data updates externally
  useEffect(() => {
    setCpuLimit(app.cpuLimit?.toString() || '');
    setMemoryLimit(app.memoryLimit?.toString() || '');
    setSwapDisabled(Boolean(app.swapDisabled));
  }, [app.cpuLimit, app.memoryLimit, app.swapDisabled]);

  const handleSave = async () => {
    const memMb = memoryLimit ? parseInt(memoryLimit, 10) : undefined;
    if (swapDisabled && (!memMb || memMb <= 0)) {
      toast.error('Disabling swap requires a positive memory limit.');
      return;
    }

    try {
      await updateApp({
        appId: app.id,
        payload: {
          ...app,
          cpuLimit: cpuLimit ? parseFloat(cpuLimit) : undefined,
          memoryLimit: memMb,
          swapDisabled,
        },
      });
      toast.success('Resource limits updated');
    } catch (error: any) {
      toast.error(error?.message || 'Failed to update resource limits');
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Resource Limits</CardTitle>
        <CardDescription>
          Restrict the amount of CPU and Memory this service can consume.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="cpuLimit">CPU Limit (Cores)</Label>
          <div className="flex items-start gap-2">
            <div className="flex-1 space-y-2">
              <Input
                id="cpuLimit"
                type="number"
                step="0.1"
                min="0.1"
                placeholder="e.g., 0.5 (Half a core) or 1 (One core)"
                value={cpuLimit}
                onChange={(e) => setCpuLimit(e.target.value)}
              />
              <p className="text-muted-foreground text-sm">
                Leave empty for unbounded CPU. Example: 0.5 = 50% of 1 core.
              </p>
            </div>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="memoryLimit">Memory Limit (MB)</Label>
          <div className="flex items-start gap-2">
            <div className="flex-1 space-y-2">
              <Input
                id="memoryLimit"
                type="number"
                step="1"
                min="64"
                placeholder="e.g., 512 or 1024"
                value={memoryLimit}
                onChange={(e) => setMemoryLimit(e.target.value)}
              />
              <p className="text-muted-foreground text-sm">
                Leave empty for unbounded memory. Minimum is 64 MB.
              </p>
            </div>
          </div>
        </div>

        <div className="flex items-start space-x-3 rounded-md border p-3">
          <Checkbox
            id="swapDisabled"
            checked={swapDisabled}
            onCheckedChange={(checked) => setSwapDisabled(Boolean(checked))}
          />
          <div className="space-y-1 leading-none">
            <Label htmlFor="swapDisabled" className="cursor-pointer font-medium text-sm">
              Disable swap memory
            </Label>
            <p className="text-muted-foreground text-xs">
              Runs this container without swap. Requires a positive memory limit and a Docker host
              with swap accounting enabled.
            </p>
          </div>
        </div>
      </CardContent>
      <CardFooter className="border-t bg-muted/50 px-6 py-4">
        <Button onClick={handleSave} disabled={isPending}>
          {isPending ? 'Saving...' : 'Save Limits'}
        </Button>
      </CardFooter>
    </Card>
  );
}
