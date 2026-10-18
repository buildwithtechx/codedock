import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Icon } from '@/components/Icon';
import type { CodedockProject } from '@/interfaces';
import { useAppTheme } from '@/theme';

export function ProjectCard({
  project,
  onPress,
}: {
  project: CodedockProject;
  onPress: () => void;
}) {
  const { colors } = useAppTheme();

  return (
    <Pressable
      style={[styles.card, { backgroundColor: colors.surface, borderColor: colors.border }]}
      onPress={onPress}
    >
      <View style={styles.cardHeader}>
        <View style={styles.titleRow}>
          <Icon name="folder" size={18} color={colors.primary} />
          <Text style={[styles.cardTitle, { color: colors.text }]}>{project.name}</Text>
        </View>
        <Icon name="chevron-right" size={18} color={colors.textMuted} />
      </View>
      {project.gitUrl ? (
        <Text style={[styles.cardSub, { color: colors.textMuted }]}>
          {project.gitUrl} {project.gitBranch ? `(${project.gitBranch})` : ''}
        </Text>
      ) : (
        <Text style={[styles.cardSub, { color: colors.textMuted }]}>Standard Project</Text>
      )}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    borderWidth: 1,
    borderRadius: 10,
    padding: 16,
    marginBottom: 10,
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
    fontWeight: '600',
  },
  cardSub: {
    fontSize: 13,
  },
});
