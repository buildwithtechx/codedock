import { Link } from '@tanstack/react-router';
import { Plus, RefreshCw, Server as ServerIcon, SlidersHorizontal } from 'lucide-react';
import { useEffect } from 'react';
import { Button } from '#/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { Skeleton } from '#/components/ui/skeleton';
import { useGetPublicSettings } from '#/features/settings';
import { useListServers } from '#/hooks/use-servers';

interface WizardDestinationProps {
  serverId: string;
  onServerChange: (serverId: string) => void;
  appName: string;
  isDeploying: boolean;
  onDeploy: () => void;
  onReview: () => void;
  disabled?: boolean;
}

export function WizardDestination({
  serverId,
  onServerChange,
  appName,
  isDeploying,
  onDeploy,
  onReview,
  disabled = false,
}: WizardDestinationProps) {
  const { data: serversData, isLoading } = useListServers();
  const { data: publicSettings } = useGetPublicSettings();
  const servers = Array.isArray(serversData) ? serversData : [];
  const isCloud = Boolean(publicSettings?.data?.cloudMode);

  useEffect(() => {
    if (!serverId && !isLoading) {
      if (servers.length > 0) {
        onServerChange(servers[0].id);
      } else {
        onServerChange('local');
      }
    }
  }, [serverId, servers, isLoading, onServerChange]);

  const activeServer =
    servers.find((s) => s.id === serverId) || (servers.length > 0 ? servers[0] : null);

  const serverDisplayName = isCloud
    ? activeServer?.name || 'Production'
    : activeServer?.name || 'Local Control Plane';

  const serverSub = isCloud
    ? '0 projects · Choose a plan'
    : activeServer?.isLocal
      ? 'local deployment runtime'
      : activeServer
        ? `${activeServer.sshUser ?? 'root'}@${activeServer.sshHost ?? activeServer.ipAddress}`
        : 'local deployment runtime';

  const needsUpgrade = isCloud;

  return (
    <div className="space-y-4">
      <div className="space-y-3 rounded-2xl bg-card p-5">
        <h3 className="font-semibold text-foreground text-sm">Where to install</h3>

        {isLoading ? (
          <Skeleton className="h-16 w-full rounded-xl" />
        ) : (
          <div className="space-y-3">
            <div className="flex min-w-0 items-center gap-3 overflow-hidden rounded-xl border border-border/50 bg-muted/20 p-3.5">
              <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                <ServerIcon className="size-4" />
              </div>
              <div className="min-w-0 flex-1 overflow-hidden">
                {servers.length > 1 ? (
                  <Select
                    value={serverId || activeServer?.id || 'local'}
                    onValueChange={onServerChange}
                    disabled={disabled}
                  >
                    <SelectTrigger className="h-7 w-full min-w-0 max-w-full justify-between overflow-hidden border-none bg-transparent p-0 font-medium text-foreground text-sm shadow-none focus:ring-0 [&>span]:block [&>span]:min-w-0 [&>span]:truncate">
                      <SelectValue placeholder={serverDisplayName} />
                    </SelectTrigger>
                    <SelectContent>
                      {servers.map((s) => (
                        <SelectItem key={s.id} value={s.id} className="truncate">
                          {s.name} ({s.ipAddress || s.sshHost || 'Host'})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                ) : (
                  <p className="truncate font-medium text-foreground text-sm">
                    {serverDisplayName}
                  </p>
                )}
                <p className="truncate text-muted-foreground text-xs">{serverSub}</p>
              </div>
            </div>

            <p className="text-muted-foreground text-xs">Uses this server's plan.</p>

            <Button
              variant="outline"
              size="sm"
              asChild
              className="w-full justify-center gap-1.5 border-dashed text-xs"
            >
              <Link to="/servers">
                <Plus className="size-3.5" />
                New dedicated server
              </Link>
            </Button>
          </div>
        )}

        {needsUpgrade && (
          <div className="space-y-3 rounded-xl border border-amber-500/20 bg-amber-500/10 p-4 text-xs">
            <p className="font-medium text-amber-500">This app needs a larger Cloud plan</p>
            <p className="text-muted-foreground leading-relaxed">
              Your plan can deploy static sites only. Service stacks need a paid plan.
            </p>
            <button
              type="button"
              onClick={() => void onDeploy()}
              className="inline-flex items-center gap-1 font-medium text-foreground hover:underline"
            >
              <RefreshCw className="size-3" />
              Check again
            </button>
          </div>
        )}
      </div>

      <div className="space-y-3">
        {needsUpgrade ? (
          <Button asChild className="w-full py-5 font-medium text-sm">
            <Link to="/billing">
              Upgrade plan
              <span className="ml-1">↗</span>
            </Link>
          </Button>
        ) : (
          <Button
            className="w-full py-5 font-medium text-sm"
            onClick={onDeploy}
            disabled={disabled || isDeploying}
          >
            {isDeploying ? 'Installing…' : `Install ${appName}`}
          </Button>
        )}

        <button
          type="button"
          onClick={onReview}
          disabled={disabled}
          className="flex w-full items-center justify-center gap-1.5 text-muted-foreground text-xs transition-colors hover:text-foreground"
        >
          <SlidersHorizontal className="size-3.5" />
          Advanced deployment options
        </button>
      </div>
    </div>
  );
}
