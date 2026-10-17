'use client';

import { ThemeProvider as NextThemesProvider, useTheme } from 'next-themes';
import type * as React from 'react';

export type Theme = 'light' | 'dim' | 'dark' | 'system';

export interface ThemeProviderProps extends React.ComponentProps<typeof NextThemesProvider> {
  children: React.ReactNode;
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
      {children}
    </NextThemesProvider>
  );
}

export { useTheme };
