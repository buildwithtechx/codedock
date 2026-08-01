import { Link, useNavigate, useRouterState } from '@tanstack/react-router';
import {
  ArrowRightLeft,
  Bell,
  Brain,
  CloudCog,
  Download,
  Lock,
  Settings as SettingsIcon,
  Wrench,
} from 'lucide-react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import { OAuthProvidersList } from '#/features/users/oauth-providers-list';
import { useAuthStore } from '#/stores/auth-store';
import { GeneralSettings } from './general-settings';

type TabId = 'general' | 'notifications' | 'oauth';

type Tab = { id: TabId; label: string; icon: React.ReactNode };

const TABS: Tab[] = [
  {
    id: 'general',
    label: 'General',
    icon: <SettingsIcon className="h-4 w-4" />,
  },
  {
    id: 'notifications',
    label: 'Notifications',
    icon: <Bell className="h-4 w-4" />,
  },
  { id: 'oauth', label: 'OAuth', icon: <Lock className="h-4 w-4" /> },
];

const instanceTools = [
  { label: 'AI', to: '/ai', icon: Brain },
  { label: 'Sources', to: '/sources', icon: CloudCog },
  { label: 'Maintenance', to: '/maintenance', icon: Wrench },
  { label: 'Updates', to: '/updates', icon: Download },
  { label: 'Migration', to: '/migrations', icon: ArrowRightLeft },
];

export const SettingsLayout = () => {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const search = useRouterState({ select: (state) => state.location.search as { tab?: TabId } });
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

      <nav aria-label="Instance tools" className="flex gap-1 overflow-x-auto xl:hidden">
        {instanceTools.map((tool) => (
          <Link
            key={tool.to}
            to={tool.to}
            className="flex shrink-0 items-center gap-2 rounded-lg border border-border/80 bg-card px-3 py-2 text-muted-foreground text-sm transition-colors hover:bg-muted hover:text-foreground"
          >
            <tool.icon className="h-4 w-4" />
            {tool.label}
          </Link>
        ))}
      </nav>

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
              <div className="space-y-1">
                {TABS.map((tab) => {
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
              <div className="my-2 h-px bg-border/70" />
              <div className="space-y-1">
                {instanceTools.map((tool) => (
                  <Link
                    key={tool.to}
                    to={tool.to}
                    className="flex items-center gap-3 rounded-xl px-3 py-2.5 font-medium text-muted-foreground text-sm transition-colors hover:bg-muted/60 hover:text-foreground"
                  >
                    <tool.icon className="h-4 w-4" />
                    {tool.label}
                  </Link>
                ))}
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
