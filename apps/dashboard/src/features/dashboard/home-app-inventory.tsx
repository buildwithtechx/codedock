import { Link } from '@tanstack/react-router';
import { ArrowRight, Box, Database, Plus, Workflow } from 'lucide-react';
import type { CanvasSummary } from '#/features/projects';

export function HomeAppInventory({
  projects,
  isLoading,
}: {
  projects: CanvasSummary[];
  isLoading: boolean;
}) {
  const services = projects.reduce((total, project) => total + project.totalServices, 0);

  return (
    <section className="rounded-2xl bg-card p-5">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Box className="h-4 w-4 text-primary" />
          <h2 className="font-semibold text-sm">Applications</h2>
        </div>
        <Link
          to="/apps/new"
          className="flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          aria-label="Deploy an application"
        >
          <Plus className="h-4 w-4" />
        </Link>
      </div>
      {isLoading ? (
        <div className="mt-5 grid grid-cols-3 gap-2">
          {[0, 1, 2].map((item) => (
            <span key={item} className="h-15 animate-pulse rounded-xl bg-muted/70" />
          ))}
        </div>
      ) : services === 0 ? (
        <div className="pt-5 text-center">
          <div className="flex items-center justify-center gap-2" aria-hidden="true">
            <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <Box className="h-4 w-4" />
            </span>
            <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <Workflow className="h-4 w-4" />
            </span>
            <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <Database className="h-4 w-4" />
            </span>
          </div>
          <p className="mt-4 font-medium text-sm">Nothing deployed yet</p>
          <p className="mt-1 text-muted-foreground/75 text-xs leading-5">
            Services and data stores will appear here.
          </p>
        </div>
      ) : (
        <div className="mt-5 divide-y divide-border/65">
          {projects.slice(0, 3).map((project) => (
            <Link
              key={project.id}
              to="/projects/$projectId"
              params={{ projectId: project.id }}
              className="flex items-center gap-3 py-2.5 first:pt-0 last:pb-0"
            >
              <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                <Box className="h-4 w-4" />
              </span>
              <span className="min-w-0 flex-1 truncate font-medium text-sm">{project.name}</span>
              <span className="text-muted-foreground text-xs">{project.totalServices}</span>
              <ArrowRight className="h-3.5 w-3.5 text-muted-foreground" />
            </Link>
          ))}
        </div>
      )}
    </section>
  );
}
