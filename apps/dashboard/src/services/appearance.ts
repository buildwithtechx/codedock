import { apiClient } from '#/lib/api-client';

export interface CustomTheme {
  id: string;
  name: string;
  mode?: 'light' | 'dim' | 'dark';
  tokens?: Record<string, string>;
}

export interface AppearancePreferences {
  activeTheme: string;
  mode: 'light' | 'dim' | 'dark' | 'system';
  fontFamily?: string;
  radius?: string;
  tokens?: Record<string, string>;
  customThemes?: CustomTheme[];
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
