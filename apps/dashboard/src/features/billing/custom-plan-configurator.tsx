import { Cpu, HardDrive, Layers, Server } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Button } from '@codedock/ui/button';

export interface CustomPlanQuote {
  cores: number;
  ramGb: number;
  diskGb: number;
  monthlyPriceCents: number;
}

export interface CustomPlanConfiguratorProps {
  initialCores?: number;
  initialRamGb?: number;
  initialDiskGb?: number;
  disabled?: boolean;
  busy?: boolean;
  onSelectPlan?: (quote: CustomPlanQuote) => void;
}

const COST_PER_CORE = 5;
const COST_PER_RAM_GB = 3;
const COST_PER_DISK_GB = 0.1;
const BASE_PLATFORM_COST = 10;

export function CustomPlanConfigurator({
  initialCores = 2,
  initialRamGb = 4,
  initialDiskGb = 50,
  disabled = false,
  busy = false,
  onSelectPlan,
}: CustomPlanConfiguratorProps) {
  const [cores, setCores] = useState(initialCores);
  const [ramGb, setRamGb] = useState(initialRamGb);
  const [diskGb, setDiskGb] = useState(initialDiskGb);

  const quote = useMemo<CustomPlanQuote>(() => {
    const rawTotal =
      BASE_PLATFORM_COST +
      cores * COST_PER_CORE +
      ramGb * COST_PER_RAM_GB +
      diskGb * COST_PER_DISK_GB;
    const monthlyPriceCents = Math.round(rawTotal * 100);
    return {
      cores,
      ramGb,
      diskGb,
      monthlyPriceCents,
    };
  }, [cores, ramGb, diskGb]);

  const monthlyPriceFormatted = `$${(quote.monthlyPriceCents / 100).toFixed(2)}`;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!disabled && !busy && onSelectPlan) {
      onSelectPlan(quote);
    }
  };

  return (
    <form
      aria-label="Custom Plan Configurator"
      onSubmit={handleSubmit}
      className="grid gap-6 rounded-2xl border border-border bg-card p-6 shadow-sm lg:grid-cols-3"
    >
      <div className="space-y-6 lg:col-span-2">
        <div>
          <h3 className="font-semibold text-foreground text-lg">Custom Dedicated Cloud Plan</h3>
          <p className="mt-1 text-muted-foreground text-sm">
            Configure exact compute and storage capacity for your cluster and standalone apps.
          </p>
        </div>

        <div className="space-y-5">
          <div>
            <div className="flex items-center justify-between text-sm">
              <span className="flex items-center gap-2 font-medium text-foreground">
                <Cpu className="size-4 text-primary" />
                Dedicated vCPU
              </span>
              <span className="font-semibold text-foreground">{cores} Cores</span>
            </div>
            <input
              type="range"
              aria-label="Dedicated vCPU Cores"
              min={1}
              max={32}
              step={1}
              value={cores}
              disabled={disabled || busy}
              onChange={(e) => setCores(Number(e.target.value))}
              className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-muted accent-primary"
            />
            <div className="mt-1 flex justify-between text-muted-foreground text-xs">
              <span>1 Core</span>
              <span>16 Cores</span>
              <span>32 Cores</span>
            </div>
          </div>

          <div>
            <div className="flex items-center justify-between text-sm">
              <span className="flex items-center gap-2 font-medium text-foreground">
                <Layers className="size-4 text-primary" />
                System Memory
              </span>
              <span className="font-semibold text-foreground">{ramGb} GB RAM</span>
            </div>
            <input
              type="range"
              aria-label="System Memory RAM"
              min={1}
              max={64}
              step={1}
              value={ramGb}
              disabled={disabled || busy}
              onChange={(e) => setRamGb(Number(e.target.value))}
              className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-muted accent-primary"
            />
            <div className="mt-1 flex justify-between text-muted-foreground text-xs">
              <span>1 GB</span>
              <span>32 GB</span>
              <span>64 GB</span>
            </div>
          </div>

          <div>
            <div className="flex items-center justify-between text-sm">
              <span className="flex items-center gap-2 font-medium text-foreground">
                <HardDrive className="size-4 text-primary" />
                Fast NVMe Storage
              </span>
              <span className="font-semibold text-foreground">{diskGb} GB SSD</span>
            </div>
            <input
              type="range"
              aria-label="Fast NVMe Storage"
              min={20}
              max={1000}
              step={10}
              value={diskGb}
              disabled={disabled || busy}
              onChange={(e) => setDiskGb(Number(e.target.value))}
              className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-muted accent-primary"
            />
            <div className="mt-1 flex justify-between text-muted-foreground text-xs">
              <span>20 GB</span>
              <span>500 GB</span>
              <span>1,000 GB</span>
            </div>
          </div>
        </div>
      </div>

      <div className="flex flex-col justify-between rounded-xl border border-border bg-muted/30 p-5">
        <div className="space-y-4">
          <div>
            <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
              Monthly Total
            </p>
            <p className="mt-1 font-bold text-3xl text-foreground tracking-tight">
              {monthlyPriceFormatted}
              <span className="font-normal text-muted-foreground text-sm">/mo</span>
            </p>
          </div>

          <div className="space-y-2 border-border/60 border-t pt-3 text-sm">
            <div className="flex justify-between text-muted-foreground">
              <span>Base Platform</span>
              <span>${BASE_PLATFORM_COST}.00</span>
            </div>
            <div className="flex justify-between text-muted-foreground">
              <span>{cores} vCPUs</span>
              <span>${(cores * COST_PER_CORE).toFixed(2)}</span>
            </div>
            <div className="flex justify-between text-muted-foreground">
              <span>{ramGb} GB RAM</span>
              <span>${(ramGb * COST_PER_RAM_GB).toFixed(2)}</span>
            </div>
            <div className="flex justify-between text-muted-foreground">
              <span>{diskGb} GB NVMe</span>
              <span>${(diskGb * COST_PER_DISK_GB).toFixed(2)}</span>
            </div>
          </div>
        </div>

        <div className="mt-6">
          <Button type="submit" className="w-full gap-2" disabled={disabled || busy}>
            <Server className="size-4" />
            Configure & Deploy
          </Button>
        </div>
      </div>
    </form>
  );
}
