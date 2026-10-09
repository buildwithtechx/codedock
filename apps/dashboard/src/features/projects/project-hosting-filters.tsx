import type { LucideIcon } from 'lucide-react';
import { HardDrive, LayoutGrid, Server } from 'lucide-react';
import type { CanvasSummary } from '#/features/projects';

export type ProjectHostingFilter =
  | { kind: 'all' }
  | { kind: 'local' }
  | { kind: 'server'; name: string };

export function projectFilterKey(filter: ProjectHostingFilter): string {
  return filter.kind === 'server' ? `server:${filter.name}` : filter.kind;
}

export function projectMatchesHostingFilter(
  project: CanvasSummary,
  filter: ProjectHostingFilter
): boolean {
  switch (filter.kind) {
    case 'all':
      return true;
    case 'local':
      return project.deployTarget === 'local';
    case 'server':
      return project.deployTarget === 'server' && (project.serverName || 'Server') === filter.name;
  }
}

export interface ProjectHostingOption {
  key: string;
  filter: ProjectHostingFilter;
  label: string;
  icon: LucideIcon;
  count: number;
}

export function buildHostingFilterOptions(projects: CanvasSummary[]): ProjectHostingOption[] {
  let local = 0;
  const servers = new Map<string, number>();
  for (const project of projects) {
    if (project.deployTarget === 'server') {
      const name = project.serverName || 'Server';
      servers.set(name, (servers.get(name) ?? 0) + 1);
    } else {
      local++;
    }
  }
  const options: ProjectHostingOption[] = [
    {
      key: 'all',
      filter: { kind: 'all' },
      label: 'All projects',
      icon: LayoutGrid,
      count: projects.length,
    },
  ];
  for (const [name, count] of [...servers.entries()].sort((a, b) => a[0].localeCompare(b[0]))) {
    options.push({
      key: `server:${name}`,
      filter: { kind: 'server', name },
      label: name,
      icon: Server,
      count,
    });
  }
  if (local > 0) {
    options.push({
      key: 'local',
      filter: { kind: 'local' },
      label: 'Local',
      icon: HardDrive,
      count: local,
    });
  }
  return options;
}

export function ProjectHostingFilters({
  options,
  active,
  onChange,
}: {
  options: ProjectHostingOption[];
  active: ProjectHostingFilter;
  onChange: (filter: ProjectHostingFilter) => void;
}) {
  const activeKey = projectFilterKey(active);
  return (
    <div className="rounded-2xl bg-card">
      <div className="border-border/50 border-b px-5 py-4">
        <h2 className="font-semibold text-[15px] text-foreground">Filters</h2>
        <p className="text-muted-foreground text-xs">Filter by deploy target</p>
      </div>
      <div className="p-2">
        {options.map((option) => {
          const isActive = option.key === activeKey;
          return (
            <button
              key={option.key}
              type="button"
              onClick={() => onChange(option.filter)}
              aria-pressed={isActive}
              className={`flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-sm transition-colors ${
                isActive
                  ? 'bg-primary/10 text-foreground'
                  : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'
              }`}
            >
              <span className={isActive ? 'text-primary' : 'text-muted-foreground'}>
                <option.icon className="h-4 w-4" />
              </span>
              <span className="flex-1 truncate text-start">{option.label}</span>
              <span
                className={`text-xs tabular-nums ${isActive ? 'text-foreground' : 'text-muted-foreground/60'}`}
              >
                {option.count}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
