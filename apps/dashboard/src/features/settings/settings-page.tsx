import { useNavigate, useRouterState } from '@tanstack/react-router';
import {
  Bell,
  Brain,
  GitBranch,
  KeyRound,
  Lock,
  Server,
  Settings as SettingsIcon,
  UsersRound,
  Wrench,
} from 'lucide-react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import { GithubIntegration, GitProviders } from '#/features/sources';
import { ApiKeysList } from '#/features/users/api-keys-list';
import { OAuthProvidersList } from '#/features/users/oauth-providers-list';
import { useAuthStore } from '#/stores/auth-store';
import { AISettings } from './ai-settings';
import { GeneralSettings } from './general-settings';
import { InstanceInfo } from './instance-info';
import { MaintenancePage } from './maintenance-settings';
import { MigrationSettings } from './migration-settings';
import { TeamSettings } from './team-settings';
import { UpdatesPage } from './update-settings';

type TabId =
  | 'general'
  | 'git'
  | 'tokens'
  | 'team'
  | 'notifications'
  | 'oauth'
  | 'ai'
  | 'maintenance'
  | 'instance';

type Tab = {
  id: TabId;
  label: string;
  icon: React.ReactNode;
};

const TABS: Tab[] = [
  {
    id: 'general',
    label: 'General',
    icon: <SettingsIcon className="h-4 w-4" />,
  },
  {
    id: 'git',
    label: 'Git',
    icon: <GitBranch className="h-4 w-4" />,
  },
  {
    id: 'tokens',
    label: 'Tokens',
    icon: <KeyRound className="h-4 w-4" />,
  },
  {
    id: 'team',
    label: 'Team',
    icon: <UsersRound className="h-4 w-4" />,
  },
  {
    id: 'notifications',
    label: 'Notifications',
    icon: <Bell className="h-4 w-4" />,
  },
  {
    id: 'oauth',
    label: 'OAuth',
    icon: <Lock className="h-4 w-4" />,
  },
  {
    id: 'ai',
    label: 'AI',
    icon: <Brain className="h-4 w-4" />,
  },
  {
    id: 'maintenance',
    label: 'Maintenance',
    icon: <Wrench className="h-4 w-4" />,
  },
  {
    id: 'instance',
    label: 'Instance',
    icon: <Server className="h-4 w-4" />,
  },
];

export const SettingsLayout = () => {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const search = useRouterState({
    select: (state) => state.location.search as { tab?: TabId; code?: string },
  });
  const activeId = TABS.some((tab) => tab.id === search.tab) ? (search.tab as TabId) : 'general';
  const setActiveId = (tab: TabId) => {
    void navigate({
      to: '/settings',
      search: tab === 'general' ? {} : ({ tab } as never),
      replace: true,
    });
  };

  const content = {
    general: <GeneralSettings />,
    git: (
      <Tabs defaultValue="github-app" className="pb-12">
        <TabsList>
          <TabsTrigger value="github-app">GitHub App</TabsTrigger>
          <TabsTrigger value="tokens">Personal tokens</TabsTrigger>
        </TabsList>
        <p className="mt-3 text-muted-foreground text-sm">
          Connect GitHub automatically with an App, or link providers manually with a personal
          access token.
        </p>
        <TabsContent value="github-app" className="mt-6">
          <GithubIntegration />
        </TabsContent>
        <TabsContent value="tokens" className="mt-6">
          <GitProviders />
        </TabsContent>
      </Tabs>
    ),
    tokens: <ApiKeysList />,
    team: <TeamSettings />,
    notifications: <NotificationsSettings />,
    oauth: <OAuthProvidersList />,
    ai: <AISettings />,
    maintenance: <MaintenancePage />,
    instance: (
      <div className="space-y-6">
        <InstanceInfo />
        <UpdatesPage />
        <MigrationSettings />
      </div>
    ),
  }[activeId];
  return (
    <div className="space-y-5">
      <PageHeader
        title="Settings"
        description="Manage instance behavior, connections, and notifications."
      />

      <div className="flex gap-1 overflow-x-auto xl:hidden">
        {TABS.map((tab) => {
          const isActive = tab.id === activeId;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveId(tab.id)}
              type="button"
              aria-current={isActive ? 'page' : undefined}
              className={`flex shrink-0 items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors ${
                isActive
                  ? 'bg-primary/12 font-medium text-foreground'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground'
              }`}
            >
              {tab.icon}
              {tab.label}
            </button>
          );
        })}
      </div>

      <PageFrame
        rail={
          <div>
            <section className="mb-3 rounded-2xl bg-card p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/12 text-primary">
                  <SettingsIcon className="h-4 w-4" />
                </div>
                <div className="min-w-0">
                  <p className="font-semibold text-sm">Settings</p>
                  <p className="truncate text-muted-foreground text-xs">
                    {user?.email || 'Instance owner'}
                  </p>
                </div>
              </div>
            </section>
            <div className="rounded-2xl bg-card p-2">
              <div className="space-y-1">
                {TABS.map((tab) => {
                  const isActive = tab.id === activeId;
                  return (
                    <button
                      key={tab.id}
                      onClick={() => setActiveId(tab.id)}
                      type="button"
                      aria-current={isActive ? 'page' : undefined}
                      className={`flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left font-medium text-sm transition-colors ${
                        isActive
                          ? 'bg-primary/12 text-foreground'
                          : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'
                      }`}
                    >
                      {tab.icon}
                      {tab.label}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        }
      >
        {content}
      </PageFrame>
    </div>
  );
};
