import { useColorScheme } from 'react-native';
import { colors, type ThemeColors } from './colors';
import { spacing } from './spacing';

export function useAppTheme(): { colors: ThemeColors; spacing: typeof spacing; isDark: boolean } {
  const scheme = useColorScheme();
  const isDark = scheme !== 'light';
  return {
    colors: isDark ? colors.dark : colors.light,
    spacing,
    isDark,
  };
}

export { colors, spacing };
