import { Link } from '@tanstack/react-router';
import {
  ArrowDownToLine,
  ArrowLeft,
  ArrowUpToLine,
  Moon,
  Plus,
  Sun,
  SunMedium,
  Tv,
} from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { useTheme } from '@codedock/ui';
import type { AppearancePreferences, CustomTheme, ThemeMode } from '#/interfaces/appearance';
import { downloadThemeFile } from '#/lib/theme-export-import';
import {
  appearanceCacheKey,
  applyAppearance,
  cacheAppearance,
  cachedAppearance,
  getAppearance,
  setAppearance,
} from '#/services/appearance';
import { useAuthStore } from '#/stores/auth-store';
import { AppearanceFloatingBar } from './appearance-floating-bar';
import { AppearanceLivePreview } from './appearance-live-preview';
import { DesignTokensSection } from './design-tokens-section';
import { ThemeImportModal } from './theme-import-modal';
import { TypographySection } from './typography-section';

const MODES: { id: ThemeMode; label: string; icon: typeof Sun }[] = [
  { id: 'light', label: 'Light', icon: Sun },
  { id: 'dim', label: 'Dim', icon: SunMedium },
  { id: 'dark', label: 'Dark', icon: Moon },
  { id: 'system', label: 'System', icon: Tv },
];

export function AppearanceEditor() {
  const { theme, setTheme } = useTheme();
  const user = useAuthStore((s) => s.user);
  const owner = user?.id || 'default';
  const cacheKey = appearanceCacheKey(owner);

  const [saved, setSaved] = useState<AppearancePreferences>(() => {
    return (
      cachedAppearance(cacheKey) ?? {
        activeTheme: 'neutral-slate',
        mode: (theme as ThemeMode) || 'dark',
        radius: '0.5rem',
        tokens: {},
        customThemes: [],
      }
    );
  });

  const [draft, setDraft] = useState<AppearancePreferences>(saved);
  const [importModalOpen, setImportModalOpen] = useState(false);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    void getAppearance()
      .then((remote) => {
        if (remote) {
          setSaved(remote);
          setDraft(remote);
        }
      })
      .catch(() => {});
  }, []);

  const isDirty = JSON.stringify(saved) !== JSON.stringify(draft);

  const updateDraft = (next: AppearancePreferences) => {
    setDraft(next);
    const resolvedMode = next.mode === 'system' ? 'dark' : (next.mode as 'light' | 'dim' | 'dark');
    applyAppearance(next, resolvedMode);
  };

  const handleModeChange = (mode: ThemeMode) => {
    setTheme(mode);
    updateDraft({ ...draft, mode });
  };

  const handleSetToken = (key: string, value: string) => {
    updateDraft({ ...draft, tokens: { ...(draft.tokens ?? {}), [key]: value } });
  };

  const handleResetToken = (key: string) => {
    const tokens = { ...(draft.tokens ?? {}) };
    delete tokens[key];
    updateDraft({ ...draft, tokens });
  };

  const handleResetAll = () => {
    updateDraft({ ...draft, tokens: {} });
    toast.info('Design tokens reset to defaults');
  };

  const handleCreateTheme = () => {
    const newTheme: CustomTheme = {
      id: `theme-${Date.now()}`,
      name: `Custom ${(draft.customThemes?.length || 0) + 1}`,
      mode: draft.mode === 'system' ? 'dark' : draft.mode,
      tokens: { ...(draft.tokens ?? {}) },
    };
    updateDraft({
      ...draft,
      activeTheme: newTheme.id,
      customThemes: [...(draft.customThemes ?? []), newTheme],
    });
    toast.success(`Created "${newTheme.name}"`);
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      const mode = draft.mode === 'system' ? 'dark' : (draft.mode as 'light' | 'dim' | 'dark');
      cacheAppearance(cacheKey, draft, mode);
      await setAppearance(draft);
      setSaved(draft);
      toast.success('Appearance & theme saved');
    } catch {
      toast.error('Failed to save appearance to server');
    } finally {
      setSaving(false);
    }
  };

  const handleDiscard = () => {
    setDraft(saved);
    const mode = saved.mode === 'system' ? 'dark' : (saved.mode as 'light' | 'dim' | 'dark');
    applyAppearance(saved, mode);
    toast.info('Draft changes discarded');
  };

  const handleDownloadTheme = () => {
    downloadThemeFile({
      version: 1,
      name: draft.activeTheme || 'Codedock Custom',
      mode: draft.mode === 'system' ? 'dark' : draft.mode,
      radius: draft.radius || '0.5rem',
      tokens: draft.tokens ?? {},
    });
    toast.success('Theme JSON downloaded');
  };

  return (
    <div className="space-y-6 pb-6">
      {/* Back button */}
      <div>
        <Link
          to="/settings"
          search={{ tab: 'general' } as never}
          className="inline-flex items-center gap-1.5 text-xs text-muted-foreground transition-colors hover:text-foreground"
        >
          <ArrowLeft className="size-3.5" />
          <span>Back to General</span>
        </Link>
      </div>

      {/* Header Bar */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="font-semibold text-foreground text-xl tracking-tight">Appearance</h2>
          <p className="mt-1 text-muted-foreground text-xs">
            Personal themes, live previews and portable JSON.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => setImportModalOpen(true)}
            className="gap-1.5 text-xs border-border/60"
          >
            <ArrowUpToLine className="size-3.5" />
            <span>Import</span>
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleDownloadTheme}
            className="gap-1.5 text-xs border-border/60"
          >
            <ArrowDownToLine className="size-3.5" />
            <span>Download theme</span>
          </Button>
        </div>
      </div>

      {/* Theme Mode Strip Card */}
      <section className="rounded-2xl border border-border/60 bg-card p-4">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div
            role="group"
            aria-label="Presets"
            className="grid w-full grid-cols-2 gap-1 rounded-xl bg-background/60 p-1 sm:flex sm:w-auto sm:flex-wrap"
          >
            {MODES.map((m) => {
              const Icon = m.icon;
              const isSelected = draft.mode === m.id;
              return (
                <button
                  key={m.id}
                  type="button"
                  aria-pressed={isSelected}
                  onClick={() => handleModeChange(m.id)}
                  className={`flex h-9 items-center gap-2 rounded-lg px-3.5 font-medium text-xs transition-colors ${
                    isSelected
                      ? 'bg-foreground/[0.08] text-foreground hover:bg-foreground/[0.08]'
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
            onClick={handleCreateTheme}
            className="gap-1.5 text-xs text-muted-foreground hover:text-foreground"
          >
            <Plus className="size-3.5" />
            <span>New theme</span>
          </Button>
        </div>
      </section>

      {/* Two Column Layout (Tokens + Typography Left, Live Preview Right) */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_360px]">
        {/* Left Column */}
        <div className="space-y-6">
          <DesignTokensSection
            tokens={draft.tokens ?? {}}
            onSetToken={handleSetToken}
            onResetToken={handleResetToken}
          />

          <TypographySection
            interfaceFont={draft.fontFamily}
            onChange={(fields) => {
              if (fields.interfaceFont) {
                updateDraft({ ...draft, fontFamily: fields.interfaceFont });
              }
            }}
          />
        </div>

        {/* Right Column (Sticky Preview) */}
        <div>
          <AppearanceLivePreview
            tokens={draft.tokens ?? {}}
            activeThemeName={user?.name || 'codedock'}
          />
        </div>
      </div>

      {/* Floating Bottom Bar */}
      <AppearanceFloatingBar
        isDirty={isDirty}
        isSaving={saving}
        onReset={handleResetAll}
        onDiscard={handleDiscard}
        onSave={handleSave}
      />

      {/* Import Modal */}
      <ThemeImportModal
        open={importModalOpen}
        onOpenChange={setImportModalOpen}
        onImport={(imported) => {
          const next: AppearancePreferences = {
            ...draft,
            activeTheme: imported.id,
            customThemes: [...(draft.customThemes ?? []), imported],
            tokens: { ...(draft.tokens ?? {}), ...(imported.tokens ?? {}) },
          };
          updateDraft(next);
        }}
      />
    </div>
  );
}
