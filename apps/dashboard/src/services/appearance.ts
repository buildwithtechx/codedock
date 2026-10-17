import type { AppearancePreferences, CustomTheme, ThemeMode } from '#/interfaces/appearance';
import { apiClient } from '#/lib/api-client';

export const APPEARANCE_CACHE_POINTER = 'codedock-appearance-owner';
export const APPEARANCE_CHANGED = 'codedock:appearance-changed';

export function appearanceCacheKey(owner: string): string {
  return `codedock-appearance-v1:${owner}`;
}

export function cachedAppearance(key: string): AppearancePreferences | null {
  try {
    const raw = localStorage.getItem(key);
    return raw ? (JSON.parse(raw) as AppearancePreferences) : null;
  } catch {
    return null;
  }
}

export function resolveAppearanceMode(
  preferences: AppearancePreferences,
  prefersDark = false
): 'light' | 'dim' | 'dark' {
  if (preferences.mode === 'system') {
    return prefersDark ? 'dark' : 'light';
  }
  if (preferences.mode === 'dim' || preferences.mode === 'dark' || preferences.mode === 'light') {
    return preferences.mode;
  }
  return 'dark';
}

const APPLIED_VARIABLES = [
  '--radius',
  '--font-sans',
  '--font-mono',
  '--appearance-font-size',
  '--appearance-mono-size',
  '--appearance-line-height',
  '--appearance-font-weight',
];

export function applyAppearance(
  preferences: AppearancePreferences,
  mode?: 'light' | 'dim' | 'dark'
): void {
  if (typeof document === 'undefined') return;

  const dark =
    typeof window !== 'undefined' && window.matchMedia?.('(prefers-color-scheme: dark)').matches;
  const resolvedMode = mode ?? resolveAppearanceMode(preferences, !!dark);
  const root = document.documentElement;

  root.setAttribute('data-theme', resolvedMode);
  root.classList.remove('light', 'dim', 'dark');
  root.classList.add(resolvedMode);

  for (const key of APPLIED_VARIABLES) {
    root.style.removeProperty(key);
  }

  if (preferences.radius) {
    root.style.setProperty('--radius', preferences.radius);
  }
  if (preferences.fontFamily) {
    root.style.setProperty('--font-sans', preferences.fontFamily);
  }

  const customTheme = preferences.customThemes?.find((t) => t.id === preferences.activeTheme);
  const allTokens = {
    ...(preferences.tokens ?? {}),
    ...(customTheme?.tokens ?? {}),
  };

  for (const [key, value] of Object.entries(allTokens)) {
    if (typeof value === 'string' && !/[;{}<>\\]|url\s*\(/i.test(value)) {
      root.style.setProperty(key, value);
    }
  }

  if (typeof window !== 'undefined') {
    window.dispatchEvent(new Event(APPEARANCE_CHANGED));
  }
}

export function cacheAppearance(
  key: string,
  preferences: AppearancePreferences,
  mode: 'light' | 'dim' | 'dark'
): void {
  try {
    localStorage.setItem(key, JSON.stringify(preferences));
    localStorage.setItem(
      `${key}:paint`,
      JSON.stringify({
        mode,
        preferences,
      })
    );
    localStorage.setItem(APPEARANCE_CACHE_POINTER, key);
  } catch {
    // Storage unavailable or quota exceeded
  }
}

export function bootAppearance(): void {
  try {
    if (typeof window === 'undefined' || typeof document === 'undefined') return;
    const dark = window.matchMedia?.('(prefers-color-scheme: dark)').matches;
    const owner = localStorage.getItem(APPEARANCE_CACHE_POINTER);
    const raw = owner?.startsWith('codedock-appearance-v1:')
      ? localStorage.getItem(`${owner}:paint`)
      : null;
    const snapshot = raw && raw.length < 64_000 ? JSON.parse(raw) : null;

    let mode: 'light' | 'dim' | 'dark' = dark ? 'dark' : 'light';
    if (snapshot?.mode && ['light', 'dim', 'dark'].includes(snapshot.mode)) {
      mode = snapshot.mode;
    }

    const root = document.documentElement;
    root.setAttribute('data-theme', mode);
    root.classList.remove('light', 'dim', 'dark');
    root.classList.add(mode);

    if (snapshot?.preferences) {
      applyAppearance(snapshot.preferences, mode);
    }
  } catch {
    document.documentElement.setAttribute('data-theme', 'dark');
    document.documentElement.classList.add('dark');
  }
}

export async function getAppearance(): Promise<AppearancePreferences> {
  const response = await apiClient.get<{ appearance: AppearancePreferences }>(
    '/api/settings/appearance'
  );
  return response.appearance;
}

export async function setAppearance(
  appearance: AppearancePreferences
): Promise<AppearancePreferences> {
  const response = await apiClient.put<{ appearance: AppearancePreferences }>(
    '/api/settings/appearance',
    { appearance }
  );
  return response.appearance;
}

export type { AppearancePreferences, CustomTheme, ThemeMode };
