'use client';

import { ThemeProvider as NextThemesProvider, useTheme } from 'next-themes';
import type * as React from 'react';
import { useEffect } from 'react';

export type Theme = 'light' | 'dim' | 'dark' | 'system';

export interface ThemeProviderProps extends React.ComponentProps<typeof NextThemesProvider> {
  children: React.ReactNode;
}

function ThemeSync() {
  const { theme, resolvedTheme } = useTheme();

  useEffect(() => {
    if (typeof document === 'undefined') return;
    const root = document.documentElement;

    let effective: 'light' | 'dim' | 'dark' = 'dim';
    if (theme === 'system') {
      const isDark =
        typeof window !== 'undefined' &&
        window.matchMedia?.('(prefers-color-scheme: dark)').matches;
      effective = isDark ? 'dim' : 'light';
    } else if (theme === 'light' || theme === 'dim' || theme === 'dark') {
      effective = theme;
    } else if (resolvedTheme === 'light' || resolvedTheme === 'dim' || resolvedTheme === 'dark') {
      effective = resolvedTheme;
    }

    root.setAttribute('data-theme', effective);
    root.classList.remove('light', 'dim', 'dark');
    root.classList.add(effective);
  }, [theme, resolvedTheme]);

  return null;
}

export function ThemeProvider({
  children,
  attribute = 'class',
  defaultTheme = 'system',
  enableSystem = true,
  disableTransitionOnChange = true,
  themes = ['light', 'dim', 'dark', 'system'],
  ...props
}: ThemeProviderProps) {
  return (
    <NextThemesProvider
      attribute={attribute}
      defaultTheme={defaultTheme}
      enableSystem={enableSystem}
      disableTransitionOnChange={disableTransitionOnChange}
      themes={themes}
      {...props}
    >
      <ThemeSync />
      {children}
    </NextThemesProvider>
  );
}

export { useTheme };
