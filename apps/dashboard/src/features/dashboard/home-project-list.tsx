import { Link } from '@tanstack/react-router';
import { ArrowUpRight, FolderKanban, Plus } from 'lucide-react';
import { Button } from '#/components/ui/button';
import type { CanvasSummary } from '#/features/projects';

export function HomeProjectList({
  projects,
  isLoading,
  onCreateProject,
}: {
  projects: CanvasSummary[];
  isLoading: boolean;
  onCreateProject: () => void;
}) {
  return (
    <section className="overflow-hidden rounded-2xl border border-border/80 bg-card shadow-sm">
      <div className="flex items-center justify-between border-border/70 border-b px-5 py-4">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/12 text-primary">
            <FolderKanban className="h-4 w-4" />
          </div>
          <div>
            <h2 className="font-semibold text-sm">Active projects</h2>
            <p className="text-muted-foreground text-xs">
              {isLoading ? 'Loading workspace...' : `${projects.length} in this organization`}
            </p>
          </div>
        </div>
        <Link
          to="/projects"
          className="flex items-center gap-1 font-medium text-muted-foreground text-xs transition-colors hover:text-foreground"
        >
          View all
          <ArrowUpRight className="h-3.5 w-3.5" />
        </Link>
      </div>

      {isLoading ? (
        <div className="space-y-px bg-border/50">
          {[0, 1, 2].map((index) => (
            <div key={index} className="flex items-center gap-3 bg-card px-5 py-4">
              <div className="h-9 w-9 animate-pulse rounded-xl bg-muted" />
              <div className="flex-1 space-y-2">
                <div className="h-3 w-32 animate-pulse rounded bg-muted" />
                <div className="h-2.5 w-48 animate-pulse rounded bg-muted" />
              </div>
            </div>
          ))}
        </div>
      ) : projects.length === 0 ? (
        <div className="flex min-h-[25rem] flex-col items-center justify-center px-6 py-10 text-center">
          <div className="grid grid-cols-[auto_2rem_auto_2rem_auto] items-center text-primary">
            <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-primary/12">
              <FolderKanban className="h-4 w-4" />
            </div>
            <div className="h-px bg-primary/35" />
            <div className="flex h-11 w-11 items-center justify-center rounded-xl border border-primary/30 bg-card">
              <span className="h-3 w-3 rounded-sm border-2 border-primary" />
            </div>
            <div className="h-px bg-primary/35" />
            <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-primary text-primary-foreground">
              <Plus className="h-4 w-4" />
            </div>
          </div>
          <h3 className="mt-5 font-semibold text-base">Start building your workspace</h3>
          <p className="mt-1 max-w-xs text-muted-foreground text-sm">
            Projects keep environments, applications, domains, and deployments together.
          </p>
          <Button className="mt-5 gap-2" onClick={onCreateProject}>
            <Plus className="h-4 w-4" />
            Create project
          </Button>
        </div>
      ) : (
        <div className="divide-y divide-border/70">
          {projects.slice(0, 5).map((project) => (
            <Link
              key={project.id}
              to="/projects/$projectId"
              params={{ projectId: project.id }}
              className="group flex items-center gap-3 px-5 py-3.5 transition-colors hover:bg-muted/35"
            >
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-muted text-muted-foreground transition-colors group-hover:bg-primary/12 group-hover:text-primary">
                <FolderKanban className="h-4 w-4" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-sm">{project.name}</p>
                <p className="mt-0.5 truncate text-muted-foreground text-xs">
                  {project.totalServices} service{project.totalServices === 1 ? '' : 's'} ·{' '}
                  {project.defaultEnvironment?.name ?? 'No environment'}
                </p>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <span
                  className={`h-1.5 w-1.5 rounded-full ${
                    project.onlineServices > 0 ? 'bg-emerald-400' : 'bg-muted-foreground/45'
                  }`}
                />
                <span className="text-muted-foreground">
                  {project.onlineServices}/{project.totalServices}
                </span>
                <ArrowUpRight className="h-3.5 w-3.5 text-muted-foreground transition-colors group-hover:text-foreground" />
              </div>
            </Link>
          ))}
        </div>
      )}
    </section>
  );
}
