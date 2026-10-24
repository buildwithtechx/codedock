import { Button } from '@codedock/ui/button';
import { Moon, Plus, Sun, SunMedium, Tv } from 'lucide-react';
import type { ThemeMode } from '#/interfaces/appearance';

const MODES: { id: ThemeMode; label: string; icon: typeof Sun }[] = [
  { id: 'light', label: 'Light', icon: Sun },
  { id: 'dim', label: 'Dim', icon: SunMedium },
  { id: 'dark', label: 'Dark', icon: Moon },
  { id: 'system', label: 'System', icon: Tv },
];

interface AppearanceModeStripProps {
  currentMode: ThemeMode;
  onSelectMode: (mode: ThemeMode) => void;
  onNewTheme: () => void;
}

export function AppearanceModeStrip({
  currentMode,
  onSelectMode,
  onNewTheme,
}: AppearanceModeStripProps) {
  return (
    <section className="rounded-2xl border border-border/60 bg-card p-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div
          role="group"
          aria-label="Presets"
          className="grid w-full grid-cols-2 gap-1 rounded-xl bg-background/60 p-1 sm:flex sm:w-auto sm:flex-wrap"
        >
          {MODES.map((m) => {
            const Icon = m.icon;
            const isSelected = currentMode === m.id;
            return (
              <button
                key={m.id}
                type="button"
                aria-pressed={isSelected}
                onClick={() => onSelectMode(m.id)}
                className={`flex h-9 items-center gap-2 rounded-lg px-3.5 font-medium text-xs transition-colors ${
                  isSelected
                    ? 'bg-foreground/8 text-foreground hover:bg-foreground/8'
                    : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
                }`}
              >
                <Icon className="size-4 shrink-0" />
                <span>{m.label}</span>
              </button>
            );
          })}
        </div>

        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={onNewTheme}
          className="gap-1.5 text-muted-foreground text-xs hover:text-foreground"
        >
          <Plus className="size-3.5" />
          <span>New theme</span>
        </Button>
      </div>
    </section>
  );
}
