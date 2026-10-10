import { useNavigate, useRouterState } from '@tanstack/react-router';
import {
  Bell,
  Bot,
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
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import { GithubIntegration, GitProviders } from '#/features/sources';
import { ApiKeysList } from '#/features/users/api-keys-list';
import { OAuthProvidersList } from '#/features/users/oauth-providers-list';
import { useAuthStore } from '#/stores/auth-store';
import { AISettings } from './ai-settings';
import { CredentialsSettings } from './credentials-settings';
import { GeneralSettings } from './general-settings';
import { InstanceInfo } from './instance-info';
import { MaintenancePage } from './maintenance-settings';
import { McpSettings } from './mcp-settings';
import { MigrationSettings } from './migration-settings';
import { TeamSettings } from './team-settings';
import { UpdatesPage } from './update-settings';

type TabId =
  | 'general'
  | 'credentials'
  | 'mcp'
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
    id: 'credentials',
    label: 'Credentials',
    icon: <KeyRound className="h-4 w-4" />,
  },
  {
    id: 'mcp',
    label: 'MCP',
    icon: <Bot className="h-4 w-4" />,
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
    credentials: <CredentialsSettings />,
    mcp: <McpSettings />,
    git: (
      <div className="space-y-6">
        <GithubIntegration />
        <GitProviders />
      </div>
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
        className="lg:grid-cols-[minmax(0,1fr)_280px]"
        rail={
          <div className="space-y-3">
            <section className="rounded-2xl border border-border/50 bg-card p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-foreground">
                  <SettingsIcon className="h-4 w-4" />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-foreground text-sm">Settings</p>
                  <p className="truncate text-muted-foreground text-xs">
                    {user?.email || 'Instance owner'}
                  </p>
                </div>
              </div>
            </section>
            <div className="rounded-2xl border border-border/50 bg-card p-3">
              <div className="space-y-1">
                {TABS.map((tab) => {
                  const isActive = tab.id === activeId;
                  return (
                    <button
                      key={tab.id}
                      onClick={() => setActiveId(tab.id)}
                      type="button"
                      aria-current={isActive ? 'page' : undefined}
                      className={`flex w-full items-center gap-3 rounded-xl px-3 py-2 text-left font-medium text-[14px] transition-colors ${
                        isActive
                          ? 'bg-foreground/[0.07] text-foreground'
                          : 'text-muted-foreground hover:bg-foreground/[0.04] hover:text-foreground'
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
