import { useNavigate, useRouterState } from '@tanstack/react-router';
import {
  ArrowRightLeft,
  Bell,
  Brain,
  CloudCog,
  Download,
  Lock,
  Settings as SettingsIcon,
  UsersRound,
  Wrench,
} from 'lucide-react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import { GithubIntegration, GitProviders } from '#/features/sources';
import { OAuthProvidersList } from '#/features/users/oauth-providers-list';
import { useAuthStore } from '#/stores/auth-store';
import { AISettings } from './ai-settings';
import { GeneralSettings } from './general-settings';
import { MaintenancePage } from './maintenance-settings';
import { MigrationSettings } from './migration-settings';
import { TeamSettings } from './team-settings';
import { UpdatesPage } from './update-settings';

type TabId =
  | 'general'
  | 'notifications'
  | 'oauth'
  | 'team'
  | 'ai'
  | 'sources'
  | 'maintenance'
  | 'updates'
  | 'migration';

type Tab = {
  id: TabId;
  label: string;
  icon: React.ReactNode;
  group: string;
};

const TABS: Tab[] = [
  {
    id: 'general',
    label: 'General',
    icon: <SettingsIcon className="h-4 w-4" />,
    group: 'General',
  },
  {
    id: 'notifications',
    label: 'Notifications',
    icon: <Bell className="h-4 w-4" />,
    group: 'General',
  },
  {
    id: 'oauth',
    label: 'OAuth',
    icon: <Lock className="h-4 w-4" />,
    group: 'General',
  },
  {
    id: 'team',
    label: 'Team',
    icon: <UsersRound className="h-4 w-4" />,
    group: 'General',
  },
  {
    id: 'ai',
    label: 'AI',
    icon: <Brain className="h-4 w-4" />,
    group: 'Integrations',
  },
  {
    id: 'sources',
    label: 'Sources',
    icon: <CloudCog className="h-4 w-4" />,
    group: 'Integrations',
  },
  {
    id: 'maintenance',
    label: 'Maintenance',
    icon: <Wrench className="h-4 w-4" />,
    group: 'Operations',
  },
  {
    id: 'updates',
    label: 'Updates',
    icon: <Download className="h-4 w-4" />,
    group: 'Operations',
  },
  {
    id: 'migration',
    label: 'Migration',
    icon: <ArrowRightLeft className="h-4 w-4" />,
    group: 'Operations',
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
    });
  };

  const content = {
    general: <GeneralSettings />,
    notifications: <NotificationsSettings />,
    oauth: <OAuthProvidersList />,
    team: <TeamSettings />,
    ai: <AISettings />,
    sources: (
      <div className="space-y-8 pb-12">
        <GithubIntegration />
        <GitProviders />
      </div>
    ),
    maintenance: <MaintenancePage />,
    updates: <UpdatesPage />,
    migration: <MigrationSettings />,
  }[activeId];
  const tabsByGroup = TABS.reduce<Record<string, Tab[]>>((groups, tab) => {
    groups[tab.group] ??= [];
    groups[tab.group].push(tab);
    return groups;
  }, {});

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
            <section className="mb-3 rounded-2xl border border-border/80 bg-card p-4 shadow-sm">
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
            <div className="rounded-2xl border border-border/80 bg-card p-2 shadow-sm">
              {Object.entries(tabsByGroup).map(([group, tabs], index) => (
                <div key={group} className={index === 0 ? 'space-y-1' : 'mt-4 space-y-1'}>
                  <p className="px-3 pb-1 font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.16em]">
                    {group}
                  </p>
                  {tabs.map((tab) => {
                    const isActive = tab.id === activeId;
                    return (
                      <button
                        key={tab.id}
                        onClick={() => setActiveId(tab.id)}
                        type="button"
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
              ))}
            </div>
          </div>
        }
      >
        {content}
      </PageFrame>
    </div>
  );
};
