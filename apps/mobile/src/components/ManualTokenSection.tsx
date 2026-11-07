import { ActivityIndicator, Pressable, StyleSheet, Text, TextInput, View } from 'react-native';
import { Icon } from '@/components/Icon';
import { useAppTheme } from '@/theme';

export function ManualTokenSection({
  mode,
  url,
  token,
  loading,
  onChangeUrl,
  onChangeToken,
  onConnect,
}: {
  mode: 'cloud' | 'self-hosted';
  url: string;
  token: string;
  loading: boolean;
  onChangeUrl: (u: string) => void;
  onChangeToken: (t: string) => void;
  onConnect: () => void;
}) {
  const { colors } = useAppTheme();

  return (
    <View style={{ marginTop: mode === 'cloud' ? 8 : 0 }}>
      {mode === 'self-hosted' && (
        <View style={styles.fieldGroup}>
          <Text style={[styles.label, { color: colors.textSecondary }]}>Instance URL</Text>
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
            onChangeText={onChangeUrl}
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
            {
              backgroundColor: colors.surface,
              borderColor: colors.border,
              color: colors.text,
            },
          ]}
          placeholder="pat_..."
          placeholderTextColor={colors.textMuted}
          value={token}
          onChangeText={onChangeToken}
          secureTextEntry
          autoCapitalize="none"
          autoCorrect={false}
        />
      </View>

      <Pressable
        style={[
          styles.primaryButton,
          { backgroundColor: colors.primary },
          loading && { opacity: 0.7 },
        ]}
        onPress={onConnect}
        disabled={loading}
      >
        {loading ? (
          <ActivityIndicator color="#FFFFFF" />
        ) : (
          <View style={styles.buttonInner}>
            <Text style={styles.primaryButtonText}>
              {mode === 'cloud' ? 'Connect with Token' : 'Connect to Instance'}
            </Text>
            <Icon name="arrow-right" size={18} color="#FFFFFF" />
          </View>
        )}
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  fieldGroup: {
    marginBottom: 16,
  },
  label: {
    fontSize: 13,
    fontWeight: '500',
    marginBottom: 6,
  },
  input: {
    borderWidth: 1,
    borderRadius: 10,
    paddingHorizontal: 14,
    paddingVertical: 11,
    fontSize: 14,
  },
  primaryButton: {
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 6,
  },
  buttonInner: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  primaryButtonText: {
    color: '#FFFFFF',
    fontWeight: '700',
    fontSize: 15,
  },
});
