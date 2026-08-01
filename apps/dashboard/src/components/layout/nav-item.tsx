import { Link, useRouterState } from '@tanstack/react-router';
import { ExternalLink } from 'lucide-react';
import type React from 'react';

export type NavItemProps = {
  title: string;
  url: string;
  icon: React.ComponentType<{ className?: string }>;
  external?: boolean;
  badge?: string;
};

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
  const isActive = exact
    ? pathname === item.url
    : pathname.startsWith(item.url) && item.url !== '/';

  return (
    <Link
      to={item.url as never}
      className={`group relative flex items-center rounded-lg font-medium text-sm transition-colors ${
        collapsed ? 'justify-center px-0 py-2.5' : 'gap-3 px-2.5 py-2.25'
      } ${
        isActive
          ? 'bg-primary/12 text-sidebar-foreground'
          : 'text-sidebar-foreground/55 hover:bg-sidebar-accent hover:text-sidebar-foreground'
      }`}
      target={item.external ? '_blank' : undefined}
      rel={item.external ? 'noopener noreferrer' : undefined}
    >
      {!collapsed && isActive && (
        <div className="absolute top-1/2 -left-2 h-4 w-0.5 -translate-y-1/2 rounded-r-full bg-primary" />
      )}

      <item.icon
        className={`h-4 w-4 shrink-0 transition-colors ${
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
        </>
      )}
    </Link>
  );
}
