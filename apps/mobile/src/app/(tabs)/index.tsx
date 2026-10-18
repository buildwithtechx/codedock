import { useRouter } from 'expo-router';
import { useState } from 'react';
import {
  ActivityIndicator,
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { Icon } from '@/components/Icon';
import { ProjectCard } from '@/features/projects/ProjectCard';
import { AppCard } from '@/features/workloads/AppCard';
import { useProjectsQuery } from '@/hooks/useProjects';
import { useAppsQuery } from '@/hooks/useWorkloads';
import type { CodedockApp, CodedockProject } from '@/interfaces';
import { useAppStore } from '@/store/useAppStore';
import { useAppTheme } from '@/theme';

export default function WorkloadsScreen() {
  const router = useRouter();
  const { colors } = useAppTheme();
  const { connection, setSelectedApp, selectedProject, setSelectedProject } = useAppStore();

  const [activeTab, setActiveTab] = useState<'projects' | 'apps'>('projects');

  const {
    data: projects = [],
    isLoading: loadingProjects,
    isRefetching: refetchingProjects,
    refetch: refetchProjects,
  } = useProjectsQuery();

  const {
    data: allApps = [],
    isLoading: loadingApps,
    isRefetching: refetchingApps,
    refetch: refetchApps,
  } = useAppsQuery();

  const displayedApps = selectedProject
    ? allApps.filter(
        (a) => a.projectId === selectedProject.id || a.projectName === selectedProject.name
      )
    : allApps;

  const onSelectApp = (app: CodedockApp) => {
    setSelectedApp(app);
    router.push('/modal');
  };

  const onSelectProject = (project: CodedockProject) => {
    setSelectedProject(project);
    setActiveTab('apps');
  };

  if (!connection) {
    return (
      <View style={[styles.center, { backgroundColor: colors.background }]}>
        <Icon name="shield-lock-outline" size={40} color={colors.primary} />
        <Text style={[styles.title, { color: colors.text }]}>Authentication Required</Text>
        <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
          Sign in to access your projects and apps.
        </Text>
        <Pressable
          style={[styles.button, { backgroundColor: colors.primary }]}
          onPress={() => router.replace('/')}
        >
          <Text style={styles.buttonText}>Sign In / Connect</Text>
        </Pressable>
      </View>
    );
  }

  const isRefreshing = activeTab === 'projects' ? refetchingProjects : refetchingApps;
  const onRefresh = () => {
    if (activeTab === 'projects') void refetchProjects();
    else void refetchApps();
  };

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <View style={[styles.segmentedControl, { backgroundColor: colors.surfaceSubtle }]}>
        <Pressable
          style={[
            styles.segmentBtn,
            activeTab === 'projects' && { backgroundColor: colors.surface },
          ]}
          onPress={() => setActiveTab('projects')}
        >
          <Icon
            name="folder-outline"
            size={16}
            color={activeTab === 'projects' ? colors.primary : colors.textMuted}
          />
          <Text
            style={[
              styles.segmentText,
              { color: activeTab === 'projects' ? colors.text : colors.textMuted },
            ]}
          >
            Projects ({projects.length})
          </Text>
        </Pressable>

        <Pressable
          style={[styles.segmentBtn, activeTab === 'apps' && { backgroundColor: colors.surface }]}
          onPress={() => setActiveTab('apps')}
        >
          <Icon
            name="cube-outline"
            size={16}
            color={activeTab === 'apps' ? colors.primary : colors.textMuted}
          />
          <Text
            style={[
              styles.segmentText,
              { color: activeTab === 'apps' ? colors.text : colors.textMuted },
            ]}
          >
            Apps ({displayedApps.length})
          </Text>
        </Pressable>
      </View>

      {selectedProject && activeTab === 'apps' && (
        <View
          style={[
            styles.filterBar,
            { backgroundColor: colors.surface, borderColor: colors.border },
          ]}
        >
          <Text style={[styles.filterText, { color: colors.textSecondary }]}>
            Project:{' '}
            <Text style={{ color: colors.text, fontWeight: '600' }}>{selectedProject.name}</Text>
          </Text>
          <Pressable onPress={() => setSelectedProject(null)} style={styles.clearFilterBtn}>
            <Text style={[styles.clearFilterText, { color: colors.primary }]}>Show all apps</Text>
          </Pressable>
        </View>
      )}

      {activeTab === 'projects' ? (
        loadingProjects ? (
          <ActivityIndicator style={{ marginTop: 32 }} color={colors.primary} />
        ) : (
          <FlatList
            data={projects}
            keyExtractor={(item) => item.id}
            refreshControl={
              <RefreshControl
                refreshing={isRefreshing}
                onRefresh={onRefresh}
                tintColor={colors.primary}
              />
            }
            ListEmptyComponent={
              <View style={styles.center}>
                <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
                  No projects found on this instance.
                </Text>
              </View>
            }
            renderItem={({ item }) => (
              <ProjectCard project={item} onPress={() => onSelectProject(item)} />
            )}
          />
        )
      ) : loadingApps ? (
        <ActivityIndicator style={{ marginTop: 32 }} color={colors.primary} />
      ) : (
        <FlatList
          data={displayedApps}
          keyExtractor={(item) => item.id}
          refreshControl={
            <RefreshControl
              refreshing={isRefreshing}
              onRefresh={onRefresh}
              tintColor={colors.primary}
            />
          }
          ListEmptyComponent={
            <View style={styles.center}>
              <Text style={[styles.subtitle, { color: colors.textSecondary }]}>No apps found.</Text>
            </View>
          }
          renderItem={({ item }) => <AppCard app={item} onPress={() => onSelectApp(item)} />}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 16 },
  center: { flex: 1, justifyContent: 'center', alignItems: 'center', padding: 24, marginTop: 40 },
  title: { fontSize: 20, fontWeight: 'bold', marginTop: 12, marginBottom: 8 },
  subtitle: { fontSize: 14, textAlign: 'center', marginBottom: 20 },
  button: { paddingVertical: 12, paddingHorizontal: 24, borderRadius: 10 },
  buttonText: { color: '#FFFFFF', fontWeight: '600' },
  segmentedControl: { flexDirection: 'row', borderRadius: 10, padding: 4, marginBottom: 12 },
  segmentBtn: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 6,
    paddingVertical: 8,
    borderRadius: 8,
  },
  segmentText: { fontSize: 13, fontWeight: '600' },
  filterBar: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 8,
    marginBottom: 12,
  },
  filterText: { fontSize: 13 },
  clearFilterBtn: { padding: 4 },
  clearFilterText: { fontSize: 12, fontWeight: '600' },
});
