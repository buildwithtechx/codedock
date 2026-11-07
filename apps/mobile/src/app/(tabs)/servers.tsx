import { ActivityIndicator, FlatList, RefreshControl, StyleSheet, Text, View } from 'react-native';
import { useServersQuery } from '@/hooks/useServers';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';
import { formatStatus } from '@/utils';

export default function ServersScreen() {
  const { colors } = useAppTheme();
  const { connection } = useAppStore();
  const { data: servers = [], isLoading, isRefetching, refetch } = useServersQuery();

  if (!connection) {
    return (
      <View style={[styles.center, { backgroundColor: colors.background }]}>
        <Text style={[styles.title, { color: colors.text }]}>Not Connected</Text>
        <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
          Configure instance settings to view connected servers.
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
          data={servers}
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
                No servers connected to this instance.
              </Text>
            </View>
          }
          renderItem={({ item }) => {
            const status = formatStatus(item.status);

            return (
              <View
                style={[
                  styles.card,
                  { backgroundColor: colors.surface, borderColor: colors.border },
                ]}
              >
                <View style={styles.cardHeader}>
                  <View style={styles.titleRow}>
                    <Text style={[styles.cardTitle, { color: colors.text }]}>{item.name}</Text>
                    {item.isControlPlane && (
                      <View style={[styles.roleBadge, { backgroundColor: colors.surfaceSubtle }]}>
                        <Text style={[styles.roleText, { color: colors.primary }]}>CONTROL</Text>
                      </View>
                    )}
                  </View>
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
                  {item.ipAddress || item.sshHost || '127.0.0.1'}
                  {item.sshPort ? ` :${item.sshPort}` : ''}
                </Text>

                <View style={styles.metaRow}>
                  <Text style={[styles.metaText, { color: colors.textSecondary }]}>
                    {item.provider || 'Self-hosted'}
                  </Text>
                  {item.region ? (
                    <Text style={[styles.metaText, { color: colors.textMuted }]}>
                      {item.region}
                    </Text>
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
  titleRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  cardTitle: {
    fontSize: 16,
    fontWeight: '600',
  },
  roleBadge: {
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 4,
  },
  roleText: {
    fontSize: 10,
    fontWeight: '700',
  },
  cardSub: {
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
