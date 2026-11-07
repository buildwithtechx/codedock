import { StyleSheet, Text, View } from 'react-native';
import { Icon } from '@/components/Icon';
import { useAppTheme } from '@/theme';

export function AuthBrandHeader() {
  const { colors } = useAppTheme();

  return (
    <View style={styles.brandHeader}>
      <View style={[styles.iconContainer, { backgroundColor: colors.primarySoft }]}>
        <Icon name="rocket-launch-outline" size={36} color={colors.primary} />
      </View>
      <Text style={[styles.brandTitle, { color: colors.text }]}>Codedock</Text>
      <Text style={[styles.brandSubtitle, { color: colors.textSecondary }]}>
        Cloud & Self-Hosted Infrastructure
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  brandHeader: {
    alignItems: 'center',
    marginBottom: 28,
  },
  iconContainer: {
    width: 64,
    height: 64,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 14,
  },
  brandTitle: {
    fontSize: 26,
    fontWeight: '800',
    letterSpacing: -0.5,
  },
  brandSubtitle: {
    fontSize: 13,
    marginTop: 4,
  },
});
