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
      <View style={styles.cardHeader}>
        <View style={styles.titleRow}>
          <Icon name="check-circle" size={20} color={colors.success} />
          <Text style={[styles.cardTitle, { color: colors.text }]}>
            {connection.mode === 'cloud' ? 'Codedock Cloud' : 'Connected Instance'}
          </Text>
        </View>
        <Pressable onPress={onDisconnect} style={styles.disconnectBtn}>
          <Text style={[styles.disconnectText, { color: colors.danger }]}>Disconnect</Text>
        </Pressable>
      </View>
      <Text style={[styles.cardDetail, { color: colors.textSecondary }]}>{connection.url}</Text>
      {connection.userEmail ? (
        <Text style={[styles.cardUser, { color: colors.textMuted }]}>
          Signed in as {connection.userEmail} {connection.version ? `(v${connection.version})` : ''}
        </Text>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  connectedCard: {
    borderWidth: 1,
    borderRadius: 12,
    padding: 16,
    marginBottom: 20,
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
    fontWeight: '700',
  },
  cardDetail: {
    fontSize: 13,
    marginBottom: 4,
  },
  cardUser: {
    fontSize: 12,
  },
  disconnectBtn: {
    paddingVertical: 4,
    paddingHorizontal: 8,
  },
  disconnectText: {
    fontSize: 12,
    fontWeight: '600',
  },
});
