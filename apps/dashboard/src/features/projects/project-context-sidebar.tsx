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

const navigation = [
  { label: 'Overview', suffix: '', icon: LayoutDashboard, exact: true },
  { label: 'Canvas', suffix: '/canvas', icon: Grid3X3 },
  { label: 'Add resource', suffix: '/new', icon: Plus },
  { label: 'Compose', suffix: '/compose', icon: Workflow },
  { label: 'Scheduled tasks', suffix: '/scheduled-tasks', icon: CalendarDays },
  { label: 'Project settings', suffix: '/settings', icon: Settings2 },
];

export function ProjectContextSidebar({ projectId }: { projectId: string }) {
  const location = useLocation();
  const { data } = useGetCanvasSummary(projectId);
  const project = data?.data;
  const basePath = `/projects/${projectId}`;

  return (
    <aside className="sticky top-8 space-y-3">
      <section className="rounded-2xl border border-border/80 bg-card p-4 shadow-sm">
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

      <nav
        className="rounded-2xl border border-border/80 bg-card p-2 shadow-sm"
        aria-label="Project sections"
      >
        {navigation.map((item) => {
          const href = `${basePath}${item.suffix}`;
          const active = item.exact
            ? location.pathname === href || location.pathname === `${href}/`
            : location.pathname.startsWith(href);

          return (
            <Link
              key={item.label}
              to={href}
              className={`flex items-center gap-3 rounded-xl px-3 py-2.5 font-medium text-sm transition-colors ${
                active
                  ? 'bg-primary/12 text-foreground'
                  : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'
              }`}
            >
              <item.icon className={`h-4 w-4 ${active ? 'text-primary' : ''}`} />
              <span className="flex-1">{item.label}</span>
              {active && <ChevronRight className="h-3.5 w-3.5 text-primary" />}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
