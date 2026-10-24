import { Link } from '@tanstack/react-router';
import {
  ArrowRight,
  Cloud,
  Plus,
  RefreshCw,
  Server as ServerIcon,
  SlidersHorizontal,
} from 'lucide-react';
import { useEffect } from 'react';
import { Button } from '@codedock/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { Skeleton } from '@codedock/ui/skeleton';
import { useListServers } from '#/hooks/use-servers';
import { useGetPublicSettings } from '#/hooks/use-settings';

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
      <div className="space-y-3.5 rounded-2xl border border-border/50 bg-card p-5 shadow-xs">
        <h3 className="font-semibold text-foreground text-sm">Where to install</h3>

        {isLoading ? (
          <Skeleton className="h-16 w-full rounded-xl" />
        ) : (
          <div className="space-y-3">
            <div className="flex min-w-0 items-center gap-3 overflow-hidden rounded-xl border border-border/50 bg-muted/20 p-3.5">
              <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                {isCloud ? <Cloud className="size-4" /> : <ServerIcon className="size-4" />}
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

            <div className="flex items-center justify-between gap-2">
              <p className="text-muted-foreground text-xs">Uses the existing plan.</p>
              <Link
                to="/servers"
                className="inline-flex shrink-0 items-center gap-1.5 rounded-xl border border-border/60 bg-muted/30 px-3 py-1.5 font-medium text-foreground text-xs transition-colors hover:bg-muted/60"
              >
                <Plus className="size-3.5" />
                {isCloud ? 'New Cloud space' : 'New dedicated server'}
              </Link>
            </div>
          </div>
        )}

        {needsUpgrade && (
          <div className="space-y-2.5 rounded-xl border border-primary/20 bg-primary/10 p-4 text-xs">
            <p className="font-medium text-primary">This app needs a larger Cloud plan</p>
            <p className="text-muted-foreground leading-relaxed">
              Your plan can deploy static sites only. Service stacks need a paid plan.
            </p>
            <button
              type="button"
              onClick={() => void onDeploy()}
              className="inline-flex items-center gap-1.5 font-medium text-foreground hover:underline"
            >
              <RefreshCw className="size-3" />
              Check again
            </button>
          </div>
        )}
      </div>

      <div className="space-y-3">
        {needsUpgrade ? (
          <Button asChild className="h-11 w-full gap-2 rounded-xl font-semibold text-sm shadow-sm">
            <Link to="/billing">
              <ArrowRight className="size-4" />
              Upgrade plan
            </Link>
          </Button>
        ) : (
          <Button
            className="h-11 w-full gap-2 rounded-xl font-semibold text-sm shadow-sm"
            onClick={onDeploy}
            disabled={disabled || isDeploying}
          >
            {!isDeploying && <ArrowRight className="size-4" />}
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
