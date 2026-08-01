import { useNavigate, useRouterState } from '@tanstack/react-router';
import { Bell, Lock, Settings as SettingsIcon } from 'lucide-react';
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import { OAuthProvidersList } from '#/features/users/oauth-providers-list';
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

export const SettingsLayout = () => {
  const navigate = useNavigate();
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
          </div>
        </aside>
      </div>
    </div>
  );
};
