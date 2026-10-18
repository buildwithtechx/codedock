import { useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import type { CodedockDatabase } from '@/interfaces';
import {
  useDatabasesQuery,
  useRestartDatabaseMutation,
  useStartDatabaseMutation,
  useStopDatabaseMutation,
} from '@/hooks/useDatabases';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';
import { formatStatus } from '@/utils';

export default function DatabasesScreen() {
  const { colors } = useAppTheme();
  const { connection } = useAppStore();
  const { data: databases = [], isLoading, isRefetching, refetch } = useDatabasesQuery();
  const restartMutation = useRestartDatabaseMutation();
  const startMutation = useStartDatabaseMutation();
  const stopMutation = useStopDatabaseMutation();
  const [actingId, setActingId] = useState<string | null>(null);

  const handleRestart = (db: CodedockDatabase) => {
    Alert.alert('Restart Database', `Restart ${db.name}?`, [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Restart',
        style: 'destructive',
        onPress: async () => {
          setActingId(db.id);
          try {
            await restartMutation.mutateAsync(db.id);
          } catch (e: any) {
            Alert.alert('Error', e.message || 'Failed to restart database');
          } finally {
            setActingId(null);
          }
        },
      },
    ]);
  };

  const handleToggle = (db: CodedockDatabase) => {
    const isRunning = db.status === 'running';
    const action = isRunning ? 'Stop' : 'Start';

    Alert.alert(`${action} Database`, `${action} ${db.name}?`, [
      { text: 'Cancel', style: 'cancel' },
      {
        text: action,
        onPress: async () => {
          setActingId(db.id);
          try {
            if (isRunning) {
              await stopMutation.mutateAsync(db.id);
            } else {
              await startMutation.mutateAsync(db.id);
            }
          } catch (e: any) {
            Alert.alert('Error', e.message || `Failed to ${action.toLowerCase()} database`);
          } finally {
            setActingId(null);
          }
        },
      },
    ]);
  };

  if (!connection) {
    return (
      <View style={[styles.center, { backgroundColor: colors.background }]}>
        <Text style={[styles.title, { color: colors.text }]}>Not Connected</Text>
        <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
          Configure instance settings to view databases.
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
          data={databases}
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
                No managed databases found on this instance.
              </Text>
            </View>
          }
          renderItem={({ item }) => {
            const status = formatStatus(item.status);
            const isBusy = actingId === item.id;

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
                    <View style={[styles.engineBadge, { backgroundColor: colors.surfaceSubtle }]}>
                      <Text style={[styles.engineText, { color: colors.textSecondary }]}>
                        {item.engine?.toUpperCase() || 'DB'}
                      </Text>
                    </View>
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
                  {item.internalHost || 'localhost'}:{item.port || 5432}
                </Text>

                <View style={styles.actionRow}>
                  <Pressable
                    style={[
                      styles.actionBtn,
                      { borderColor: colors.border, backgroundColor: colors.surfaceSubtle },
                    ]}
                    onPress={() => handleToggle(item)}
                    disabled={isBusy}
                  >
                    <Text style={[styles.actionBtnText, { color: colors.text }]}>
                      {item.status === 'running' ? 'Stop' : 'Start'}
                    </Text>
                  </Pressable>
                  <Pressable
                    style={[
                      styles.actionBtn,
                      { borderColor: colors.border, backgroundColor: colors.surfaceSubtle },
                    ]}
                    onPress={() => handleRestart(item)}
                    disabled={isBusy}
                  >
                    <Text style={[styles.actionBtnText, { color: colors.text }]}>Restart</Text>
                  </Pressable>
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
  engineBadge: {
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 4,
  },
  engineText: {
    fontSize: 11,
    fontWeight: '700',
  },
  cardSub: {
    fontSize: 13,
    marginBottom: 12,
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
  actionRow: {
    flexDirection: 'row',
    gap: 8,
  },
  actionBtn: {
    borderWidth: 1,
    paddingVertical: 6,
    paddingHorizontal: 14,
    borderRadius: 6,
  },
  actionBtnText: {
    fontSize: 12,
    fontWeight: '600',
  },
});
