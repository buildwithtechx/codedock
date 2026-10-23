import { Button } from '@codedock/ui/button';
import {
  ArrowLeft,
  ArrowRight,
  Database,
  Infinity as InfinityIcon,
  Plus,
  Sliders,
} from 'lucide-react';
import { useState } from 'react';
import { RollbackRetentionCards } from '#/components/rollback/rollback-retention-cards';
import { AddServerModal } from '#/features/servers/add-server-modal';

interface DeployTargetStepProps {
  onContinue: () => void;
}

export function DeployTargetStep({ onContinue }: DeployTargetStepProps) {
  const [machinePower, setMachinePower] = useState<'full' | 'custom'>('full');
  const [rollbackReuse, setRollbackReuse] = useState(false);
  const [rollbackHistory, setRollbackHistory] = useState(5);
  const [addServerOpen, setAddServerOpen] = useState(false);

  return (
    <div className="grid gap-5 lg:grid-cols-[1fr_340px]">
      <div className="space-y-5">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="font-semibold text-2xl text-foreground tracking-tight">
              Where do you want to deploy?
            </h1>
            <p className="mt-1 text-muted-foreground text-sm">
              Choose where your application will run
            </p>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={onContinue}
            className="gap-1.5 rounded-xl text-xs"
          >
            <ArrowLeft className="size-3.5" />
            Back
          </Button>
        </div>

        <section className="rounded-2xl border border-border/60 bg-card p-6">
          <h2 className="font-semibold text-foreground text-sm">Destination</h2>
          <div className="mt-4 flex flex-col gap-3 rounded-xl border border-border/60 bg-muted/20 p-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-3.5">
              <div className="flex size-10 items-center justify-center rounded-xl bg-muted text-muted-foreground">
                <Database className="size-5" />
              </div>
              <div>
                <p className="font-medium text-foreground text-sm">Production</p>
                <p className="text-muted-foreground text-xs">0 projects · Choose a plan</p>
              </div>
            </div>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setAddServerOpen(true)}
              className="gap-1.5 rounded-xl text-xs"
            >
              <Plus className="size-3.5" />
              New dedicated server
            </Button>
          </div>
          <p className="mt-3 text-[11px] text-muted-foreground">Uses this server&apos;s plan.</p>
        </section>

        <section className="rounded-2xl border border-border/60 bg-card p-6">
          <h2 className="font-semibold text-foreground text-sm">Machine Power</h2>
          <div className="mt-4 grid gap-3 sm:grid-cols-2">
            <button
              type="button"
              onClick={() => setMachinePower('full')}
              className={`flex items-start justify-between rounded-xl border p-4 text-left transition-all ${
                machinePower === 'full'
                  ? 'border-primary/60 bg-primary/5 ring-1 ring-primary/20'
                  : 'border-border/60 bg-muted/15 hover:bg-muted/30'
              }`}
            >
              <div className="flex items-start gap-3">
                <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                  <InfinityIcon className="size-4" />
                </div>
                <div>
                  <p className="font-semibold text-foreground text-xs">Full capacity</p>
                  <p className="mt-1 text-[11px] text-muted-foreground leading-relaxed">
                    Use the server&apos;s available CPU and RAM (default)
                  </p>
                </div>
              </div>
              <span
                className={`mt-0.5 flex size-4 items-center justify-center rounded-full border ${
                  machinePower === 'full'
                    ? 'border-primary bg-primary text-primary-foreground'
                    : 'border-muted-foreground/40'
                }`}
              >
                {machinePower === 'full' && (
                  <span className="size-1.5 rounded-full bg-background" />
                )}
              </span>
            </button>

            <button
              type="button"
              onClick={() => setMachinePower('custom')}
              className={`flex items-start justify-between rounded-xl border p-4 text-left transition-all ${
                machinePower === 'custom'
                  ? 'border-primary/60 bg-primary/5 ring-1 ring-primary/20'
                  : 'border-border/60 bg-muted/15 hover:bg-muted/30'
              }`}
            >
              <div className="flex items-start gap-3">
                <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                  <Sliders className="size-4" />
                </div>
                <div>
                  <p className="font-semibold text-foreground text-xs">Customized</p>
                  <p className="mt-1 text-[11px] text-muted-foreground leading-relaxed">
                    Choose a preset or set your own limits
                  </p>
                </div>
              </div>
              <span
                className={`mt-0.5 flex size-4 items-center justify-center rounded-full border ${
                  machinePower === 'custom'
                    ? 'border-primary bg-primary text-primary-foreground'
                    : 'border-muted-foreground/40'
                }`}
              >
                {machinePower === 'custom' && (
                  <span className="size-1.5 rounded-full bg-background" />
                )}
              </span>
            </button>
          </div>
        </section>

        <section className="rounded-2xl border border-border/60 bg-card p-6">
          <h2 className="font-semibold text-foreground text-sm">Rollback & backups</h2>
          <p className="mt-0.5 text-muted-foreground text-xs">
            How many past versions stay restorable on this server, and what gets backed up before a
            deploy.
          </p>

          <div className="mt-5">
            <RollbackRetentionCards
              strategy={rollbackReuse ? 'snapshot' : 'git'}
              capacity={{
                window: rollbackHistory,
                source: rollbackHistory === 5 ? 'instance-default' : 'explicit',
                explicit: rollbackHistory === 5 ? null : rollbackHistory,
                snapshotSizeBytes: null,
                measuredAt: null,
                diskFreeBytes: null,
                diskTotalBytes: null,
                maxWindow: 20,
                diskBudgetFraction: 0.2,
                strategy: rollbackReuse ? 'snapshot' : 'git',
              }}
              onToggleStrategy={() => setRollbackReuse((prev) => !prev)}
              onChangeWindow={(next) => setRollbackHistory(next ?? 5)}
            />
          </div>
        </section>
      </div>

      <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
        <div className="rounded-2xl border border-border/60 bg-card p-5">
          <h3 className="font-semibold text-foreground text-sm">Deployment preview</h3>
          <div className="mt-3">
            <p className="text-muted-foreground text-xs">Destination</p>
            <p className="font-semibold text-foreground text-sm">Production</p>
          </div>

          <div className="mt-4 space-y-3 rounded-xl border border-border/50 bg-muted/20 p-3.5 text-xs">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Runtime</span>
              <span className="font-medium text-foreground">Sandboxed</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Machine Power</span>
              <span className="font-medium text-foreground">
                {machinePower === 'full' ? 'Full capacity' : 'Customized'}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Rollback</span>
              <span className="font-medium text-foreground">auto</span>
            </div>
          </div>

          <Button className="mt-6 w-full gap-2 font-semibold" size="lg" onClick={onContinue}>
            Continue
            <ArrowRight className="size-4" />
          </Button>
          <p className="mt-2 text-center text-[11px] text-muted-foreground leading-relaxed">
            Continue to configure your app. Deployment starts from the next page.
          </p>
        </div>
      </div>

      <AddServerModal open={addServerOpen} onClose={() => setAddServerOpen(false)} />
    </div>
  );
}
