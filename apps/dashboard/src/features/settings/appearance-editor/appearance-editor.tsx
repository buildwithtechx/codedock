import { useTheme } from '@codedock/ui';
import { Button } from '@codedock/ui/button';
import { Link } from '@tanstack/react-router';
import { ArrowDownToLine, ArrowLeft, ArrowUpToLine } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
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
import { AppearanceModeStrip } from './appearance-mode-strip';
import { getEffectiveThemeTokens } from './appearance-tokens-data';
import { DesignTokensSection } from './design-tokens-section';
import { ThemeImportModal } from './theme-import-modal';
import { TypographySection } from './typography-section';

export function AppearanceEditor() {
  const { theme, setTheme } = useTheme();
  const user = useAuthStore((s) => s.user);
  const owner = user?.id || 'default';
  const cacheKey = appearanceCacheKey(owner);

  const [saved, setSaved] = useState<AppearancePreferences>(
    () =>
      cachedAppearance(cacheKey) ?? {
        activeTheme: 'neutral-slate',
        mode: (theme as ThemeMode) || 'dark',
        radius: '0.5rem',
        tokens: {},
        customThemes: [],
      }
  );

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
    const mode = draft.mode === 'system' ? 'dark' : (draft.mode as 'light' | 'dim' | 'dark');
    downloadThemeFile({
      version: 1,
      id: draft.activeTheme || `theme-${Date.now()}`,
      name: draft.activeTheme || `${mode.charAt(0).toUpperCase() + mode.slice(1)} Custom`,
      base: mode,
      mode,
      radius: draft.radius || '0.5rem',
      tokens: getEffectiveThemeTokens(mode, draft.tokens ?? {}),
      typography: {
        sans: draft.fontFamily || 'system',
        mono: 'mono',
        size: 14,
        monoSize: 13,
        lineHeight: 1.5,
        weight: 400,
      },
    });
    toast.success('Theme JSON downloaded with full tokens');
  };

  return (
    <div className="space-y-6 pb-6">
      {/* Back button */}
      <div>
        <Link
          to="/settings"
          search={{ tab: 'general' } as never}
          className="inline-flex items-center gap-1.5 text-muted-foreground text-xs transition-colors hover:text-foreground"
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
            className="gap-1.5 border-border/60 text-xs"
          >
            <ArrowUpToLine className="size-3.5" />
            <span>Import</span>
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleDownloadTheme}
            className="gap-1.5 border-border/60 text-xs"
          >
            <ArrowDownToLine className="size-3.5" />
            <span>Download theme</span>
          </Button>
        </div>
      </div>

      {/* Theme Mode Strip Card */}
      <AppearanceModeStrip
        currentMode={draft.mode}
        onSelectMode={handleModeChange}
        onNewTheme={handleCreateTheme}
      />

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
          updateDraft({
            ...draft,
            activeTheme: imported.id,
            customThemes: [...(draft.customThemes ?? []), imported],
            tokens: { ...(draft.tokens ?? {}), ...(imported.tokens ?? {}) },
          });
        }}
      />
    </div>
  );
}
