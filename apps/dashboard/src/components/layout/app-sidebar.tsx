import { Link } from '@tanstack/react-router';
import { LogOut, Moon, PanelLeftClose, PanelLeftOpen, Plus, Sun, X } from 'lucide-react';
import { useTheme } from 'next-themes';
import { useLogout } from '#/features/auth';
import { useGetPublicSettings } from '#/features/settings';
import {
  getSystemNavigation,
  helpNavigation,
  infrastructureNavigation,
  primaryNavigation,
} from './dashboard-navigation';
import { NavItem, type NavItemProps } from './nav-item';
import { OrganizationSwitcher } from './organization-switcher';
import { useSidebarCounts } from './use-sidebar-counts';

type NavGroup = {
  title?: string;
  items: (NavItemProps & { exact?: boolean })[];
};

const countForUrl: Record<string, keyof ReturnType<typeof useSidebarCounts>> = {
  '/projects': 'projects',
  '/apps': 'apps',
  '/servers': 'servers',
};

interface AppSidebarProps {
  collapsed: boolean;
  onToggle: () => void;
  mobileOpen: boolean;
  onMobileClose: () => void;
}

export function AppSidebar({ collapsed, onToggle, mobileOpen, onMobileClose }: AppSidebarProps) {
  const navCollapsed = collapsed && !mobileOpen;
  const { theme, resolvedTheme, setTheme } = useTheme();
  const counts = useSidebarCounts();
  const { data: publicRes } = useGetPublicSettings();
  const { mutate: logout, isPending: isLoggingOut } = useLogout();
  const isCloud = publicRes?.data?.cloudMode ?? false;

  const mainGroup: NavGroup = {
    title: 'Main',
    items: primaryNavigation.map(({ title, to, icon, exact }) => ({
      title,
      url: to,
      icon,
      exact,
    })),
  };

  const infraGroup: NavGroup = {
    title: 'Infrastructure',
    items: infrastructureNavigation.map(({ title, to, icon, exact }) => ({
      title,
      url: to,
      icon,
      exact,
    })),
  };

  const settingsGroup: NavGroup = {
    title: 'Settings',
    items: getSystemNavigation(isCloud).map(({ title, to, icon, exact }) => ({
      title,
      url: to,
      icon,
      exact,
    })),
  };

  const helpGroup: NavGroup = {
    title: 'Help',
    items: helpNavigation.map(({ title, to, icon, exact }) => ({
      title,
      url: to,
      icon,
      exact,
    })),
  };

  const navGroups: NavGroup[] = isCloud
    ? [mainGroup, settingsGroup, infraGroup, helpGroup]
    : [mainGroup, infraGroup, settingsGroup, helpGroup];

  const toggleTheme = () => {
    if (theme === 'light') setTheme('dim');
    else if (theme === 'dim') setTheme('dark');
    else if (theme === 'dark') setTheme('light');
    else setTheme(resolvedTheme === 'dark' ? 'light' : 'dim');
  };

  return (
    <>
      {mobileOpen && (
        <button
          type="button"
          className="fixed inset-0 z-30 cursor-default bg-black/50 md:hidden"
          onClick={onMobileClose}
          aria-label="Close menu"
        />
      )}

      <aside
        id="dashboard-sidebar"
        className={`z-40 my-3 ms-3 flex shrink-0 flex-col overflow-hidden rounded-2xl border border-border/50 bg-card shadow-2xl shadow-black/15 transition-[width] duration-200 md:static md:z-auto md:shadow-none ${
          collapsed ? 'md:w-[72px]' : 'md:w-[260px]'
        } fixed inset-y-3 left-3 ${
          mobileOpen
            ? 'w-[min(19rem,calc(100vw-1.5rem))] translate-x-0'
            : 'w-[min(19rem,calc(100vw-1.5rem))] -translate-x-[calc(100%+1rem)] md:translate-x-0'
        }`}
      >
        <div className="flex items-center justify-between px-3 pt-3 pb-2 md:hidden">
          <div className="flex items-center gap-2.5 py-2">
            <img src="/apple-touch-icon.png" alt="" className="size-6.5 shrink-0 rounded-lg" />
            <span className="truncate font-semibold tracking-tight text-foreground text-sm">Codedock</span>
          </div>
          <button
            type="button"
            onClick={onMobileClose}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-foreground/[0.06] hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="hidden md:block">
          <div
            className={`app-sidebar-header flex items-center px-5 py-6 ${
              navCollapsed ? 'flex-col gap-3 pb-3' : 'justify-between'
            }`}
          >
            <div className="flex min-w-0 items-center gap-2.5">
              <img src="/apple-touch-icon.png" alt="" className="size-6.5 shrink-0 rounded-lg" />
              {!navCollapsed && (
                <span className="truncate text-base font-semibold tracking-tight text-foreground">
                  Codedock
                </span>
              )}
            </div>
            <div className={`flex items-center ${navCollapsed ? 'flex-col gap-1' : 'gap-1'}`}>
              <button
                type="button"
                onClick={toggleTheme}
                className="flex size-8 items-center justify-center rounded-xl text-muted-foreground transition-colors hover:bg-foreground/[0.06] hover:text-foreground"
                aria-label="Toggle theme"
                title="Toggle theme"
              >
                {theme === 'dark' || (theme === 'system' && resolvedTheme === 'dark') ? (
                  <Sun className="size-4" />
                ) : theme === 'dim' ? (
                  <Moon className="size-4 text-primary" />
                ) : (
                  <Moon className="size-4" />
                )}
              </button>
              <button
                type="button"
                onClick={onToggle}
                className="flex size-8 items-center justify-center rounded-xl text-muted-foreground transition-colors hover:bg-foreground/[0.06] hover:text-foreground"
                aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
                title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
              >
                {collapsed ? (
                  <PanelLeftOpen className="size-4" />
                ) : (
                  <PanelLeftClose className="size-4" />
                )}
              </button>
            </div>
          </div>
          <div className="mx-3 h-px bg-border/60" />
        </div>

        <div className="relative min-h-0 flex-1">
          <nav className="scrollbar-none h-full space-y-5 overflow-y-auto px-3 pt-3 pb-12 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
            {navGroups.map((group, i) => (
              <div key={group.title ?? i} className={i > 0 ? 'mt-5' : undefined}>
                {!navCollapsed && group.title && (
                  <p className="mb-2 px-2 text-[11px] font-semibold uppercase tracking-widest text-muted-foreground/60">
                    {group.title}
                  </p>
                )}
                {navCollapsed && i > 0 && <div className="mx-2 my-3 h-px bg-border/60" />}
                <div className="space-y-1">
                  {group.items.map((item) => {
                    const countKey = countForUrl[item.url];
                    return (
                      <NavItem
                        key={item.url}
                        item={{ ...item, count: countKey ? counts[countKey] : null }}
                        exact={item.exact}
                        collapsed={navCollapsed}
                      />
                    );
                  })}
                </div>
              </div>
            ))}
          </nav>
          <div className="pointer-events-none absolute inset-x-0 bottom-0 h-14 bg-gradient-to-t from-card to-transparent" />
        </div>

        <div className="px-3 pb-2">
          <Link
            to={'/library' as never}
            aria-label="New project"
            title={navCollapsed ? 'New project' : undefined}
            className="th-btn-accent flex items-center justify-center gap-2.5 rounded-xl px-3 py-2.5 text-sm font-semibold transition-all overflow-hidden hover:brightness-110"
          >
            <Plus className="size-4" />
            {!navCollapsed && <span>New project</span>}
          </Link>
        </div>

        <div
          className={`mt-auto px-3 pt-1 pb-3 ${navCollapsed ? 'flex flex-col items-center px-2' : ''}`}
        >
          <div className="mx-2 mb-3 h-px bg-border/60" />
          {!navCollapsed && (
            <p className="mb-2 px-2 text-[11px] font-semibold uppercase tracking-widest text-muted-foreground/60">
              Account
            </p>
          )}
          <OrganizationSwitcher collapsed={navCollapsed} />
          {navCollapsed && (
            <button
              type="button"
              onClick={() => logout()}
              disabled={isLoggingOut}
              className="mt-1.5 flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-foreground/[0.06] hover:text-foreground disabled:opacity-50"
              aria-label="Log out"
              title="Log out"
            >
              <LogOut className="h-4 w-4" />
            </button>
          )}
        </div>
      </aside>
    </>
  );
}
