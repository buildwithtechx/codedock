import { Link } from '@tanstack/react-router';
import {
  ArrowUpRight,
  Box,
  ExternalLink,
  Globe,
  Play,
  RotateCw,
  Square,
  Terminal,
} from 'lucide-react';
import { useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '#/components/ui/sheet';
import type { AppService } from '#/features/services';
import { appsService } from '#/services/apps';

interface ServiceDetailDrawerProps {
  service: AppService | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onRefresh?: () => void;
}

export function ServiceDetailDrawer({
  service,
  open,
  onOpenChange,
  onRefresh,
}: ServiceDetailDrawerProps) {
  const [actionLoading, setActionLoading] = useState(false);

  if (!service) return null;

  const handleAction = async (action: 'restart' | 'redeploy' | 'stop') => {
    setActionLoading(true);
    try {
      if (action === 'restart') await appsService.restartApp(service.id);
      if (action === 'redeploy') await appsService.redeployApp(service.id);
      if (action === 'stop') await appsService.stopApp(service.id);
      onRefresh?.();
    } catch {
      // Error handled
    } finally {
      setActionLoading(false);
    }
  };

  const isRunning =
    service.status?.toLowerCase() === 'running' ||
    service.status?.toLowerCase() === 'online' ||
    service.status?.toLowerCase() === 'healthy';

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full overflow-y-auto sm:max-w-lg">
        <SheetHeader className="space-y-3 border-border/60 border-b pb-4">
          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-2.5">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Box className="h-5 w-5" />
              </div>
              <div>
                <SheetTitle className="font-semibold text-base">{service.name}</SheetTitle>
                <SheetDescription className="text-xs">
                  {service.runtimeMode || 'Web Service'}
                </SheetDescription>
              </div>
            </div>
            <Badge variant={isRunning ? 'default' : 'secondary'} className="text-xs uppercase">
              {service.status || 'Unknown'}
            </Badge>
          </div>

          <div className="flex flex-wrap gap-2 pt-2">
            <Button
              size="sm"
              variant="outline"
              className="h-8 gap-1.5 text-xs"
              onClick={() => handleAction('redeploy')}
              disabled={actionLoading}
            >
              <RotateCw className="h-3.5 w-3.5" />
              Redeploy
            </Button>
            <Button
              size="sm"
              variant="outline"
              className="h-8 gap-1.5 text-xs"
              onClick={() => handleAction('restart')}
              disabled={actionLoading}
            >
              <Play className="h-3.5 w-3.5" />
              Restart
            </Button>
            {isRunning && (
              <Button
                size="sm"
                variant="outline"
                className="h-8 gap-1.5 text-destructive text-xs hover:bg-destructive/10"
                onClick={() => handleAction('stop')}
                disabled={actionLoading}
              >
                <Square className="h-3.5 w-3.5" />
                Stop
              </Button>
            )}
            <Button size="sm" variant="secondary" className="h-8 gap-1.5 text-xs" asChild>
              <Link to="/services/$serviceId" params={{ serviceId: service.id }}>
                <Terminal className="h-3.5 w-3.5" />
                Logs & Terminal
              </Link>
            </Button>
          </div>
        </SheetHeader>

        <div className="space-y-5 py-5 text-sm">
          {service.domain && (
            <div className="space-y-2 rounded-xl border border-border/70 bg-card p-3.5">
              <span className="flex items-center gap-1.5 font-medium text-muted-foreground text-xs">
                <Globe className="h-3.5 w-3.5" /> Public Domain
              </span>
              <div className="flex items-center justify-between gap-2">
                <a
                  href={`https://${service.domain}`}
                  target="_blank"
                  rel="noreferrer"
                  className="flex items-center gap-1 truncate font-mono text-primary text-xs hover:underline"
                >
                  https://{service.domain}
                  <ExternalLink className="h-3 w-3 shrink-0" />
                </a>
              </div>
            </div>
          )}

          <div className="grid grid-cols-2 gap-3">
            <div className="rounded-xl border border-border/70 bg-card p-3">
              <span className="block font-medium text-[11px] text-muted-foreground">
                Internal Port
              </span>
              <span className="mt-1 block font-mono font-semibold text-sm">
                {service.internalPort || 3000}
              </span>
            </div>
            <div className="rounded-xl border border-border/70 bg-card p-3">
              <span className="block font-medium text-[11px] text-muted-foreground">Replicas</span>
              <span className="mt-1 block font-mono font-semibold text-sm">
                {service.replicas ?? 1}
              </span>
            </div>
          </div>

          <div className="space-y-3 rounded-xl border border-border/70 bg-card p-4">
            <h4 className="font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              Deployment Info
            </h4>
            <div className="space-y-2.5 text-xs">
              {service.branch && (
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Git Branch</span>
                  <span className="font-medium font-mono">{service.branch}</span>
                </div>
              )}
              {service.containerId && (
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Container ID</span>
                  <span className="max-w-[180px] truncate font-medium font-mono">
                    {service.containerId.slice(0, 12)}
                  </span>
                </div>
              )}
              {service.healthCheckPath && (
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Health Check</span>
                  <span className="font-medium font-mono">{service.healthCheckPath}</span>
                </div>
              )}
              {service.buildEngine && (
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Build Engine</span>
                  <span className="font-medium capitalize">{service.buildEngine}</span>
                </div>
              )}
            </div>
          </div>

          <div className="space-y-3 rounded-xl border border-border/70 bg-card p-4">
            <h4 className="font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              Resource Allocations
            </h4>
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div>
                <span className="text-muted-foreground">CPU Limit</span>
                <p className="mt-0.5 font-semibold">
                  {service.cpuLimit ? `${service.cpuLimit} vCPU` : 'Unlimited'}
                </p>
              </div>
              <div>
                <span className="text-muted-foreground">Memory Limit</span>
                <p className="mt-0.5 font-semibold">
                  {service.memoryLimit ? `${service.memoryLimit} MB` : 'Unlimited'}
                </p>
              </div>
            </div>
          </div>

          <div className="pt-2">
            <Button variant="outline" className="w-full justify-between" asChild>
              <Link to="/services/$serviceId" params={{ serviceId: service.id }}>
                <span>Manage Service Configuration</span>
                <ArrowUpRight className="h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
