import { useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  Pressable,
  ScrollView,
  StyleSheet,
  Switch,
  Text,
  TextInput,
  View,
} from 'react-native';
import { MobileApi } from '@/services/api';
import { registerForPushNotificationsAsync } from '@/services/notifications';
import { MobileStorage } from '@/services/storage';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';

export default function SettingsScreen() {
  const { colors } = useAppTheme();
  const { connection, setConnection, setApps } = useAppStore();

  const [url, setUrl] = useState(connection?.url || 'http://localhost:8080');
  const [token, setToken] = useState(connection?.token || '');
  const [loading, setLoading] = useState(false);
  const [pushEnabled, setPushEnabled] = useState(true);

  const onConnect = async () => {
    if (!url || !token) {
      Alert.alert('Validation Error', 'Instance URL and PAT token are required.');
      return;
    }

    setLoading(true);
    try {
      const meta = await MobileApi.getMeta(url);
      const me = await MobileApi.getMe(url, token);

      await MobileStorage.saveConnection(url, token, me.email);
      setConnection({
        url,
        token,
        userEmail: me.email,
        userName: me.name,
        version: meta.version,
      });

      if (pushEnabled) {
        await registerForPushNotificationsAsync(url, token);
      }

      Alert.alert('Connected', `Connected to Codedock (${meta.version}) as ${me.email}`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      Alert.alert('Connection Failed', msg);
    } finally {
      setLoading(false);
    }
  };

  const onDisconnect = async () => {
    Alert.alert('Disconnect', 'Clear stored credentials from this device?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Disconnect',
        style: 'destructive',
        onPress: async () => {
          await MobileStorage.clearConnection();
          setConnection(null);
          setApps([]);
          setUrl('');
          setToken('');
        },
      },
    ]);
  };

  return (
    <ScrollView
      style={[styles.container, { backgroundColor: colors.background }]}
      contentContainerStyle={styles.content}
    >
      <Text style={[styles.headerTitle, { color: colors.text }]}>Instance Connection</Text>

      <Text style={[styles.label, { color: colors.textSecondary }]}>Codedock Instance URL</Text>
      <TextInput
        style={[
          styles.input,
          {
            backgroundColor: colors.surface,
            borderColor: colors.border,
            color: colors.text,
          },
        ]}
        placeholder="https://codedock.example.com"
        placeholderTextColor={colors.textMuted}
        value={url}
        onChangeText={setUrl}
        autoCapitalize="none"
        autoCorrect={false}
      />

      <Text style={[styles.label, { color: colors.textSecondary }]}>
        Personal Access Token (PAT)
      </Text>
      <TextInput
        style={[
          styles.input,
          {
            backgroundColor: colors.surface,
            borderColor: colors.border,
            color: colors.text,
          },
        ]}
        placeholder="pat_..."
        placeholderTextColor={colors.textMuted}
        value={token}
        onChangeText={setToken}
        secureTextEntry
        autoCapitalize="none"
        autoCorrect={false}
      />

      <View style={styles.switchRow}>
        <View>
          <Text style={[styles.switchLabel, { color: colors.text }]}>Push Notifications</Text>
          <Text style={[styles.switchSub, { color: colors.textMuted }]}>
            Receive alerts on deployment and server events
          </Text>
        </View>
        <Switch
          value={pushEnabled}
          onValueChange={setPushEnabled}
          trackColor={{ false: colors.border, true: colors.primary }}
        />
      </View>

      <Pressable
        style={[styles.button, { backgroundColor: colors.primary }]}
        onPress={onConnect}
        disabled={loading}
      >
        {loading ? (
          <ActivityIndicator color="#FFFFFF" />
        ) : (
          <Text style={styles.buttonText}>
            {connection ? 'Update Connection' : 'Connect to Instance'}
          </Text>
        )}
      </Pressable>

      {connection ? (
        <View
          style={[
            styles.connectedCard,
            { backgroundColor: colors.surface, borderColor: colors.border },
          ]}
        >
          <Text style={[styles.connectedTitle, { color: colors.success }]}>
            Connected to Codedock
          </Text>
          <Text style={[styles.connectedText, { color: colors.textSecondary }]}>
            User: {connection.userEmail || 'Active'}
          </Text>
          <Text style={[styles.connectedText, { color: colors.textSecondary }]}>
            Version: {connection.version || 'v1.0.0'}
          </Text>

          <Pressable
            style={[styles.dangerButton, { backgroundColor: colors.dangerBackground }]}
            onPress={onDisconnect}
          >
            <Text style={[styles.dangerButtonText, { color: colors.danger }]}>Disconnect</Text>
          </Pressable>
        </View>
      ) : null}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  content: {
    padding: 20,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  label: {
    fontSize: 13,
    marginBottom: 6,
  },
  input: {
    borderWidth: 1,
    borderRadius: 8,
    padding: 12,
    fontSize: 15,
    marginBottom: 16,
  },
  switchRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 24,
    paddingVertical: 8,
  },
  switchLabel: {
    fontSize: 15,
    fontWeight: '500',
  },
  switchSub: {
    fontSize: 12,
    marginTop: 2,
  },
  button: {
    borderRadius: 8,
    padding: 14,
    alignItems: 'center',
    marginBottom: 20,
  },
  buttonText: {
    color: '#FFFFFF',
    fontWeight: '600',
    fontSize: 15,
  },
  connectedCard: {
    borderWidth: 1,
    borderRadius: 8,
    padding: 16,
    marginTop: 10,
  },
  connectedTitle: {
    fontWeight: '600',
    fontSize: 15,
    marginBottom: 6,
  },
  connectedText: {
    fontSize: 13,
    marginBottom: 4,
  },
  dangerButton: {
    marginTop: 12,
    paddingVertical: 10,
    borderRadius: 6,
    alignItems: 'center',
  },
  dangerButtonText: {
    fontWeight: '600',
    fontSize: 14,
  },
});
