import { Link, useRouterState } from '@tanstack/react-router';
import { ExternalLink } from 'lucide-react';
import type React from 'react';
import { Tooltip, TooltipContent, TooltipTrigger } from '@codedock/ui';

export type NavItemProps = {
  title: string;
  url: string;
  icon: React.ComponentType<{ className?: string }>;
  search?: Record<string, string>;
  exact?: boolean;
  external?: boolean;
  badge?: string;
  count?: number | null;
  statusDot?: 'healthy' | 'warning' | 'error' | 'running' | null;
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

function StatusIndicator({ status }: { status: NavItemProps['statusDot'] }) {
  if (!status) return null;
  const colors = {
    healthy: 'bg-emerald-500',
    warning: 'bg-amber-500',
    error: 'bg-rose-500',
    running: 'bg-sky-500 animate-pulse',
  };
  return <span className={`size-2 shrink-0 rounded-full ${colors[status]}`} />;
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

  const content = (
    <Link
      to={item.url as never}
      search={item.search as never}
      aria-current={isActive ? 'page' : undefined}
      className={`group relative flex items-center rounded-lg font-medium text-[13px] transition-colors ${
        collapsed ? 'justify-center px-0 py-2' : 'gap-2.5 px-2.5 py-1.5'
      } ${
        isActive
          ? 'bg-sidebar-accent text-sidebar-foreground font-medium'
          : 'text-sidebar-foreground/75 hover:bg-sidebar-accent/60 hover:text-sidebar-foreground'
      }`}
      target={item.external ? '_blank' : undefined}
      rel={item.external ? 'noopener noreferrer' : undefined}
    >
      {!collapsed && isActive && (
        <div className="absolute top-1/2 -left-2.5 h-4 w-1 -translate-y-1/2 rounded-r-full bg-primary" />
      )}

      <div className="relative">
        <item.icon
          className={`size-4 shrink-0 transition-colors ${
            isActive ? 'text-primary' : 'text-muted-foreground group-hover:text-sidebar-foreground'
          }`}
        />
        {collapsed && item.statusDot && (
          <span className="absolute -top-0.5 -right-0.5">
            <StatusIndicator status={item.statusDot} />
          </span>
        )}
      </div>

      {!collapsed && (
        <>
          <span className="flex-1 truncate">{item.title}</span>
          {item.statusDot && <StatusIndicator status={item.statusDot} />}
          {item.external && <ExternalLink className="h-3 w-3 shrink-0 opacity-50" />}
          {item.badge && (
            <span className="rounded-full bg-primary/10 px-1.5 py-0.5 font-medium text-[9px] text-primary">
              {item.badge}
            </span>
          )}
          {showCount && (
            <span className="shrink-0 text-[12px] text-muted-foreground/60 tabular-nums">
              {item.count}
            </span>
          )}
        </>
      )}
    </Link>
  );

  if (!collapsed) {
    return content;
  }

  return (
    <Tooltip delayDuration={100}>
      <TooltipTrigger asChild>{content}</TooltipTrigger>
      <TooltipContent side="right" className="flex items-center gap-2">
        <span>{item.title}</span>
        {showCount && <span className="text-muted-foreground text-xs">({item.count})</span>}
      </TooltipContent>
    </Tooltip>
  );
}
