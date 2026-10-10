import { Sliders } from 'lucide-react';
import { Switch } from '#/components/ui/switch';
import { useDemoMode } from '#/lib/demo-mode';

export function PreferencesSetting() {
  const [demoMode, setDemoMode] = useDemoMode();

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5">
      <div className="mb-5 flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <Sliders className="h-4 w-4" />
        </div>
        <div>
          <h2 className="font-semibold text-foreground text-sm">Preferences</h2>
          <p className="text-muted-foreground text-xs">
            Personal display preferences for this browser — nothing is sent to the server.
          </p>
        </div>
      </div>

      <div className="flex items-center justify-between gap-4 rounded-xl border border-border/50 bg-muted/20 p-4">
        <div className="min-w-0">
          <p className="font-medium text-foreground text-sm">Demo mode</p>
          <p className="text-muted-foreground text-xs">
            Blur IP addresses, hosts, and environment values during screen-shares and recordings.
          </p>
        </div>
        <Switch checked={demoMode} onCheckedChange={setDemoMode} aria-label="Toggle demo mode" />
      </div>
    </section>
  );
}
