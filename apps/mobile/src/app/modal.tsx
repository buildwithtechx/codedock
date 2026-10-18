import { useRouter } from 'expo-router';
import { useState } from 'react';
import { ActivityIndicator, Pressable, StyleSheet, Text, View } from 'react-native';
import Toast from 'react-native-toast-message';
import { Icon } from '@/components/Icon';
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
      <View style={[styles.container, { backgroundColor: colors.surface }]}>
        <View style={[styles.grabber, { backgroundColor: colors.border }]} />
        <Text style={[styles.title, { color: colors.text }]}>No App Selected</Text>
        <Pressable style={styles.closeButton} onPress={() => router.back()}>
          <Text style={[styles.closeButtonText, { color: colors.textMuted }]}>Close</Text>
        </Pressable>
      </View>
    );
  }

  const handleAction = async (action: 'redeploy' | 'restart' | 'stop') => {
    if (!connection) {
      Toast.show({
        type: 'error',
        text1: 'Authentication Error',
        text2: 'Not connected to Codedock instance',
      });
      return;
    }

    setBusyAction(action);
    try {
      if (action === 'redeploy') {
        await MobileApi.redeployApp(connection.url, connection.token, selectedApp.id);
        Toast.show({
          type: 'success',
          text1: 'Redeploy Triggered',
          text2: `Queued rebuild for ${selectedApp.name}`,
        });
      } else if (action === 'restart') {
        await MobileApi.restartApp(connection.url, connection.token, selectedApp.id);
        Toast.show({
          type: 'success',
          text1: 'Restart Triggered',
          text2: `Restarting container for ${selectedApp.name}`,
        });
      } else if (action === 'stop') {
        await MobileApi.stopApp(connection.url, connection.token, selectedApp.id);
        Toast.show({
          type: 'info',
          text1: 'Service Stopped',
          text2: `${selectedApp.name} has been stopped`,
        });
      }
      router.back();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      Toast.show({
        type: 'error',
        text1: 'Action Failed',
        text2: msg,
      });
    } finally {
      setBusyAction(null);
    }
  };

  return (
    <View style={[styles.container, { backgroundColor: colors.surface }]}>
      <View style={[styles.grabber, { backgroundColor: colors.border }]} />

      <View style={styles.header}>
        <View style={{ flex: 1 }}>
          <Text style={[styles.title, { color: colors.text }]}>{selectedApp.name}</Text>
          <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
            {selectedApp.projectName || 'Default Project'} · {selectedApp.status || 'stopped'}
          </Text>
        </View>
        <Pressable onPress={() => router.back()} style={styles.dismissBtn}>
          <Icon name="close" size={20} color={colors.textMuted} />
        </Pressable>
      </View>

      <View
        style={[
          styles.infoCard,
          { backgroundColor: colors.background, borderColor: colors.border },
        ]}
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
        {selectedApp.branch ? (
          <View style={styles.infoRow}>
            <Text style={[styles.infoLabel, { color: colors.textSecondary }]}>Branch</Text>
            <Text style={[styles.infoValue, { color: colors.text }]}>{selectedApp.branch}</Text>
          </View>
        ) : null}
      </View>

      <Text style={[styles.sectionHeader, { color: colors.text }]}>Service Operations</Text>

      <Pressable
        style={[styles.actionBtn, { backgroundColor: colors.primary }]}
        onPress={() => handleAction('redeploy')}
        disabled={!!busyAction}
      >
        {busyAction === 'redeploy' ? (
          <ActivityIndicator color="#FFFFFF" />
        ) : (
          <View style={styles.btnInner}>
            <Icon name="rocket-launch-outline" size={18} color="#FFFFFF" />
            <Text style={styles.actionBtnText}>Redeploy Service</Text>
          </View>
        )}
      </Pressable>

      <Pressable
        style={[styles.actionBtn, { backgroundColor: colors.surfaceSubtle }]}
        onPress={() => handleAction('restart')}
        disabled={!!busyAction}
      >
        {busyAction === 'restart' ? (
          <ActivityIndicator color={colors.text} />
        ) : (
          <View style={styles.btnInner}>
            <Icon name="restart" size={18} color={colors.text} />
            <Text style={[styles.actionBtnText, { color: colors.text }]}>Restart Service</Text>
          </View>
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
          <View style={styles.btnInner}>
            <Icon name="stop-circle-outline" size={18} color={colors.danger} />
            <Text style={[styles.actionBtnText, { color: colors.danger }]}>Stop Service</Text>
          </View>
        )}
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
  },
  grabber: {
    width: 40,
    height: 4,
    borderRadius: 2,
    alignSelf: 'center',
    marginBottom: 16,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 16,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  subtitle: {
    fontSize: 13,
    marginTop: 2,
  },
  dismissBtn: {
    padding: 6,
  },
  infoCard: {
    borderWidth: 1,
    borderRadius: 12,
    padding: 14,
    marginBottom: 20,
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: 5,
  },
  infoLabel: {
    fontSize: 13,
  },
  infoValue: {
    fontSize: 13,
    fontWeight: '500',
  },
  sectionHeader: {
    fontSize: 15,
    fontWeight: '600',
    marginBottom: 12,
  },
  actionBtn: {
    borderRadius: 10,
    paddingVertical: 13,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 10,
  },
  btnInner: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  actionBtnText: {
    color: '#FFFFFF',
    fontWeight: '600',
    fontSize: 14,
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
