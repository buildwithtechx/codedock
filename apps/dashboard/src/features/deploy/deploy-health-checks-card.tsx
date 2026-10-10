import { Activity, ChevronDown, ChevronUp } from 'lucide-react';
import { useState } from 'react';
import { Switch } from '#/components/ui/switch';

interface DeployHealthChecksCardProps {
  waitForPort?: boolean;
  onWaitForPortChange?: (v: boolean) => void;
  watchCrashLoop?: boolean;
  onWatchCrashLoopChange?: (v: boolean) => void;
}

export function DeployHealthChecksCard({
  waitForPort = false,
  onWaitForPortChange,
  watchCrashLoop = false,
  onWatchCrashLoopChange,
}: DeployHealthChecksCardProps) {
  const [open, setOpen] = useState(false);
  const [waitPort, setWaitPort] = useState(waitForPort);
  const [crashLoop, setCrashLoop] = useState(watchCrashLoop);

  const handleWaitChange = (val: boolean) => {
    setWaitPort(val);
    onWaitForPortChange?.(val);
  };

  const handleCrashChange = (val: boolean) => {
    setCrashLoop(val);
    onWatchCrashLoopChange?.(val);
  };

  const isEnabled = waitPort || crashLoop;

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3.5">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted text-muted-foreground">
            <Activity className="size-5" />
          </div>
          <div>
            <h2 className="font-semibold text-foreground text-sm">Health checks</h2>
            <p className="text-muted-foreground text-xs">
              {isEnabled
                ? 'On — health probes active during rollout'
                : 'Off — the deploy finishes as soon as your app is running'}
            </p>
          </div>
        </div>

        <button
          type="button"
          onClick={() => setOpen((prev) => !prev)}
          className="p-1 text-muted-foreground transition-colors hover:text-foreground"
          aria-label={open ? 'Collapse health checks' : 'Expand health checks'}
        >
          {open ? <ChevronUp className="size-4" /> : <ChevronDown className="size-4" />}
        </button>
      </div>

      {open && (
        <div className="mt-5 space-y-5 border-border/50 border-t pt-5">
          <p className="text-muted-foreground text-xs leading-relaxed">
            Codedock doesn't wait for your app after it starts, so a slow boot never delays or fails
            a deploy. Once it's live, the port is checked anyway and reported without blocking
            anything. Turn these on only if you want the deploy itself to refuse to go green.
          </p>

          <div className="space-y-4">
            <div className="flex items-center justify-between gap-4 rounded-xl border border-border/40 bg-muted/15 p-4">
              <div className="space-y-0.5">
                <p className="font-medium text-foreground text-sm">Wait for the app to answer</p>
                <p className="text-muted-foreground text-xs">
                  After the container starts, dial the app's port and wait for it to accept a
                  connection.
                </p>
              </div>
              <Switch checked={waitPort} onCheckedChange={handleWaitChange} />
            </div>

            <div className="flex items-center justify-between gap-4 rounded-xl border border-border/40 bg-muted/15 p-4">
              <div className="space-y-0.5">
                <p className="font-medium text-foreground text-sm">Watch for restart loops</p>
                <p className="text-muted-foreground text-xs">
                  Watch the container after it starts and report a crash loop. Reads the runtime, so
                  this works on remote servers too.
                </p>
              </div>
              <Switch checked={crashLoop} onCheckedChange={handleCrashChange} />
            </div>
          </div>
        </div>
      )}
    </section>
  );
}
