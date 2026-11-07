import { StyleSheet, Text, View } from 'react-native';
import { Icon } from '@/components/Icon';
import { useAppTheme } from '@/theme';

export function CloudNoticeCard({ cloudUrl }: { cloudUrl: string }) {
  const { colors } = useAppTheme();

  return (
    <View
      style={[styles.cloudNotice, { backgroundColor: colors.surface, borderColor: colors.border }]}
    >
      <View style={styles.cloudNoticeHeader}>
        <Icon name="check-circle-outline" size={17} color={colors.success} />
        <Text style={[styles.cloudNoticeTitle, { color: colors.text }]}>
          Managed Cloud Platform
        </Text>
      </View>
      <Text style={[styles.cloudNoticeText, { color: colors.textSecondary }]}>
        Connected to {cloudUrl}. Authorize seamlessly via your browser without entering a token, or
        provide a Personal Access Token below.
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  cloudNotice: {
    borderWidth: 1,
    borderRadius: 12,
    padding: 14,
    marginBottom: 18,
  },
  cloudNoticeHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    marginBottom: 4,
  },
  cloudNoticeTitle: {
    fontSize: 13,
    fontWeight: '600',
  },
  cloudNoticeText: {
    fontSize: 12,
    lineHeight: 17,
  },
});
