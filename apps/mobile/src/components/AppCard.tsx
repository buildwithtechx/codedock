import { Pressable, StyleSheet, Text, View } from 'react-native';
import type { CodedockApp } from '@/interfaces';
import { useAppTheme } from '@/theme';
import { formatStatus } from '@/utils';

export function AppCard({ app, onPress }: { app: CodedockApp; onPress: () => void }) {
  const { colors } = useAppTheme();
  const status = formatStatus(app.status);

  return (
    <Pressable
      style={[styles.card, { backgroundColor: colors.surface, borderColor: colors.border }]}
      onPress={onPress}
    >
      <View style={styles.cardHeader}>
        <Text style={[styles.cardTitle, { color: colors.text }]}>{app.name}</Text>
        <View
          style={[
            styles.badge,
            status.isSuccess && { backgroundColor: colors.successBackground },
            status.isWarning && { backgroundColor: colors.warningBackground },
            status.isDanger && { backgroundColor: colors.dangerBackground },
          ]}
        >
          <Text
            style={[
              styles.badgeText,
              status.isSuccess && { color: colors.success },
              status.isWarning && { color: colors.warning },
              status.isDanger && { color: colors.danger },
            ]}
          >
            {status.label}
          </Text>
        </View>
      </View>
      <Text style={[styles.cardSub, { color: colors.textMuted }]}>
        {app.domain || app.projectName || 'Default Project'}
      </Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    borderWidth: 1,
    borderRadius: 10,
    padding: 16,
    marginBottom: 10,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 6,
  },
  cardTitle: {
    fontSize: 16,
    fontWeight: '600',
  },
  cardSub: {
    fontSize: 13,
  },
  badge: {
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 12,
  },
  badgeText: {
    fontSize: 11,
    fontWeight: '600',
    textTransform: 'uppercase',
  },
});
