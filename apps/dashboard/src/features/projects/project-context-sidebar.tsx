import { Link, useLocation } from '@tanstack/react-router';
import {
  CalendarDays,
  ChevronRight,
  FolderKanban,
  Grid3X3,
  LayoutDashboard,
  Plus,
  Settings2,
  Workflow,
} from 'lucide-react';
import { useGetCanvasSummary } from '#/hooks/use-canvas';

export const getProjectContextNavigation = (projectId: string) => {
  const basePath = `/projects/${projectId}`;

  return [
    { title: 'Overview', to: basePath, icon: LayoutDashboard, exact: true },
    { title: 'Canvas', to: `${basePath}/canvas`, icon: Grid3X3 },
    { title: 'Add resource', to: `${basePath}/new`, icon: Plus },
    { title: 'Compose', to: `${basePath}/compose`, icon: Workflow },
    { title: 'Scheduled tasks', to: `${basePath}/scheduled-tasks`, icon: CalendarDays },
    { title: 'Project settings', to: `${basePath}/settings`, icon: Settings2 },
  ];
};

export function ProjectContextSidebar({ projectId }: { projectId: string }) {
  const location = useLocation();
  const { data } = useGetCanvasSummary(projectId);
  const project = data?.data;
  const navigation = getProjectContextNavigation(projectId);

  return (
    <div className="space-y-3">
      <section className="rounded-2xl bg-card p-4">
        <div className="flex items-start gap-3">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-primary/12 text-primary">
            <FolderKanban className="h-4 w-4" />
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-muted-foreground text-xs">Project</p>
            <h2 className="mt-0.5 truncate font-semibold text-sm">
              {project?.name || 'Loading project'}
            </h2>
          </div>
        </div>
        <div className="mt-4 space-y-2.5 border-border/70 border-t pt-3 text-xs">
          <div className="flex items-center justify-between gap-3">
            <span className="text-muted-foreground">Environment</span>
            <span className="truncate font-medium">
              {project?.defaultEnvironment?.name || 'Not configured'}
            </span>
          </div>
          <div className="flex items-center justify-between gap-3">
            <span className="text-muted-foreground">Services</span>
            <span className="font-medium">{project?.totalServices ?? 0}</span>
          </div>
        </div>
      </section>

      <nav className="rounded-2xl bg-card p-2" aria-label="Project sections">
        {navigation.map((item) => {
          const active = item.exact
            ? location.pathname === item.to || location.pathname === `${item.to}/`
            : location.pathname.startsWith(item.to);

          return (
            <Link
              key={item.title}
              to={item.to}
              className={`flex items-center gap-3 rounded-xl px-3 py-2.5 font-medium text-sm transition-colors ${
                active
                  ? 'bg-primary/12 text-foreground'
                  : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'
              }`}
            >
              <item.icon className={`h-4 w-4 ${active ? 'text-primary' : ''}`} />
              <span className="flex-1">{item.title}</span>
              {active && <ChevronRight className="h-3.5 w-3.5 text-primary" />}
            </Link>
          );
        })}
      </nav>
    </div>
  );
}
