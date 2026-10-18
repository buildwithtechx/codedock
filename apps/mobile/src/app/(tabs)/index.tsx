import { useRouter } from 'expo-router';
import {
  ActivityIndicator,
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import type { CodedockApp } from '@/interfaces';
import { useAppsQuery } from '@/hooks/useWorkloads';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';
import { formatStatus } from '@/utils';

export default function WorkloadsScreen() {
  const router = useRouter();
  const { colors } = useAppTheme();
  const { connection, setSelectedApp } = useAppStore();
  const { data: apps = [], isLoading, isRefetching, refetch } = useAppsQuery();

  const onSelectApp = (app: CodedockApp) => {
    setSelectedApp(app);
    router.push('/modal');
  };

  if (!connection) {
    return (
      <View style={[styles.center, { backgroundColor: colors.background }]}>
        <Text style={[styles.title, { color: colors.text }]}>Connect to Codedock</Text>
        <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
          Configure your instance URL and token to view workloads.
        </Text>
        <Pressable
          style={[styles.button, { backgroundColor: colors.primary }]}
          onPress={() => router.push('/settings')}
        >
          <Text style={styles.buttonText}>Open Settings</Text>
        </Pressable>
      </View>
    );
  }

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      {isLoading ? (
        <ActivityIndicator style={{ marginTop: 24 }} color={colors.primary} />
      ) : (
        <FlatList
          data={apps}
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
                No apps deployed on this instance.
              </Text>
            </View>
          }
          renderItem={({ item }) => {
            const status = formatStatus(item.status);

            return (
              <Pressable
                style={[
                  styles.card,
                  { backgroundColor: colors.surface, borderColor: colors.border },
                ]}
                onPress={() => onSelectApp(item)}
              >
                <View style={styles.cardHeader}>
                  <Text style={[styles.cardTitle, { color: colors.text }]}>{item.name}</Text>
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
                  {item.domain || item.projectName || 'Default Project'}
                </Text>
              </Pressable>
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
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 20,
  },
  button: {
    paddingVertical: 12,
    paddingHorizontal: 24,
    borderRadius: 8,
  },
  buttonText: {
    color: '#FFFFFF',
    fontWeight: '600',
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
