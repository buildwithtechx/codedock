import { useQuery } from '@tanstack/react-query';
import {
  Boxes,
  CheckCircle2,
  Container,
  GitBranch,
  Globe,
  Loader2,
  RefreshCw,
  ShieldAlert,
  ShieldCheck,
  TerminalSquare,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import type { ServerComponentStatus } from '#/interfaces/server';
import { cn } from '#/lib/utils';
import { serverService } from '#/services/servers';

interface ServerComponentsTabProps {
  serverId: string;
}

function componentIcon(name: string) {
  switch (name) {
    case 'docker':
      return Container;
    case 'traefik':
      return Globe;
    case 'git':
      return GitBranch;
    case 'nixpacks':
      return Boxes;
    case 'ssh-bridge':
      return TerminalSquare;
    default:
      return Container;
  }
}

export function ServerComponentsTab({ serverId }: ServerComponentsTabProps) {
  const {
    data: components,
    isLoading,
    isError,
    refetch,
    isRefetching,
  } = useQuery({
    queryKey: ['server-components', serverId],
    queryFn: () => serverService.getComponents(serverId),
  });

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError || !components) {
    return (
      <QueryErrorState
        title="Could not load server components"
        description="Failed to query live component statuses from this node."
        onRetry={() => void refetch()}
      />
    );
  }

  const allHealthy = components.every((c) => c.healthy);

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-card">
        <div className="flex flex-wrap items-center justify-between gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex items-center gap-3">
            <div
              className={cn(
                'flex h-9 w-9 items-center justify-center rounded-xl',
                allHealthy ? 'bg-success/10 text-success' : 'bg-warning/10 text-warning'
              )}
            >
              {allHealthy ? (
                <ShieldCheck className="size-[18px]" />
              ) : (
                <ShieldAlert className="size-[18px]" />
              )}
            </div>
            <div>
              <h2 className="font-semibold text-[15px] text-foreground">System health</h2>
              <p className="text-muted-foreground text-xs">
                {allHealthy
                  ? 'All core runtime components are operational'
                  : 'Some components require attention or are initializing'}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={isRefetching}
              onClick={() => void refetch()}
            >
              <RefreshCw className={cn('mr-1.5 size-3.5', isRefetching && 'animate-spin')} />
              Check again
            </Button>
          </div>
        </div>

        <div className="divide-y divide-border/40">
          {components.map((item: ServerComponentStatus) => {
            const Icon = componentIcon(item.name);
            return (
              <div
                key={item.name}
                className="flex flex-col gap-3 p-5 sm:flex-row sm:items-center sm:justify-between"
              >
                <div className="flex min-w-0 items-start gap-3.5">
                  <div
                    className={cn(
                      'flex size-9 shrink-0 items-center justify-center rounded-xl',
                      item.healthy
                        ? 'bg-primary/10 text-primary'
                        : 'bg-destructive/10 text-destructive'
                    )}
                  >
                    <Icon className="size-4" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <h3 className="font-medium text-foreground text-sm">{item.label}</h3>
                      {item.version && (
                        <span className="rounded-md bg-muted px-2 py-0.5 font-mono text-[11px] text-muted-foreground">
                          {item.version}
                        </span>
                      )}
                      <span
                        className={cn(
                          'inline-flex items-center gap-1 rounded-full px-2 py-0.5 font-medium text-[11px]',
                          item.healthy
                            ? 'bg-success/10 text-success'
                            : 'bg-destructive/10 text-destructive'
                        )}
                      >
                        <CheckCircle2 className="size-3" />
                        {item.healthy ? 'Healthy' : 'Degraded'}
                      </span>
                    </div>
                    <p className="mt-1 text-muted-foreground text-xs">{item.description}</p>
                    <p className="mt-1 font-mono text-[11px] text-muted-foreground/80">
                      {item.message}
                    </p>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
