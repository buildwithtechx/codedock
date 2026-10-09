import { Link } from '@tanstack/react-router';
import { Moon, PanelLeftClose, PanelLeftOpen, Plus, Sun, X } from 'lucide-react';
import { useTheme } from 'next-themes';
import {
  infrastructureNavigation,
  primaryNavigation,
  systemNavigation,
} from './dashboard-navigation';
import { NavItem, type NavItemProps } from './nav-item';
import { OrganizationSwitcher } from './organization-switcher';
import { useSidebarCounts } from './use-sidebar-counts';

type NavGroup = {
  title?: string;
  items: (NavItemProps & { exact?: boolean })[];
};

const navGroups: NavGroup[] = [
  {
    title: 'Main',
    items: primaryNavigation.map(({ title, to, icon, exact }) => ({
      title,
      url: to,
      icon,
      exact,
    })),
  },
  {
    title: 'Infrastructure',
    items: infrastructureNavigation.map(({ title, to, icon, exact }) => ({
      title,
      url: to,
      icon,
      exact,
    })),
  },
  {
    title: 'Settings',
    items: systemNavigation.map(({ title, to, icon, exact }) => ({
      title,
      url: to,
      icon,
      exact,
    })),
  },
];

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
  const { resolvedTheme, setTheme } = useTheme();
  const counts = useSidebarCounts();

  const toggleTheme = () => setTheme(resolvedTheme === 'dark' ? 'light' : 'dark');

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
        className={`fixed inset-y-3 left-3 z-40 flex flex-col overflow-hidden rounded-2xl border border-sidebar-border/70 bg-sidebar shadow-2xl shadow-black/15 transition-all duration-300 md:z-20 ${
          collapsed ? 'md:w-[72px]' : 'md:w-[260px]'
        } ${mobileOpen ? 'w-[min(19rem,calc(100vw-1.5rem))] translate-x-0' : 'w-[min(19rem,calc(100vw-1.5rem))] -translate-x-[calc(100%+1rem)] md:translate-x-0'}`}
      >
        <div className="flex items-center justify-between px-3 pt-3 pb-2 md:hidden">
          <div className="flex items-center gap-2.5 py-2">
            <img src="/apple-touch-icon.png" alt="" className="h-7 w-7 shrink-0 rounded-lg" />
            <span className="truncate font-medium text-sidebar-foreground text-sm">Codedock</span>
          </div>
          <button
            type="button"
            onClick={onMobileClose}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="hidden md:block">
          <div className={navCollapsed ? 'px-2 py-3' : 'px-5 py-5'}>
            <div
              className={`flex ${navCollapsed ? 'flex-col items-center gap-2' : 'items-center justify-between gap-2.5 py-2'}`}
            >
              <div className="flex min-w-0 items-center gap-2.5">
                <img src="/apple-touch-icon.png" alt="" className="h-7 w-7 shrink-0 rounded-lg" />
                {!navCollapsed && (
                  <span className="flex-1 truncate font-medium text-sidebar-foreground text-sm">
                    Codedock
                  </span>
                )}
              </div>
              <div className={`flex items-center ${navCollapsed ? 'flex-col gap-1' : 'gap-1'}`}>
                <button
                  type="button"
                  onClick={toggleTheme}
                  className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  aria-label="Toggle theme"
                >
                  {resolvedTheme === 'dark' ? (
                    <Sun className="h-4 w-4" />
                  ) : (
                    <Moon className="h-4 w-4" />
                  )}
                </button>
                <button
                  type="button"
                  onClick={onToggle}
                  className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
                >
                  {collapsed ? (
                    <PanelLeftOpen className="h-4 w-4" />
                  ) : (
                    <PanelLeftClose className="h-4 w-4" />
                  )}
                </button>
              </div>
            </div>
          </div>
          <div className="mx-3 h-px bg-sidebar-border" />
        </div>

        <div className="relative min-h-0 flex-1">
          <nav className="h-full space-y-5 overflow-y-auto px-3 pt-3 pb-12 [scrollbar-width:thin]">
            {navGroups.map((group, i) => (
              <div key={group.title ?? i} className="flex flex-col gap-0.5">
                {!navCollapsed && group.title && (
                  <h4 className="px-2 pb-1.5 font-semibold text-[11px] text-sidebar-foreground/50 uppercase tracking-[0.14em]">
                    {group.title}
                  </h4>
                )}
                {navCollapsed && i > 0 && <div className="mx-2 my-3 h-px bg-sidebar-border/60" />}
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
            ))}
          </nav>
          <div className="pointer-events-none absolute inset-x-0 bottom-0 h-14 bg-gradient-to-t from-sidebar to-transparent" />
        </div>

        <div className="px-3 pb-2">
          <Link
            to={'/library' as never}
            aria-label="New project"
            title={navCollapsed ? 'New project' : undefined}
            className={`flex items-center justify-center gap-2.5 overflow-hidden rounded-xl bg-gradient-to-r from-violet-600 to-fuchsia-600 px-3 py-2.5 font-semibold text-sm text-white transition-all hover:brightness-110 ${navCollapsed ? 'px-0' : 'px-3'}`}
          >
            <Plus className="h-4 w-4" />
            {!navCollapsed && 'New project'}
          </Link>
        </div>

        <div className={`mt-auto px-3 pt-1 pb-3 ${navCollapsed ? 'px-2' : ''}`}>
          <div className="mx-2 mb-3 h-px bg-sidebar-border/60" />
          {!navCollapsed && (
            <p className="mb-2 px-2 font-semibold text-[11px] text-sidebar-foreground/50 uppercase tracking-[0.14em]">
              Account
            </p>
          )}
          <OrganizationSwitcher collapsed={navCollapsed} />
        </div>
      </aside>
    </>
  );
}
