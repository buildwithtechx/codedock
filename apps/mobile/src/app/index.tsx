import { Redirect, useRouter } from 'expo-router';
import { useEffect, useState } from 'react';
import {
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Switch,
  Text,
  TextInput,
  View,
} from 'react-native';
import Toast from 'react-native-toast-message';
import { AuthModeSelector } from '@/components/AuthModeSelector';
import { CloudNoticeCard } from '@/components/CloudNoticeCard';
import { Icon } from '@/components/Icon';
import { MobileApi } from '@/services/api';
import { registerForPushNotificationsAsync } from '@/services/notifications';
import { MobileStorage } from '@/services/storage';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';

const CLOUD_DEFAULT_URL = 'https://api.codedock.io';

export default function IndexAuthScreen() {
  const router = useRouter();
  const { colors } = useAppTheme();
  const connection = useAppStore((state) => state.connection);
  const setConnection = useAppStore((state) => state.setConnection);

  const [checkingAuth, setCheckingAuth] = useState(true);
  const [mode, setMode] = useState<'cloud' | 'self-hosted'>('cloud');
  const [url, setUrl] = useState(CLOUD_DEFAULT_URL);
  const [token, setToken] = useState('');
  const [pushEnabled, setPushEnabled] = useState(true);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    MobileStorage.getConnection().then((conn) => {
      if (conn?.url && conn?.token) setConnection(conn);
      setCheckingAuth(false);
    });
  }, [setConnection]);

  if (checkingAuth) {
    return (
      <View style={[styles.center, { backgroundColor: colors.background }]}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  if (connection?.url && connection?.token) {
    return <Redirect href="/(tabs)" />;
  }

  const handleModeChange = (newMode: 'cloud' | 'self-hosted') => {
    setMode(newMode);
    setUrl(newMode === 'cloud' ? CLOUD_DEFAULT_URL : '');
  };

  const handleConnect = async () => {
    const targetUrl = mode === 'cloud' ? url || CLOUD_DEFAULT_URL : url.trim();
    const cleanToken = token.trim();

    if (!targetUrl || !cleanToken) {
      Toast.show({
        type: 'error',
        text1: 'Validation Error',
        text2: !targetUrl ? 'Instance URL is required' : 'Personal Access Token is required',
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

      if (pushEnabled) await registerForPushNotificationsAsync(targetUrl, cleanToken);

      Toast.show({
        type: 'success',
        text1: `Connected to ${mode === 'cloud' ? 'Codedock Cloud' : 'Instance'}`,
        text2: `Signed in as ${me.email}`,
      });
      router.replace('/(tabs)');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Connection failed';
      Toast.show({ type: 'error', text1: 'Authentication Failed', text2: msg });
    } finally {
      setLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      style={[styles.container, { backgroundColor: colors.background }]}
    >
      <ScrollView contentContainerStyle={styles.scrollContent} keyboardShouldPersistTaps="handled">
        <View style={styles.brandHeader}>
          <View style={[styles.iconContainer, { backgroundColor: colors.primarySoft }]}>
            <Icon name="rocket-launch-outline" size={36} color={colors.primary} />
          </View>
          <Text style={[styles.brandTitle, { color: colors.text }]}>Codedock</Text>
          <Text style={[styles.brandSubtitle, { color: colors.textSecondary }]}>
            Cloud & Self-Hosted Infrastructure
          </Text>
        </View>

        <AuthModeSelector mode={mode} onModeChange={handleModeChange} />

        {mode === 'cloud' ? (
          <CloudNoticeCard cloudUrl={CLOUD_DEFAULT_URL} />
        ) : (
          <View style={styles.fieldGroup}>
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
          </View>
        )}

        <View style={styles.fieldGroup}>
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
          <Text style={[styles.helperText, { color: colors.textMuted }]}>
            Generate a token from your Codedock Dashboard under Settings &gt; Personal Tokens.
          </Text>
        </View>

        <View style={styles.switchRow}>
          <View style={{ flex: 1, marginRight: 12 }}>
            <Text style={[styles.switchLabel, { color: colors.text }]}>Push Notifications</Text>
            <Text style={[styles.switchSub, { color: colors.textMuted }]}>
              Alerts on deployment and server events
            </Text>
          </View>
          <Switch
            value={pushEnabled}
            onValueChange={setPushEnabled}
            trackColor={{ false: colors.border, true: colors.primary }}
          />
        </View>

        <Pressable
          style={[
            styles.primaryButton,
            { backgroundColor: colors.primary },
            loading && { opacity: 0.7 },
          ]}
          onPress={handleConnect}
          disabled={loading}
        >
          {loading ? (
            <ActivityIndicator color="#FFFFFF" />
          ) : (
            <View style={styles.buttonInner}>
              <Text style={styles.primaryButtonText}>
                {mode === 'cloud' ? 'Sign In to Cloud' : 'Connect to Instance'}
              </Text>
              <Icon name="arrow-right" size={18} color="#FFFFFF" />
            </View>
          )}
        </Pressable>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1 },
  center: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  scrollContent: { padding: 24, paddingTop: 56 },
  brandHeader: { alignItems: 'center', marginBottom: 28 },
  iconContainer: {
    width: 64,
    height: 64,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 14,
  },
  brandTitle: { fontSize: 26, fontWeight: '800', letterSpacing: -0.5 },
  brandSubtitle: { fontSize: 13, marginTop: 4 },
  fieldGroup: { marginBottom: 16 },
  label: { fontSize: 13, fontWeight: '500', marginBottom: 6 },
  input: {
    borderWidth: 1,
    borderRadius: 10,
    paddingHorizontal: 14,
    paddingVertical: 11,
    fontSize: 14,
  },
  helperText: { fontSize: 11, marginTop: 5, lineHeight: 15 },
  switchRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginVertical: 10,
  },
  switchLabel: { fontSize: 14, fontWeight: '500' },
  switchSub: { fontSize: 11, marginTop: 2 },
  primaryButton: {
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 16,
  },
  buttonInner: { flexDirection: 'row', alignItems: 'center', gap: 8 },
  primaryButtonText: { color: '#FFFFFF', fontWeight: '700', fontSize: 15 },
});
