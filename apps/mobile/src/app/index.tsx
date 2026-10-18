import * as Linking from 'expo-linking';
import { Redirect, useRouter } from 'expo-router';
import * as WebBrowser from 'expo-web-browser';
import { useEffect, useRef, useState } from 'react';
import {
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Switch,
  Text,
  View,
} from 'react-native';
import Toast from 'react-native-toast-message';
import { Icon } from '@/components/Icon';
import { AuthBrandHeader } from '@/features/auth/AuthBrandHeader';
import { AuthModeSelector } from '@/features/auth/AuthModeSelector';
import { CloudNoticeCard } from '@/features/auth/CloudNoticeCard';
import { ManualTokenSection } from '@/features/auth/ManualTokenSection';
import { DEFAULT_CLOUD_URL, MobileApi } from '@/services/api';
import { registerForPushNotificationsAsync } from '@/services/notifications';
import { MobileStorage } from '@/services/storage';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';

WebBrowser.maybeCompleteAuthSession();

export default function IndexAuthScreen() {
  const router = useRouter();
  const { colors } = useAppTheme();
  const connection = useAppStore((state) => state.connection);
  const setConnection = useAppStore((state) => state.setConnection);

  const [checkingAuth, setCheckingAuth] = useState(true);
  const [mode, setMode] = useState<'cloud' | 'self-hosted'>('cloud');
  const [url, setUrl] = useState(DEFAULT_CLOUD_URL);
  const [token, setToken] = useState('');
  const [showManualToken, setShowManualToken] = useState(false);
  const [pushEnabled, setPushEnabled] = useState(true);
  const [loading, setLoading] = useState(false);
  const [authorizingDevice, setAuthorizingDevice] = useState(false);
  const pollTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    MobileStorage.getConnection().then((conn) => {
      if (conn?.url && conn?.token) setConnection(conn);
      setCheckingAuth(false);
    });

    const sub = Linking.addEventListener('url', (event) => {
      void handleDeepLink(event.url);
    });

    return () => {
      sub.remove();
      if (pollTimerRef.current) clearInterval(pollTimerRef.current);
    };
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

  const handleDeepLink = async (incomingUrl: string) => {
    const parsed = Linking.parse(incomingUrl);
    const incomingToken = parsed.queryParams?.token as string | undefined;
    const incomingEmail = parsed.queryParams?.email as string | undefined;
    if (incomingToken) {
      if (pollTimerRef.current) clearInterval(pollTimerRef.current);
      setAuthorizingDevice(false);
      await finishLogin(DEFAULT_CLOUD_URL, incomingToken, incomingEmail, 'cloud');
    }
  };

  const finishLogin = async (
    targetUrl: string,
    authToken: string,
    email?: string,
    authMode: 'cloud' | 'self-hosted' = 'cloud'
  ) => {
    setLoading(true);
    try {
      const meta = await MobileApi.getMeta(targetUrl);
      const me = await MobileApi.getMe(targetUrl, authToken);
      const userEmail = email || me.email;

      await MobileStorage.saveConnection(targetUrl, authToken, userEmail, authMode);
      setConnection({
        url: targetUrl,
        token: authToken,
        mode: authMode,
        userEmail,
        userName: me.name,
        version: meta.version,
      });

      if (pushEnabled) await registerForPushNotificationsAsync(targetUrl, authToken);

      Toast.show({
        type: 'success',
        text1: `Connected to ${authMode === 'cloud' ? 'Codedock Cloud' : 'Instance'}`,
        text2: `Signed in as ${userEmail}`,
      });
      router.replace('/(tabs)');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Connection failed';
      Toast.show({ type: 'error', text1: 'Authentication Failed', text2: msg });
    } finally {
      setLoading(false);
    }
  };

  const handleStartCloudAuth = async () => {
    setAuthorizingDevice(true);
    try {
      const targetUrl = DEFAULT_CLOUD_URL;
      const session = await MobileApi.startDeviceAuth(targetUrl, 'Codedock Mobile');
      const authUrl = `${targetUrl}${session.verificationUriComplete}`;

      void WebBrowser.openAuthSessionAsync(authUrl, 'codedock://auth/callback');

      let attempts = 0;
      if (pollTimerRef.current) clearInterval(pollTimerRef.current);
      pollTimerRef.current = setInterval(async () => {
        attempts += 1;
        if (attempts > 150) {
          if (pollTimerRef.current) clearInterval(pollTimerRef.current);
          setAuthorizingDevice(false);
          return;
        }

        try {
          const res = await MobileApi.pollDeviceToken(targetUrl, session.deviceCode);
          if (res.status === 'approved' && res.token) {
            if (pollTimerRef.current) clearInterval(pollTimerRef.current);
            setAuthorizingDevice(false);
            WebBrowser.dismissAuthSession();
            await finishLogin(targetUrl, res.token, res.email, 'cloud');
          } else if (res.status === 'denied') {
            if (pollTimerRef.current) clearInterval(pollTimerRef.current);
            setAuthorizingDevice(false);
            WebBrowser.dismissAuthSession();
            Toast.show({ type: 'error', text1: 'Authorization Declined' });
          }
        } catch {
          // Continue polling until timeout or success
        }
      }, 2000);
    } catch (err: unknown) {
      setAuthorizingDevice(false);
      const msg = err instanceof Error ? err.message : 'Failed to launch cloud login';
      Toast.show({ type: 'error', text1: 'Login Error', text2: msg });
    }
  };

  const handleManualConnect = async () => {
    const targetUrl = mode === 'cloud' ? url || DEFAULT_CLOUD_URL : url.trim();
    const cleanToken = token.trim();
    if (!targetUrl || !cleanToken) {
      Toast.show({
        type: 'error',
        text1: 'Validation Error',
        text2: !targetUrl ? 'Instance URL is required' : 'Token is required',
      });
      return;
    }
    await finishLogin(targetUrl, cleanToken, undefined, mode);
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      style={[styles.container, { backgroundColor: colors.background }]}
    >
      <ScrollView contentContainerStyle={styles.scrollContent} keyboardShouldPersistTaps="handled">
        <AuthBrandHeader />

        <AuthModeSelector
          mode={mode}
          onModeChange={(m) => {
            setMode(m);
            setUrl(m === 'cloud' ? DEFAULT_CLOUD_URL : '');
          }}
        />

        {mode === 'cloud' ? (
          <>
            <CloudNoticeCard cloudUrl={DEFAULT_CLOUD_URL} />
            <Pressable
              style={[
                styles.primaryButton,
                { backgroundColor: colors.primary },
                authorizingDevice && { opacity: 0.7 },
              ]}
              onPress={handleStartCloudAuth}
              disabled={authorizingDevice || loading}
            >
              {authorizingDevice ? (
                <ActivityIndicator color="#FFFFFF" />
              ) : (
                <View style={styles.buttonInner}>
                  <Icon name="shield-check-outline" size={19} color="#FFFFFF" />
                  <Text style={styles.primaryButtonText}>Authorize with Codedock Cloud</Text>
                </View>
              )}
            </Pressable>

            <Pressable
              onPress={() => setShowManualToken((v) => !v)}
              style={styles.manualTokenToggle}
            >
              <Text style={[styles.manualTokenToggleText, { color: colors.primary }]}>
                {showManualToken ? 'Hide manual token entry' : 'Or connect with Access Token (PAT)'}
              </Text>
            </Pressable>
          </>
        ) : null}

        {(mode === 'self-hosted' || showManualToken) && (
          <ManualTokenSection
            mode={mode}
            url={url}
            token={token}
            loading={loading}
            onChangeUrl={setUrl}
            onChangeToken={setToken}
            onConnect={handleManualConnect}
          />
        )}

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
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1 },
  center: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  scrollContent: { padding: 24, paddingTop: 56 },
  switchRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 20,
  },
  switchLabel: { fontSize: 14, fontWeight: '500' },
  switchSub: { fontSize: 11, marginTop: 2 },
  primaryButton: {
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 6,
  },
  buttonInner: { flexDirection: 'row', alignItems: 'center', gap: 8 },
  primaryButtonText: { color: '#FFFFFF', fontWeight: '700', fontSize: 15 },
  manualTokenToggle: { alignItems: 'center', marginTop: 14, paddingVertical: 4 },
  manualTokenToggleText: { fontSize: 12, fontWeight: '600' },
});
