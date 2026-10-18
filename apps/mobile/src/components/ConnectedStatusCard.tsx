import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Icon } from '@/components/Icon';
import type { CodedockConnection } from '@/interfaces';
import { useAppTheme } from '@/theme';

export function ConnectedStatusCard({
  connection,
  onDisconnect,
}: {
  connection: CodedockConnection;
  onDisconnect: () => void;
}) {
  const { colors } = useAppTheme();

  return (
    <View
      style={[
        styles.connectedCard,
        { backgroundColor: colors.surface, borderColor: colors.border },
      ]}
    >
      <View style={styles.header}>
        <Icon name="check-circle" size={18} color={colors.success} />
        <Text style={[styles.connectedTitle, { color: colors.success }]}>
          {connection.mode === 'cloud' ? 'Codedock Cloud Connected' : 'Codedock Instance Active'}
        </Text>
      </View>
      <Text style={[styles.connectedText, { color: colors.textSecondary }]}>
        User: {connection.userEmail || 'Active'}
      </Text>
      <Text style={[styles.connectedText, { color: colors.textSecondary }]}>
        Endpoint: {connection.url}
      </Text>
      <Text style={[styles.connectedText, { color: colors.textSecondary }]}>
        Version: {connection.version || 'v1.0.0'}
      </Text>

      <Pressable
        style={[styles.dangerButton, { backgroundColor: colors.dangerBackground }]}
        onPress={onDisconnect}
      >
        <Text style={[styles.dangerButtonText, { color: colors.danger }]}>
          Sign Out / Disconnect
        </Text>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  connectedCard: {
    borderWidth: 1,
    borderRadius: 10,
    padding: 16,
    marginTop: 8,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    marginBottom: 8,
  },
  connectedTitle: {
    fontWeight: '600',
    fontSize: 15,
  },
  connectedText: {
    fontSize: 13,
    marginBottom: 3,
  },
  dangerButton: {
    marginTop: 14,
    paddingVertical: 10,
    borderRadius: 8,
    alignItems: 'center',
  },
  dangerButtonText: {
    fontWeight: '600',
    fontSize: 14,
  },
});
