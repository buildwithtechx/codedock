import { useRouter } from 'expo-router';
import { useState } from 'react';
import { ActivityIndicator, Alert, Pressable, StyleSheet, Text, View } from 'react-native';
import { MobileApi } from '@/services/api';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';

export default function AppDetailsModal() {
  const router = useRouter();
  const { colors } = useAppTheme();
  const { selectedApp, connection } = useAppStore();
  const [busyAction, setBusyAction] = useState<string | null>(null);

  if (!selectedApp) {
    return (
      <View style={[styles.container, { backgroundColor: colors.background }]}>
        <Text style={[styles.title, { color: colors.text }]}>No app selected</Text>
        <Pressable style={styles.closeButton} onPress={() => router.back()}>
          <Text style={[styles.closeButtonText, { color: colors.textMuted }]}>Close</Text>
        </Pressable>
      </View>
    );
  }

  const handleAction = async (action: 'redeploy' | 'restart' | 'stop') => {
    if (!connection) {
      Alert.alert('Error', 'Not connected to Codedock instance');
      return;
    }

    setBusyAction(action);
    try {
      if (action === 'redeploy') {
        await MobileApi.redeployApp(connection.url, connection.token, selectedApp.id);
        Alert.alert('Success', `Redeploy triggered for ${selectedApp.name}`);
      } else if (action === 'restart') {
        await MobileApi.restartApp(connection.url, connection.token, selectedApp.id);
        Alert.alert('Success', `Restart triggered for ${selectedApp.name}`);
      } else if (action === 'stop') {
        await MobileApi.stopApp(connection.url, connection.token, selectedApp.id);
        Alert.alert('Success', `Service ${selectedApp.name} stopped`);
      }
      router.back();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      Alert.alert('Action Failed', msg);
    } finally {
      setBusyAction(null);
    }
  };

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <Text style={[styles.title, { color: colors.text }]}>{selectedApp.name}</Text>
      <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
        {selectedApp.projectName || 'Default Project'} · {selectedApp.status || 'stopped'}
      </Text>

      <View
        style={[styles.infoCard, { backgroundColor: colors.surface, borderColor: colors.border }]}
      >
        <View style={styles.infoRow}>
          <Text style={[styles.infoLabel, { color: colors.textSecondary }]}>Status</Text>
          <Text style={[styles.infoValue, { color: colors.text }]}>
            {selectedApp.status || 'stopped'}
          </Text>
        </View>
        <View style={styles.infoRow}>
          <Text style={[styles.infoLabel, { color: colors.textSecondary }]}>Domain</Text>
          <Text style={[styles.infoValue, { color: colors.text }]}>
            {selectedApp.domain || 'None'}
          </Text>
        </View>
        <View style={styles.infoRow}>
          <Text style={[styles.infoLabel, { color: colors.textSecondary }]}>Git Repo</Text>
          <Text style={[styles.infoValue, { color: colors.text }]}>
            {selectedApp.gitRepo || 'None'}
          </Text>
        </View>
      </View>

      <Text style={[styles.sectionHeader, { color: colors.text }]}>Service Actions</Text>

      <Pressable
        style={[styles.actionBtn, { backgroundColor: colors.primary }]}
        onPress={() => handleAction('redeploy')}
        disabled={!!busyAction}
      >
        {busyAction === 'redeploy' ? (
          <ActivityIndicator color="#FFFFFF" />
        ) : (
          <Text style={styles.actionBtnText}>Redeploy Service</Text>
        )}
      </Pressable>

      <Pressable
        style={[styles.actionBtn, { backgroundColor: colors.surfaceSubtle }]}
        onPress={() => handleAction('restart')}
        disabled={!!busyAction}
      >
        {busyAction === 'restart' ? (
          <ActivityIndicator color="#FFFFFF" />
        ) : (
          <Text style={styles.actionBtnText}>Restart Service</Text>
        )}
      </Pressable>

      <Pressable
        style={[
          styles.actionBtn,
          {
            backgroundColor: colors.dangerBackground,
            borderWidth: 1,
            borderColor: colors.danger,
          },
        ]}
        onPress={() => handleAction('stop')}
        disabled={!!busyAction}
      >
        {busyAction === 'stop' ? (
          <ActivityIndicator color={colors.danger} />
        ) : (
          <Text style={[styles.actionBtnText, { color: colors.danger }]}>Stop Service</Text>
        )}
      </Pressable>

      <Pressable style={styles.closeButton} onPress={() => router.back()}>
        <Text style={[styles.closeButtonText, { color: colors.textMuted }]}>Done</Text>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 24,
  },
  title: {
    fontSize: 22,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  subtitle: {
    fontSize: 14,
    marginBottom: 20,
  },
  infoCard: {
    borderWidth: 1,
    borderRadius: 8,
    padding: 16,
    marginBottom: 24,
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: 6,
  },
  infoLabel: {
    fontSize: 14,
  },
  infoValue: {
    fontSize: 14,
    fontWeight: '500',
  },
  sectionHeader: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 12,
  },
  actionBtn: {
    borderRadius: 8,
    paddingVertical: 14,
    alignItems: 'center',
    marginBottom: 12,
  },
  actionBtnText: {
    color: '#FFFFFF',
    fontWeight: '600',
    fontSize: 15,
  },
  closeButton: {
    marginTop: 'auto',
    paddingVertical: 14,
    alignItems: 'center',
  },
  closeButtonText: {
    fontWeight: '500',
    fontSize: 15,
  },
});
