export type ThemeMode = 'light' | 'dark' | 'dim' | 'system';

export interface CustomTheme {
  id: string;
  name: string;
  mode?: 'light' | 'dim' | 'dark';
  tokens?: Record<string, string>;
}

export interface AppearancePreferences {
  activeTheme: string;
  mode: ThemeMode;
  fontFamily?: string;
  radius?: string;
  codeTheme?: string;
  compactMode?: boolean;
  tokens?: Record<string, string>;
  customThemes?: CustomTheme[];
}
