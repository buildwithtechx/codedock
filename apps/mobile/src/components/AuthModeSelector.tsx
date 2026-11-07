import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Icon } from '@/components/Icon';
import { useAppTheme } from '@/theme';

export function AuthModeSelector({
  mode,
  onModeChange,
}: {
  mode: 'cloud' | 'self-hosted';
  onModeChange: (m: 'cloud' | 'self-hosted') => void;
}) {
  const { colors } = useAppTheme();

  return (
    <View style={[styles.modeSelector, { backgroundColor: colors.surfaceSubtle }]}>
      <Pressable
        style={[styles.modeTab, mode === 'cloud' && { backgroundColor: colors.surface }]}
        onPress={() => onModeChange('cloud')}
      >
        <Icon
          name="cloud-outline"
          size={16}
          color={mode === 'cloud' ? colors.primary : colors.textMuted}
        />
        <Text
          style={[styles.modeTabText, { color: mode === 'cloud' ? colors.text : colors.textMuted }]}
        >
          Codedock Cloud
        </Text>
      </Pressable>
      <Pressable
        style={[styles.modeTab, mode === 'self-hosted' && { backgroundColor: colors.surface }]}
        onPress={() => onModeChange('self-hosted')}
      >
        <Icon
          name="server-outline"
          size={16}
          color={mode === 'self-hosted' ? colors.primary : colors.textMuted}
        />
        <Text
          style={[
            styles.modeTabText,
            { color: mode === 'self-hosted' ? colors.text : colors.textMuted },
          ]}
        >
          Self-Hosted
        </Text>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  modeSelector: {
    flexDirection: 'row',
    borderRadius: 10,
    padding: 3,
    marginBottom: 16,
  },
  modeTab: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 6,
    paddingVertical: 8,
    borderRadius: 8,
  },
  modeTabText: {
    fontSize: 13,
    fontWeight: '600',
  },
});
