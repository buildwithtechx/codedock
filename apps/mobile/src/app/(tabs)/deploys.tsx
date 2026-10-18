import { ActivityIndicator, FlatList, RefreshControl, StyleSheet, Text, View } from 'react-native';
import { useDeploymentsQuery } from '@/hooks/useDeployments';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';
import { formatDuration, formatStatus } from '@/utils';

export default function DeploymentsScreen() {
  const { colors } = useAppTheme();
  const { connection } = useAppStore();
  const { data: deployments = [], isLoading, isRefetching, refetch } = useDeploymentsQuery();

  if (!connection) {
    return (
      <View style={[styles.center, { backgroundColor: colors.background }]}>
        <Text style={[styles.title, { color: colors.text }]}>Not Connected</Text>
        <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
          Configure instance settings to view deployments.
        </Text>
      </View>
    );
  }

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      {isLoading ? (
        <ActivityIndicator style={{ marginTop: 24 }} color={colors.primary} />
      ) : (
        <FlatList
          data={deployments}
          keyExtractor={(item) => item.id}
          refreshControl={
            <RefreshControl
              refreshing={isRefetching}
              onRefresh={refetch}
              tintColor={colors.primary}
            />
          }
          ListEmptyComponent={
            <View style={styles.center}>
              <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
                No deployment history found.
              </Text>
            </View>
          }
          renderItem={({ item }) => {
            const status = formatStatus(item.status);
            const duration = formatDuration(item.durationSeconds);

            return (
              <View
                style={[
                  styles.card,
                  { backgroundColor: colors.surface, borderColor: colors.border },
                ]}
              >
                <View style={styles.cardHeader}>
                  <Text style={[styles.cardTitle, { color: colors.text }]}>
                    {item.serviceName || `Deploy ${item.id.slice(0, 8)}`}
                  </Text>
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
                <Text style={[styles.cardCommit, { color: colors.textSecondary }]}>
                  {item.commitMessage || item.commitHash || 'Manual deployment'}
                </Text>
                <View style={styles.metaRow}>
                  <Text style={[styles.metaText, { color: colors.textMuted }]}>
                    {item.branch || 'main'}
                  </Text>
                  {duration ? (
                    <Text style={[styles.metaText, { color: colors.textMuted }]}>{duration}</Text>
                  ) : null}
                </View>
              </View>
            );
          }}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 14,
    textAlign: 'center',
  },
  card: {
    borderWidth: 1,
    borderRadius: 10,
    padding: 16,
    marginBottom: 12,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 6,
  },
  cardTitle: {
    fontSize: 15,
    fontWeight: '600',
  },
  cardCommit: {
    fontSize: 13,
    marginBottom: 8,
  },
  metaRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  metaText: {
    fontSize: 12,
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
