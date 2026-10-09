import {
  ArrowRight,
  Cloud,
  HardDrive,
  History,
  Minus,
  Plus,
  RotateCcw,
  Server,
} from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Switch } from '#/components/ui/switch';

interface DeployTargetStepProps {
  onContinue: () => void;
  runtimeMode: 'web' | 'worker' | 'static';
  onRuntimeModeChange: (mode: 'web' | 'worker' | 'static') => void;
}

export function DeployTargetStep({
  onContinue,
  runtimeMode,
  onRuntimeModeChange,
}: DeployTargetStepProps) {
  const [staticHosting, setStaticHosting] = useState<'pages' | 'server'>('pages');
  const [rollbackReuse, setRollbackReuse] = useState(false);
  const [rollbackHistory, setRollbackHistory] = useState(5);

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_340px]">
      <div className="space-y-6">
        <div>
          <h1 className="font-semibold text-2xl text-foreground tracking-tight">
            Where do you want to deploy?
          </h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Choose where your application will run
          </p>
        </div>

        <section className="rounded-2xl border border-border/60 bg-card p-6">
          <div className="flex items-center justify-between">
            <h2 className="font-semibold text-foreground text-sm">Destination</h2>
          </div>
          <div className="mt-4 flex flex-col gap-3 rounded-xl border border-border/60 bg-muted/20 p-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-3">
              <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Server className="size-4" />
              </div>
              <div>
                <p className="font-medium text-foreground text-sm">Production</p>
                <p className="text-muted-foreground text-xs">
                  Primary control plane · Default capacity
                </p>
              </div>
            </div>
            <Button variant="outline" size="sm" className="gap-1.5 text-xs">
              <Plus className="size-3.5" />
              New dedicated server
            </Button>
          </div>
          <p className="mt-3 text-[11px] text-muted-foreground/70">
            Uses this server&apos;s allocated resources.
          </p>
        </section>

        <section className="rounded-2xl border border-border/60 bg-card p-6">
          <h2 className="font-semibold text-foreground text-sm">Hosting & Runtime Strategy</h2>
          <p className="mt-0.5 text-muted-foreground text-xs">
            Choose where the built files and processes are served.
          </p>
          <div className="mt-4 grid gap-3 sm:grid-cols-2">
            <button
              type="button"
              onClick={() => {
                setStaticHosting('pages');
                onRuntimeModeChange('static');
              }}
              className={`flex items-start gap-3 rounded-xl border p-4 text-left transition-all ${
                staticHosting === 'pages'
                  ? 'border-primary/60 bg-primary/5 ring-1 ring-primary/20'
                  : 'border-border/60 bg-muted/15 hover:bg-muted/30'
              }`}
            >
              <div className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Cloud className="size-4" />
              </div>
              <div className="min-w-0">
                <p className="font-semibold text-foreground text-xs">Cloud Pages (default)</p>
                <p className="mt-1 text-[11px] text-muted-foreground leading-relaxed">
                  Publish files to the Cloud edge. No running app container.
                </p>
              </div>
            </button>

            <button
              type="button"
              onClick={() => {
                setStaticHosting('server');
                onRuntimeModeChange('web');
              }}
              className={`flex items-start gap-3 rounded-xl border p-4 text-left transition-all ${
                staticHosting === 'server'
                  ? 'border-primary/60 bg-primary/5 ring-1 ring-primary/20'
                  : 'border-border/60 bg-muted/15 hover:bg-muted/30'
              }`}
            >
              <div className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                <HardDrive className="size-4" />
              </div>
              <div className="min-w-0">
                <p className="font-semibold text-foreground text-xs">Managed server</p>
                <p className="mt-1 text-[11px] text-muted-foreground leading-relaxed">
                  Serve files from an nginx container on your managed server.
                </p>
              </div>
            </button>
          </div>
        </section>

        <section className="rounded-2xl border border-border/60 bg-card p-6">
          <h2 className="font-semibold text-foreground text-sm">Rollback & backups</h2>
          <p className="mt-0.5 text-muted-foreground text-xs">
            How many past versions stay restorable on this server, and what gets backed up before a
            deploy.
          </p>

          <div className="mt-4 space-y-4">
            <div className="flex items-center justify-between gap-4 rounded-xl border border-border/50 bg-muted/20 p-4">
              <div className="flex items-start gap-3">
                <RotateCcw className="mt-0.5 size-4 text-muted-foreground" />
                <div>
                  <p className="font-medium text-foreground text-sm">Rollback strategy</p>
                  <p className="font-semibold text-foreground/80 text-xs">
                    Reuse images or rebuild
                  </p>
                  <p className="mt-0.5 text-[11px] text-muted-foreground">
                    Rebuilds the old commit, unless its files are still on the host.
                  </p>
                </div>
              </div>
              <Switch checked={rollbackReuse} onCheckedChange={setRollbackReuse} />
            </div>

            <div className="flex items-center justify-between gap-4 rounded-xl border border-border/50 bg-muted/20 p-4">
              <div className="flex items-start gap-3">
                <History className="mt-0.5 size-4 text-muted-foreground" />
                <div>
                  <p className="font-medium text-foreground text-sm">Rollback history</p>
                  <p className="font-semibold text-foreground/80 text-xs">
                    {rollbackHistory} versions
                  </p>
                  <p className="mt-0.5 max-w-md text-[11px] text-muted-foreground leading-relaxed">
                    Maximum number of past releases to keep. Active and pinned releases are extra.
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="icon"
                  className="size-8"
                  onClick={() => setRollbackHistory((prev) => Math.max(1, prev - 1))}
                  disabled={rollbackHistory <= 1}
                >
                  <Minus className="size-3" />
                </Button>
                <span className="w-6 text-center font-semibold text-xs tabular-nums">
                  {rollbackHistory}
                </span>
                <Button
                  variant="outline"
                  size="icon"
                  className="size-8"
                  onClick={() => setRollbackHistory((prev) => Math.min(20, prev + 1))}
                  disabled={rollbackHistory >= 20}
                >
                  <Plus className="size-3" />
                </Button>
              </div>
            </div>
          </div>
          <p className="mt-3 text-[11px] text-muted-foreground/70">
            built files on disk · instance default · Editable in project settings after first
            deploy.
          </p>
        </section>
      </div>

      <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
        <div className="rounded-2xl border border-border/60 bg-card p-5">
          <h3 className="font-semibold text-foreground text-sm">Deployment preview</h3>
          <p className="mt-1 text-muted-foreground text-xs">Build & deploy:</p>
          <p className="font-semibold text-foreground text-sm">Production</p>

          <div className="mt-5 space-y-3 border-border/60 border-t pt-4 text-xs">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Runtime</span>
              <span className="font-medium text-foreground capitalize">{runtimeMode}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Machine Power</span>
              <span className="font-medium text-foreground">Full capacity</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Rollback</span>
              <span className="font-medium text-foreground">auto</span>
            </div>
          </div>

          <Button className="mt-6 w-full gap-2" size="lg" onClick={onContinue}>
            Continue
            <ArrowRight className="size-4" />
          </Button>
          <p className="mt-2 text-center text-[11px] text-muted-foreground">
            Continue to configure your app. Deployment starts from the next page.
          </p>
        </div>
      </div>
    </div>
  );
}
