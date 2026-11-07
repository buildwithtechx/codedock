import { Tabs } from 'expo-router';
import { Icon } from '@/components/Icon';
import { useAppTheme } from '@/theme';

export default function TabLayout() {
  const { colors } = useAppTheme();

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: colors.tint,
        tabBarInactiveTintColor: colors.textMuted,
        headerStyle: {
          backgroundColor: colors.surface,
        },
        headerTintColor: colors.text,
        tabBarStyle: {
          backgroundColor: colors.tabBar,
          borderTopColor: colors.border,
        },
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: 'Workloads',
          tabBarIcon: ({ color }) => <Icon name="cube-outline" color={color as string} size={22} />,
        }}
      />
      <Tabs.Screen
        name="databases"
        options={{
          title: 'Databases',
          tabBarIcon: ({ color }) => (
            <Icon name="database-outline" color={color as string} size={22} />
          ),
        }}
      />
      <Tabs.Screen
        name="servers"
        options={{
          title: 'Servers',
          tabBarIcon: ({ color }) => (
            <Icon name="server-outline" color={color as string} size={22} />
          ),
        }}
      />
      <Tabs.Screen
        name="deploys"
        options={{
          title: 'Deployments',
          tabBarIcon: ({ color }) => (
            <Icon name="rocket-launch-outline" color={color as string} size={22} />
          ),
        }}
      />
      <Tabs.Screen
        name="settings"
        options={{
          title: 'Settings',
          tabBarIcon: ({ color }) => <Icon name="cog-outline" color={color as string} size={22} />,
        }}
      />
    </Tabs>
  );
}
