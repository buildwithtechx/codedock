import { Link } from '@tanstack/react-router';
import { ArrowRight, Box, Database, Plus, Workflow } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import type { CanvasSummary } from '#/interfaces/projects';

export function HomeAppInventory({
  projects,
  isLoading,
}: {
  projects: CanvasSummary[];
  isLoading: boolean;
}) {
  const services = projects.reduce((total, project) => total + project.totalServices, 0);

  return (
    <section className="rounded-2xl border border-border/50 bg-card p-5 shadow-xs">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Box className="size-4 text-primary" />
          <h2 className="font-semibold text-foreground text-sm">Applications</h2>
        </div>
        <Button
          asChild
          variant="ghost"
          size="icon"
          className="size-7 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          <Link to="/library" aria-label="Deploy an application">
            <Plus className="size-4" />
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <div className="mt-4 grid grid-cols-3 gap-2">
          {[0, 1, 2].map((item) => (
            <span key={item} className="h-14 animate-pulse rounded-xl bg-muted/60" />
          ))}
        </div>
      ) : services === 0 ? (
        <div className="pt-4 text-center">
          <div className="flex items-center justify-center gap-2" aria-hidden="true">
            <span className="flex size-9 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
              <Box className="size-4" />
            </span>
            <span className="flex size-9 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
              <Workflow className="size-4" />
            </span>
            <span className="flex size-9 items-center justify-center rounded-xl bg-muted/50 text-muted-foreground">
              <Database className="size-4" />
            </span>
          </div>
          <p className="mt-3.5 font-medium text-foreground text-sm">Nothing deployed yet</p>
          <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
            Deploy services, workers, and databases to this workspace.
          </p>
          <Button asChild variant="outline" size="sm" className="mt-4 h-8 w-full text-xs font-medium">
            <Link to="/library">
              <Plus className="mr-1.5 size-3.5" />
              <span>Explore Catalog & Deploy</span>
            </Link>
          </Button>
        </div>
      ) : (
        <div className="mt-3.5 divide-y divide-border/40">
          {projects.slice(0, 4).map((project) => (
            <Link
              key={project.id}
              to="/projects/$projectId"
              params={{ projectId: project.id }}
              className="group flex items-center justify-between gap-3 py-2.5 transition-colors first:pt-1 last:pb-0"
            >
              <div className="flex min-w-0 items-center gap-2.5">
                <span className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-muted/60 text-muted-foreground group-hover:text-primary transition-colors">
                  <Box className="size-4" />
                </span>
                <div className="min-w-0 flex-1 truncate">
                  <p className="truncate font-medium text-foreground text-xs group-hover:text-primary transition-colors">
                    {project.name}
                  </p>
                  <p className="text-[11px] text-muted-foreground">
                    {project.totalServices} {project.totalServices === 1 ? 'service' : 'services'}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-2 shrink-0">
                <span
                  className={`size-1.5 rounded-full ${
                    project.onlineServices > 0 ? 'bg-emerald-500' : 'bg-muted-foreground/30'
                  }`}
                />
                <ArrowRight className="size-3.5 text-muted-foreground opacity-60 group-hover:opacity-100 group-hover:translate-x-0.5 transition-all" />
              </div>
            </Link>
          ))}
        </div>
      )}
    </section>
  );
}
