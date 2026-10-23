import { Link } from '@tanstack/react-router';
import { Check, SunMedium } from 'lucide-react';
import { useId } from 'react';
import { Button } from '@codedock/ui/button';
import { useTheme } from '@codedock/ui';
import type { AppearancePreferences, ThemeMode } from '#/interfaces/appearance';
import {
  appearanceCacheKey,
  applyAppearance,
  cacheAppearance,
  cachedAppearance,
  setAppearance,
} from '#/services/appearance';
import { useAuthStore } from '#/stores/auth-store';
import { SettingsSection } from './settings-section';
import { ThemePreview } from './theme-preview';

const BASE_THEMES = [
  { id: 'light' as const, label: 'Light' },
  { id: 'dim' as const, label: 'Dim' },
  { id: 'dark' as const, label: 'Dark' },
  { id: 'system' as const, label: 'System' },
];

export function AppearanceSetting() {
  const { theme, setTheme } = useTheme();
  const groupName = useId();
  const user = useAuthStore((s) => s.user);

  const owner = user?.id || 'default';
  const cacheKey = appearanceCacheKey(owner);

  const handleBaseModeChange = (id: ThemeMode) => {
    setTheme(id);
    const cached = cachedAppearance(cacheKey);
    const isDark =
      typeof window !== 'undefined' && window.matchMedia?.('(prefers-color-scheme: dark)').matches;
    const resolvedMode = id === 'system' ? (isDark ? 'dim' : 'light') : id;
    const next: AppearancePreferences = {
      ...(cached ?? { activeTheme: id, tokens: {} }),
      activeTheme: id,
      tokens: {},
      mode: id,
    };
    applyAppearance(next, resolvedMode);
    cacheAppearance(cacheKey, next, resolvedMode);
    void setAppearance(next).catch(() => {});
  };

  return (
    <SettingsSection
      icon={<SunMedium className="size-4 text-primary" />}
      title="Appearance"
      description="Choose a theme, or follow your device's appearance. Changes apply immediately."
      action={
        <Button asChild variant="secondary" size="sm">
          <Link to="/settings" search={{ tab: 'general', view: 'appearance' } as never}>
            Customize
          </Link>
        </Button>
      }
    >
      <fieldset className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <legend className="sr-only">Appearance themes</legend>
        {BASE_THEMES.map(({ id, label }) => {
          const selected = theme === id;
          return (
            <label key={id} className="min-w-0 cursor-pointer">
              <input
                type="radio"
                name={groupName}
                value={id}
                checked={selected}
                onChange={() => handleBaseModeChange(id)}
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
                      <ThemePreview theme="dim" />
                    </div>
                  </div>
                ) : (
                  <ThemePreview theme={id} />
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
    </SettingsSection>
  );
}
