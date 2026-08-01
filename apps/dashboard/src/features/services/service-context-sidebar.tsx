import { Link, useLocation } from '@tanstack/react-router';
import { ChevronRight, type LucideIcon } from 'lucide-react';

export type ServiceContextTab = {
  name: string;
  href: string;
  icon: LucideIcon;
  group: string;
  exact?: boolean;
};

export function ServiceContextSidebar({
  name,
  type,
  status,
  tabs,
}: {
  name: string;
  type: string;
  status?: string;
  tabs: ServiceContextTab[];
}) {
  const location = useLocation();
  const groups = tabs.reduce<Record<string, ServiceContextTab[]>>((result, tab) => {
    result[tab.group] ??= [];
    result[tab.group].push(tab);
    return result;
  }, {});

  return (
    <div className="space-y-3">
      <section className="rounded-2xl border border-border/80 bg-card p-4 shadow-sm">
        <div className="min-w-0">
          <p className="text-muted-foreground text-xs">{type}</p>
          <h2 className="mt-0.5 truncate font-semibold text-sm">{name}</h2>
        </div>
        <div className="mt-4 flex items-center gap-2 border-border/70 border-t pt-3 text-xs">
          <span
            className={`h-1.5 w-1.5 rounded-full ${
              status === 'running'
                ? 'bg-emerald-500'
                : status === 'error'
                  ? 'bg-rose-500'
                  : 'bg-amber-400'
            }`}
          />
          <span className="text-muted-foreground capitalize">{status || 'Unknown state'}</span>
        </div>
      </section>

      <nav
        className="rounded-2xl border border-border/80 bg-card p-2 shadow-sm"
        aria-label="Service sections"
      >
        {Object.entries(groups).map(([group, groupTabs], index) => (
          <div key={group} className={index === 0 ? '' : 'mt-4'}>
            <p className="px-3 pb-1.5 font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.16em]">
              {group}
            </p>
            {groupTabs.map((tab) => {
              const active = tab.exact
                ? location.pathname === tab.href || location.pathname === `${tab.href}/`
                : location.pathname.startsWith(tab.href);

              return (
                <Link
                  key={tab.name}
                  to={tab.href}
                  className={`flex items-center gap-3 rounded-xl px-3 py-2.5 font-medium text-sm transition-colors ${
                    active
                      ? 'bg-primary/12 text-foreground'
                      : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'
                  }`}
                >
                  <tab.icon className={`h-4 w-4 ${active ? 'text-primary' : ''}`} />
                  <span className="flex-1">{tab.name}</span>
                  {active && <ChevronRight className="h-3.5 w-3.5 text-primary" />}
                </Link>
              );
            })}
          </div>
        ))}
      </nav>
    </div>
  );
}
