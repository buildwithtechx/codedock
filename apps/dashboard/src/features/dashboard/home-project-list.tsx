import { Link } from '@tanstack/react-router';
import { ArrowUpRight, FolderKanban } from 'lucide-react';
import type { CanvasSummary } from '#/features/projects';
import { HomeFirstProject } from './home-first-project';

export function HomeProjectList({
  projects,
  isLoading,
  onCreateProject,
}: {
  projects: CanvasSummary[];
  isLoading: boolean;
  onCreateProject: () => void;
}) {
  if (!isLoading && projects.length === 0) {
    return <HomeFirstProject onCreateProject={onCreateProject} />;
  }

  return (
    <section className="overflow-hidden rounded-2xl bg-card">
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
