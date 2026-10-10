import { Check, Languages } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { SettingsSection } from './settings-section';

const LANGUAGES = [
  { code: 'en', name: 'English', badge: 'EN', dir: 'ltr' },
  { code: 'ar', name: 'العربية', badge: 'ع', dir: 'rtl' },
  { code: 'es', name: 'Español', badge: 'ES', dir: 'ltr' },
  { code: 'fr', name: 'Français', badge: 'FR', dir: 'ltr' },
  { code: 'de', name: 'Deutsch', badge: 'DE', dir: 'ltr' },
  { code: 'pt', name: 'Português', badge: 'PT', dir: 'ltr' },
  { code: 'ja', name: '日本語', badge: '日', dir: 'ltr' },
  { code: 'zh', name: '中文', badge: '中', dir: 'ltr' },
  { code: 'tr', name: 'Türkçe', badge: 'TR', dir: 'ltr' },
] as const;

type LocaleCode = (typeof LANGUAGES)[number]['code'];

export function LanguageSetting() {
  const [locale, setLocale] = useState<LocaleCode>('en');

  useEffect(() => {
    const saved = localStorage.getItem('codedock_locale') as LocaleCode | null;
    if (saved && LANGUAGES.some((l) => l.code === saved)) {
      setLocale(saved);
    }
  }, []);

  const handleSelect = (lang: (typeof LANGUAGES)[number]) => {
    setLocale(lang.code);
    localStorage.setItem('codedock_locale', lang.code);
    if (typeof document !== 'undefined') {
      document.documentElement.lang = lang.code;
      document.documentElement.dir = lang.dir;
    }
    toast.success(`Language set to ${lang.name}`);
  };

  return (
    <SettingsSection
      icon={<Languages className="size-4 text-primary" />}
      title="Language"
      description="Select your preferred language for the dashboard interface."
      collapsible
    >
      <div className="grid gap-3 sm:grid-cols-2">
        {LANGUAGES.map((l) => {
          const active = l.code === locale;
          return (
            <button
              key={l.code}
              type="button"
              onClick={() => handleSelect(l)}
              aria-pressed={active}
              lang={l.code}
              dir={l.dir}
              className={`group relative flex items-center gap-3 rounded-xl border p-4 text-start transition-all ${
                active
                  ? 'border-primary/50 bg-primary/[0.06] ring-1 ring-primary/25'
                  : 'border-border/50 bg-muted/10 hover:border-border hover:bg-muted/30'
              }`}
            >
              <div
                className={`flex size-11 shrink-0 items-center justify-center rounded-lg font-semibold text-[15px] transition-colors ${
                  active
                    ? 'bg-primary/15 text-primary'
                    : 'bg-muted text-muted-foreground group-hover:text-foreground'
                }`}
              >
                {l.badge}
              </div>
              <div className="min-w-0 flex-1">
                <p className="font-medium text-foreground text-sm">{l.name}</p>
              </div>
              {active && (
                <span className="flex size-5 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground">
                  <Check className="size-3.5" />
                </span>
              )}
            </button>
          );
        })}
      </div>
    </SettingsSection>
  );
}
