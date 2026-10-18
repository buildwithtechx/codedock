import { useRouter } from 'expo-router';
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
import Toast from 'react-native-toast-message';
import { AuthModeSelector } from '@/features/auth/AuthModeSelector';
import { ConnectedStatusCard } from '@/features/auth/ConnectedStatusCard';
import { DEFAULT_CLOUD_URL, MobileApi } from '@/services/api';
import { registerForPushNotificationsAsync } from '@/services/notifications';
import { MobileStorage } from '@/services/storage';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';

export default function SettingsScreen() {
  const router = useRouter();
  const { colors } = useAppTheme();
  const { connection, setConnection, setApps, setProjects } = useAppStore();

  const [mode, setMode] = useState<'cloud' | 'self-hosted'>(
    connection?.mode || (connection?.url?.includes('codedock') ? 'cloud' : 'self-hosted')
  );
  const [url, setUrl] = useState(connection?.url || DEFAULT_CLOUD_URL);
  const [token, setToken] = useState(connection?.token || '');
  const [loading, setLoading] = useState(false);
  const [pushEnabled, setPushEnabled] = useState(true);

  const handleModeChange = (newMode: 'cloud' | 'self-hosted') => {
    setMode(newMode);
    setUrl(newMode === 'cloud' ? DEFAULT_CLOUD_URL : url === DEFAULT_CLOUD_URL ? '' : url);
  };

  const onConnect = async () => {
    const targetUrl = mode === 'cloud' ? url || DEFAULT_CLOUD_URL : url.trim();
    const cleanToken = token.trim();

    if (!targetUrl || !cleanToken) {
      Toast.show({
        type: 'error',
        text1: 'Validation Error',
        text2: 'Instance URL and PAT token are required.',
      });
      return;
    }

    setLoading(true);
    try {
      const meta = await MobileApi.getMeta(targetUrl);
      const me = await MobileApi.getMe(targetUrl, cleanToken);

      await MobileStorage.saveConnection(targetUrl, cleanToken, me.email, mode);
      setConnection({
        url: targetUrl,
        token: cleanToken,
        mode,
        userEmail: me.email,
        userName: me.name,
        version: meta.version,
      });

      if (pushEnabled) {
        await registerForPushNotificationsAsync(targetUrl, cleanToken);
      }

      Toast.show({
        type: 'success',
        text1: 'Connection Saved',
        text2: `Connected to Codedock (${meta.version}) as ${me.email}`,
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      Toast.show({ type: 'error', text1: 'Connection Failed', text2: msg });
    } finally {
      setLoading(false);
    }
  };

  const onDisconnect = () => {
    Alert.alert('Sign Out / Disconnect', 'Clear stored credentials from this device?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Sign Out',
        style: 'destructive',
        onPress: async () => {
          await MobileStorage.clearConnection();
          setConnection(null);
          setApps([]);
          setProjects([]);
          Toast.show({ type: 'info', text1: 'Signed Out', text2: 'Credentials cleared.' });
          router.replace('/');
        },
      },
    ]);
  };

  return (
    <ScrollView
      style={[styles.container, { backgroundColor: colors.background }]}
      contentContainerStyle={styles.content}
      keyboardShouldPersistTaps="handled"
    >
      <Text style={[styles.headerTitle, { color: colors.text }]}>Platform Connection</Text>

      <AuthModeSelector mode={mode} onModeChange={handleModeChange} />

      {mode === 'cloud' ? (
        <View
          style={[
            styles.cloudBanner,
            { backgroundColor: colors.surface, borderColor: colors.border },
          ]}
        >
          <Text style={[styles.cloudBannerText, { color: colors.textSecondary }]}>
            Connected to Codedock Cloud service ({DEFAULT_CLOUD_URL}).
          </Text>
        </View>
      ) : (
        <>
          <Text style={[styles.label, { color: colors.textSecondary }]}>Instance URL</Text>
          <TextInput
            style={[
              styles.input,
              { backgroundColor: colors.surface, borderColor: colors.border, color: colors.text },
            ]}
            placeholder="https://codedock.example.com"
            placeholderTextColor={colors.textMuted}
            value={url}
            onChangeText={setUrl}
            autoCapitalize="none"
            autoCorrect={false}
          />
        </>
      )}

      <Text style={[styles.label, { color: colors.textSecondary }]}>
        Personal Access Token (PAT)
      </Text>
      <TextInput
        style={[
          styles.input,
          { backgroundColor: colors.surface, borderColor: colors.border, color: colors.text },
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
        <View style={{ flex: 1, marginRight: 12 }}>
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
        <ConnectedStatusCard connection={connection} onDisconnect={onDisconnect} />
      ) : null}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1 },
  content: { padding: 20 },
  headerTitle: { fontSize: 18, fontWeight: 'bold', marginBottom: 14 },
  cloudBanner: { borderWidth: 1, borderRadius: 8, padding: 12, marginBottom: 16 },
  cloudBannerText: { fontSize: 13, lineHeight: 18 },
  label: { fontSize: 13, marginBottom: 6 },
  input: { borderWidth: 1, borderRadius: 8, padding: 12, fontSize: 14, marginBottom: 14 },
  switchRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 20,
  },
  switchLabel: { fontSize: 14, fontWeight: '500' },
  switchSub: { fontSize: 11, marginTop: 2 },
  button: { borderRadius: 10, padding: 14, alignItems: 'center', marginBottom: 16 },
  buttonText: { color: '#FFFFFF', fontWeight: '600', fontSize: 15 },
});
