import { Check, SunMedium } from 'lucide-react';
import { useTheme } from 'next-themes';
import { useId } from 'react';
import { ThemePreview } from './theme-preview';

const THEMES = [
  { id: 'light', label: 'Light' },
  { id: 'dim', label: 'Dim' },
  { id: 'dark', label: 'Dark' },
  { id: 'system', label: 'System' },
] as const;

export function AppearanceSetting() {
  const { theme, setTheme } = useTheme();
  const groupName = useId();

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="mb-5 flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <SunMedium className="h-4 w-4" />
        </div>
        <div>
          <h2 className="font-semibold text-foreground text-sm">Appearance</h2>
          <p className="text-muted-foreground text-xs">
            Choose a theme, or follow your device&apos;s appearance. Changes apply immediately.
          </p>
        </div>
      </div>

      <fieldset className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <legend className="sr-only">Appearance themes</legend>
        {THEMES.map(({ id, label }) => {
          const selected = theme === id;
          return (
            <label key={id} className="min-w-0 cursor-pointer">
              <input
                type="radio"
                name={groupName}
                value={id}
                checked={selected}
                onChange={() => setTheme(id)}
                className="peer sr-only"
              />
              <div
                className={`rounded-xl p-2.5 transition-colors peer-focus-visible:ring-2 peer-focus-visible:ring-ring ${
                  selected
                    ? 'border border-primary/50 bg-primary/5'
                    : 'border border-border/50 bg-muted/20 hover:bg-muted/40'
                }`}
              >
                {id === 'system' ? (
                  <div className="relative overflow-hidden rounded-lg" aria-hidden="true">
                    <ThemePreview theme="light" />
                    <div className="absolute inset-0 [clip-path:inset(0_0_0_50%)]">
                      <ThemePreview theme="dark" />
                    </div>
                  </div>
                ) : (
                  <ThemePreview theme={id === 'dim' ? 'dim' : id === 'light' ? 'light' : 'dark'} />
                )}
                <div className="mt-3 flex min-w-0 items-center justify-between gap-2 px-1 font-medium text-foreground text-sm">
                  <span>{label}</span>
                  <Check
                    className={`h-4 w-4 shrink-0 text-primary ${selected ? 'visible' : 'invisible'}`}
                  />
                </div>
              </div>
            </label>
          );
        })}
      </fieldset>
    </section>
  );
}
