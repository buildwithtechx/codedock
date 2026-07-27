import { createContext, type ReactNode, useContext, useEffect, useState } from "react";

type Theme = "dark";

type ThemeContextValue = {
  theme: Theme;
  toggleTheme: () => void;
  setTheme: (theme: Theme) => void;
};

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);
const FIXED_THEME: Theme = "dark";

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(FIXED_THEME);

  useEffect(() => {
    const root = document.documentElement;
    root.classList.remove("light", "dark");
    root.classList.add(FIXED_THEME);
  }, []);

  const value: ThemeContextValue = {
    theme,
    setTheme: () => setThemeState(FIXED_THEME),
    toggleTheme: () => setThemeState(FIXED_THEME),
  };

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useTheme must be used within a ThemeProvider");
  return ctx;
}
