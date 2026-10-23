import { Link, useNavigate, useRouterState } from '@tanstack/react-router';
import {
  ArrowLeft,
  Bell,
  Bot,
  Brain,
  GitBranch,
  Globe,
  Lock,
  Network,
  Server,
  Settings as SettingsIcon,
  User,
  Wrench,
} from 'lucide-react';
import { useEffect } from 'react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { DnsSettings } from '#/features/dns/dns-settings';
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import { GithubIntegration, GitProviders } from '#/features/sources';
import { OAuthProvidersList } from '#/features/users/oauth-providers-list';
import { useAuthStore } from '#/stores/auth-store';
import { AISettings } from './ai-settings';
import { AppearanceEditor } from './appearance-editor/appearance-editor';
import { GeneralSettings } from './general-settings';
import { InstanceInfo } from './instance-info';
import { MaintenancePage } from './maintenance-settings';
import { McpSettings } from './mcp-settings';
import { MigrationSettings } from './migration-settings';
import { ProfileSettingsTab } from './profile-settings-tab';
import { ServerSettingsTab } from './server-settings-tab';
import { UpdatesPage } from './update-settings';

export type TabId =
  | 'profile'
  | 'general'
  | 'notifications'
  | 'dns'
  | 'git'
  | 'oauth'
  | 'ai'
  | 'mcp'
  | 'server'
  | 'maintenance'
  | 'instance';

interface TabItem {
  id: TabId;
  label: string;
  icon: React.ReactNode;
}

interface TabGroup {
  label: string;
  tabs: TabItem[];
}

const TAB_GROUPS: TabGroup[] = [
  {
    label: 'Workspace & Preferences',
    tabs: [
      { id: 'general', label: 'General', icon: <SettingsIcon className="h-4 w-4" /> },
      { id: 'profile', label: 'Profile & Security', icon: <User className="h-4 w-4" /> },
      { id: 'notifications', label: 'Notifications', icon: <Bell className="h-4 w-4" /> },
    ],
  },
  {
    label: 'Integrations & Access',
    tabs: [
      { id: 'dns', label: 'DNS Providers', icon: <Globe className="h-4 w-4" /> },
      { id: 'git', label: 'Git', icon: <GitBranch className="h-4 w-4" /> },
      { id: 'oauth', label: 'OAuth & SSO', icon: <Lock className="h-4 w-4" /> },
    ],
  },
  {
    label: 'AI & Automation',
    tabs: [
      { id: 'ai', label: 'AI Models', icon: <Brain className="h-4 w-4" /> },
      { id: 'mcp', label: 'MCP Servers', icon: <Bot className="h-4 w-4" /> },
    ],
  },
  {
    label: 'Instance & Fleet',
    tabs: [
      { id: 'server', label: 'Server & Routing', icon: <Network className="h-4 w-4" /> },
      { id: 'maintenance', label: 'Maintenance', icon: <Wrench className="h-4 w-4" /> },
      { id: 'instance', label: 'Instance & Updates', icon: <Server className="h-4 w-4" /> },
    ],
  },
];

const ALL_TABS = TAB_GROUPS.flatMap((g) => g.tabs);

export const SettingsLayout = () => {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const search = useRouterState({
    select: (state) => state.location.search as { tab?: string; view?: string; code?: string },
  });

  const rawTab = search.tab === 'credentials' ? 'dns' : search.tab;

  // Handle redirects for standalone pages
  useEffect(() => {
    if (rawTab === 'team') {
      void navigate({ to: '/team', replace: true });
    } else if (rawTab === 'tokens') {
      void navigate({ to: '/tokens', replace: true });
    }
  }, [rawTab, navigate]);

  const activeId: TabId = ALL_TABS.some((tab) => tab.id === rawTab) ? (rawTab as TabId) : 'general';
  const isAppearanceView = activeId === 'general' && search.view === 'appearance';

  const setActiveId = (tab: TabId) => {
    void navigate({
      to: '/settings',
      search: tab === 'general' ? {} : ({ tab } as never),
      replace: true,
    });
  };

  const content: Record<TabId, React.ReactNode> = {
    profile: <ProfileSettingsTab />,
    general: isAppearanceView ? (
      <div className="space-y-6">
        <Link
          to="/settings"
          search={{ tab: 'general' } as never}
          className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft className="size-4 rtl:rotate-180" />
          <span>Back to General</span>
        </Link>
        <AppearanceEditor />
      </div>
    ) : (
      <GeneralSettings />
    ),
    notifications: <NotificationsSettings />,
    dns: <DnsSettings />,
    git: (
      <div className="space-y-6">
        <GithubIntegration />
        <GitProviders />
      </div>
    ),
    oauth: <OAuthProvidersList />,
    ai: <AISettings />,
    mcp: <McpSettings />,
    server: <ServerSettingsTab />,
    maintenance: <MaintenancePage />,
    instance: (
      <div className="space-y-6">
        <InstanceInfo />
        <UpdatesPage />
        <MigrationSettings />
      </div>
    ),
  };

  return (
    <div className="space-y-5">
      <PageHeader
        title="Settings"
        description="Manage workspace preferences, developer integrations, AI tooling, and fleet configuration."
      />

      <div className="flex gap-1 overflow-x-auto xl:hidden">
        {ALL_TABS.map((tab) => {
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

            <div className="space-y-3 rounded-2xl border border-border/50 bg-card p-3">
              {TAB_GROUPS.map((group) => (
                <div key={group.label} className="space-y-1">
                  <p className="px-3 pt-2 pb-1 font-semibold text-[10px] text-muted-foreground/80 uppercase tracking-wider">
                    {group.label}
                  </p>
                  {group.tabs.map((tab) => {
                    const isActive = tab.id === activeId;
                    return (
                      <button
                        key={tab.id}
                        onClick={() => setActiveId(tab.id)}
                        type="button"
                        aria-current={isActive ? 'page' : undefined}
                        className={`flex w-full items-center gap-2.5 rounded-xl px-3 py-1.5 text-left font-medium text-[13px] transition-colors ${
                          isActive
                            ? 'bg-foreground/[0.08] text-foreground'
                            : 'text-muted-foreground hover:bg-foreground/[0.04] hover:text-foreground'
                        }`}
                      >
                        <span className={isActive ? 'text-primary' : 'text-muted-foreground'}>
                          {tab.icon}
                        </span>
                        <span>{tab.label}</span>
                      </button>
                    );
                  })}
                </div>
              ))}
            </div>
          </div>
        }
      >
        {content[activeId]}
      </PageFrame>
    </div>
  );
};
