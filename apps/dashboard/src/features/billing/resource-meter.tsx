import { Cpu, HardDrive, Layers, Network } from 'lucide-react';

export interface ResourceMeterProps {
  label: string;
  hint?: string;
  used: number | null;
  max?: number | null;
  unit?: string;
  icon?: React.ReactNode;
}

export function ResourceMeter({ label, hint, used, max, unit = '', icon }: ResourceMeterProps) {
  const isKnown = used !== null && Number.isFinite(used);
  const hasLimit = typeof max === 'number' && Number.isFinite(max) && max > 0;
  const fraction = isKnown && hasLimit ? Math.min(1, Math.max(0, used / max)) : 0;
  const percentage = Math.round(fraction * 100);

  const barColor =
    percentage >= 100 ? 'bg-destructive' : percentage >= 80 ? 'bg-amber-500' : 'bg-primary';

  const usedDisplay = isKnown ? used.toLocaleString() : '0';
  const maxDisplay = hasLimit ? max.toLocaleString() : 'Unlimited';

  return (
    <div
      className="flex flex-col justify-between rounded-xl border border-border bg-card p-4 shadow-sm"
      role={hasLimit ? 'meter' : undefined}
      aria-label={label}
      aria-valuenow={hasLimit ? Math.min(max!, Math.max(0, used ?? 0)) : undefined}
      aria-valuemin={hasLimit ? 0 : undefined}
      aria-valuemax={hasLimit ? max! : undefined}
    >
      <div>
        <div className="flex items-center justify-between">
          <span className="flex items-center gap-2 font-medium text-foreground text-sm">
            {icon}
            {label}
          </span>
          <span className="font-semibold text-foreground text-sm">
            {usedDisplay}{' '}
            {unit && <span className="font-normal text-muted-foreground text-xs">{unit}</span>}
          </span>
        </div>

        {hint && <p className="mt-1 text-muted-foreground text-xs">{hint}</p>}

        <div className="mt-3.5 h-2 w-full overflow-hidden rounded-full bg-muted">
          <div
            className={`h-full rounded-full transition-all duration-300 ${barColor}`}
            style={{ width: `${hasLimit ? percentage : 0}%` }}
          />
        </div>
      </div>

      <div className="mt-3 flex items-center justify-between text-muted-foreground text-xs">
        <span>
          Limit: {maxDisplay} {unit}
        </span>
        {hasLimit && <span>{percentage}% used</span>}
      </div>
    </div>
  );
}

export interface ResourceUsageStats {
  cpuUsed: number;
  cpuLimit: number;
  memoryUsedGb: number;
  memoryLimitGb: number;
  storageUsedGb: number;
  storageLimitGb: number;
  bandwidthUsedGb: number;
  bandwidthLimitGb: number;
}

export function ResourceMetersOverview({ stats }: { stats?: Partial<ResourceUsageStats> }) {
  return (
    <div className="space-y-4 rounded-2xl border border-border bg-card p-6 shadow-sm">
      <div>
        <h3 className="font-semibold text-foreground text-lg">Resource Utilization</h3>
        <p className="mt-1 text-muted-foreground text-sm">
          Live meters tracking allocated compute against purchased plan limits.
        </p>
      </div>

      <div className="grid gap-3.5 sm:grid-cols-2 lg:grid-cols-4">
        <ResourceMeter
          label="Dedicated vCPU"
          hint="Active vCPU allocations"
          used={stats?.cpuUsed ?? 0}
          max={stats?.cpuLimit ?? 8}
          unit="Cores"
          icon={<Cpu className="size-4 text-primary" />}
        />
        <ResourceMeter
          label="System Memory"
          hint="Allocated RAM across apps"
          used={stats?.memoryUsedGb ?? 0}
          max={stats?.memoryLimitGb ?? 16}
          unit="GB"
          icon={<Layers className="size-4 text-primary" />}
        />
        <ResourceMeter
          label="NVMe Storage"
          hint="Persisted database & app volumes"
          used={stats?.storageUsedGb ?? 0}
          max={stats?.storageLimitGb ?? 200}
          unit="GB"
          icon={<HardDrive className="size-4 text-primary" />}
        />
        <ResourceMeter
          label="Network Egress"
          hint="Monthly outbound data transfer"
          used={stats?.bandwidthUsedGb ?? 0}
          max={stats?.bandwidthLimitGb ?? 1000}
          unit="GB"
          icon={<Network className="size-4 text-primary" />}
        />
      </div>
    </div>
  );
}
