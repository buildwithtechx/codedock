import { Link, useNavigate, useRouterState } from '@tanstack/react-router';
import {
  Bell,
  Database,
  Download,
  HardDrive,
  Lock,
  RefreshCw,
  Settings as SettingsIcon,
  Users,
  Wrench,
} from 'lucide-react';
import { BackupsList } from '#/features/backups/backups-list';
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import { OAuthProvidersList } from '#/features/users/oauth-providers-list';
import { GeneralSettings } from './general-settings';

type TabId = 'general' | 'notifications' | 'oauth' | 'backups';

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
  { id: 'backups', label: 'Backups', icon: <Database className="h-4 w-4" /> },
];

const instanceTools = [
  { label: 'Team members', to: '/users', icon: Users },
  { label: 'Storage', to: '/s3-destinations', icon: HardDrive },
  { label: 'Maintenance', to: '/maintenance', icon: Wrench },
  { label: 'Updates', to: '/updates', icon: RefreshCw },
  { label: 'Migration', to: '/migrations', icon: Download },
];

export const SettingsLayout = () => {
  const navigate = useNavigate();
  const search = useRouterState({ select: (state) => state.location.search as { tab?: TabId } });
  const activeId = TABS.some((tab) => tab.id === search.tab) ? (search.tab as TabId) : 'general';
  const activeTab = TABS.find((tab) => tab.id === activeId) || TABS[0];

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
    backups: <BackupsList />,
  }[activeId];

  return (
    <div className="space-y-6">
      <header>
        <p className="font-medium text-muted-foreground text-sm">Instance control</p>
        <h1 className="mt-1 font-semibold text-2xl tracking-tight">Settings</h1>
        <p className="mt-1 text-muted-foreground text-sm">
          Configure the Codedock instance and its integrations.
        </p>
      </header>

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

      <div className="grid min-w-0 gap-6 xl:grid-cols-[minmax(0,1fr)_17rem]">
        <section className="min-w-0">{content}</section>
        <aside className="hidden xl:block">
          <div className="sticky top-8 rounded-2xl border border-border/80 bg-card p-2 shadow-sm">
            <div className="flex items-center gap-3 px-3 py-3">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/12 text-primary">
                {activeTab.icon}
              </div>
              <div>
                <p className="font-semibold text-sm">{activeTab.label}</p>
                <p className="text-muted-foreground text-xs">Settings section</p>
              </div>
            </div>
            <div className="mt-1 space-y-1 border-border/70 border-t pt-2">
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
          </div>
          <div className="mt-3 rounded-2xl border border-border/80 bg-card p-2 shadow-sm">
            <p className="px-3 pt-2 pb-1.5 font-medium text-[10px] text-muted-foreground uppercase tracking-[0.14em]">
              Instance tools
            </p>
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
        </aside>
      </div>
    </div>
  );
};
