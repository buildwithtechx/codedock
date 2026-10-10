import { Cloud, Cpu, Server as ServerIcon } from 'lucide-react';
import { useEffect } from 'react';
import { Badge } from '#/components/ui/badge';
import { Label } from '#/components/ui/label';
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
  targetMode: 'self-hosted' | 'cloud';
  onTargetModeChange: (mode: 'self-hosted' | 'cloud') => void;
  serverId: string;
  onServerChange: (serverId: string) => void;
  minMemoryMb?: number;
  minCpuCores?: number;
  disabled?: boolean;
}

export function WizardDestination({
  targetMode,
  onTargetModeChange,
  serverId,
  onServerChange,
  minMemoryMb,
  minCpuCores,
  disabled = false,
}: WizardDestinationProps) {
  const { data: serversData, isLoading, isError } = useListServers();
  const { data: publicSettings } = useGetPublicSettings();
  const servers = Array.isArray(serversData) ? serversData : [];
  const cloudModeEnabled = Boolean(publicSettings?.data?.cloudMode);

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
    serverId === 'local' || !serverId ? null : servers.find((s) => s.id === serverId) || servers[0];

  return (
    <div className="space-y-4">
      <div className="flex rounded-lg bg-muted/40 p-1">
        <button
          type="button"
          disabled={disabled}
          onClick={() => onTargetModeChange('self-hosted')}
          className={`flex flex-1 items-center justify-center gap-1.5 rounded-md py-1.5 font-medium text-xs transition-colors ${
            targetMode === 'self-hosted'
              ? 'bg-card text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          <ServerIcon className="size-3.5" />
          Self-Hosted
        </button>
        <button
          type="button"
          disabled={disabled}
          onClick={() => onTargetModeChange('cloud')}
          className={`flex flex-1 items-center justify-center gap-1.5 rounded-md py-1.5 font-medium text-xs transition-colors ${
            targetMode === 'cloud'
              ? 'bg-card text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          <Cloud className="size-3.5" />
          Cloud
        </button>
      </div>

      {(minMemoryMb || minCpuCores) && (
        <div className="flex items-center gap-1.5 rounded-lg border border-border/50 bg-muted/20 px-3 py-2 text-muted-foreground text-xs">
          <Cpu className="size-3.5 shrink-0 text-primary" />
          <span className="font-medium text-foreground">Required capacity:</span>
          <span>
            {[
              minMemoryMb ? `${minMemoryMb} MB RAM` : null,
              minCpuCores ? `${minCpuCores} CPU core${minCpuCores > 1 ? 's' : ''}` : null,
            ]
              .filter(Boolean)
              .join(' · ')}
          </span>
        </div>
      )}

      {targetMode === 'cloud' ? (
        <div className="space-y-2">
          <Label>Cloud Environment</Label>
          <div className="flex items-center justify-between rounded-xl border border-border/60 bg-muted/20 p-3">
            <div className="flex items-center gap-2.5">
              <div className="flex size-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Cloud className="size-4" />
              </div>
              <div>
                <p className="font-medium text-foreground text-xs">Codedock Managed Cloud</p>
                <p className="font-mono text-[11px] text-muted-foreground">
                  {cloudModeEnabled
                    ? 'Managed Infrastructure · Auto-scaling'
                    : 'Global Edge Network'}
                </p>
              </div>
            </div>
            <Badge variant="secondary" className="text-[10px]">
              Ready
            </Badge>
          </div>
          <p className="text-[11px] text-muted-foreground">
            Workload will run on high-availability managed cloud runners.
          </p>
        </div>
      ) : isLoading ? (
        <Skeleton className="h-16 w-full rounded-xl" />
      ) : isError ? (
        <div className="rounded-xl border border-border/60 bg-muted/30 p-3 text-xs">
          <p className="text-muted-foreground">Default local container runtime will be used.</p>
        </div>
      ) : servers.length <= 1 ? (
        <div className="space-y-2">
          <Label>Target Server</Label>
          <div className="flex items-center justify-between rounded-xl border border-border/60 bg-muted/20 p-3">
            <div className="flex items-center gap-2.5">
              <div className="flex size-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <ServerIcon className="size-4" />
              </div>
              <div>
                <p className="font-medium text-foreground text-xs">
                  {activeServer?.name || 'Local Server'}
                </p>
                <p className="font-mono text-[11px] text-muted-foreground">
                  {activeServer?.ipAddress || activeServer?.sshHost || '127.0.0.1 (Docker Host)'}
                </p>
              </div>
            </div>
            <Badge variant="secondary" className="text-[10px]">
              {activeServer?.status || 'Ready'}
            </Badge>
          </div>
          <p className="text-[11px] text-muted-foreground">
            Application containers will be provisioned directly on this host.
          </p>
        </div>
      ) : (
        <div className="space-y-2">
          <Label htmlFor="wizard-server">Target Server</Label>
          <Select value={serverId || 'local'} onValueChange={onServerChange} disabled={disabled}>
            <SelectTrigger id="wizard-server" className="w-full">
              <SelectValue placeholder="Select server" />
            </SelectTrigger>
            <SelectContent>
              {servers.map((s) => (
                <SelectItem key={s.id} value={s.id}>
                  {s.name} ({s.ipAddress || s.sshHost || 'Host'})
                </SelectItem>
              ))}
              <SelectItem value="local">Local Server (This Machine)</SelectItem>
            </SelectContent>
          </Select>
          <p className="text-[11px] text-muted-foreground">
            Application containers will be provisioned directly on this host.
          </p>
        </div>
      )}
    </div>
  );
}
