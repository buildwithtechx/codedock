import { Link, useRouterState } from '@tanstack/react-router';
import { ExternalLink } from 'lucide-react';
import type React from 'react';

export type NavItemProps = {
  title: string;
  url: string;
  icon: React.ComponentType<{ className?: string }>;
  search?: Record<string, string>;
  external?: boolean;
  badge?: string;
  count?: number | null;
};

export function isNavigationActive(url: string, pathname: string, exact?: boolean): boolean {
  if (url === '/') return pathname === '/';
  if (exact) return pathname === url;

  if (url === '/apps') {
    return (
      pathname === '/apps' || pathname.startsWith('/apps/') || pathname.startsWith('/services/')
    );
  }
  if (url === '/deployments') {
    return (
      pathname === '/deployments' ||
      pathname.startsWith('/deployments/') ||
      pathname.startsWith('/deploy/')
    );
  }
  if (url === '/monitoring') {
    return (
      pathname === '/monitoring' ||
      pathname.startsWith('/monitoring/') ||
      pathname.startsWith('/issues')
    );
  }
  if (url === '/settings') {
    return pathname === '/settings' || pathname.startsWith('/settings/');
  }

  return pathname === url || pathname.startsWith(`${url}/`);
}

export function NavItem({
  item,
  exact = false,
  collapsed = false,
}: {
  item: NavItemProps;
  exact?: boolean;
  collapsed?: boolean;
}) {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const currentSearch = routerState.location.search as Record<string, string | undefined>;
  const searchMatches = item.search
    ? Object.entries(item.search).every(([key, value]) => currentSearch[key] === value)
    : true;
  const isActive = isNavigationActive(item.url, pathname, exact) && searchMatches;
  const showCount = item.count != null && item.count > 0;

  return (
    <Link
      to={item.url as never}
      search={item.search as never}
      title={collapsed ? item.title : undefined}
      aria-current={isActive ? 'page' : undefined}
      className={`group relative flex items-center rounded-xl font-medium text-sm transition-colors ${
        collapsed ? 'justify-center px-0 py-2' : 'gap-3 px-3 py-2'
      } ${
        isActive
          ? 'bg-primary/12 text-sidebar-foreground'
          : 'text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-foreground'
      }`}
      target={item.external ? '_blank' : undefined}
      rel={item.external ? 'noopener noreferrer' : undefined}
    >
      {!collapsed && isActive && (
        <div className="absolute top-1/2 -left-3 h-4 w-0.5 -translate-y-1/2 rounded-r-full bg-primary" />
      )}

      <item.icon
        className={`h-4.5 w-4.5 shrink-0 transition-colors ${
          isActive ? 'text-primary' : 'text-muted-foreground group-hover:text-sidebar-foreground'
        }`}
      />

      {!collapsed && (
        <>
          <span className="flex-1 truncate">{item.title}</span>
          {item.external && <ExternalLink className="h-3 w-3 shrink-0 opacity-50" />}
          {item.badge && (
            <span className="rounded-full bg-primary/10 px-1.5 py-0.5 font-medium text-[9px] text-primary">
              {item.badge}
            </span>
          )}
          {showCount && (
            <span className="shrink-0 text-[13px] text-muted-foreground/60 tabular-nums">
              {item.count}
            </span>
          )}
        </>
      )}
    </Link>
  );
}
