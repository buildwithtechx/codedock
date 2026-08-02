import { Link } from '@tanstack/react-router';
import { ArrowLeft, ArrowRight, Box } from 'lucide-react';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { ProjectEmptyState } from './project-empty-state';

export function AppCreationPage() {
  const { data, isLoading, isError, refetch } = useListCanvasSummaries();
  const projects = data?.data || [];

  return (
    <div className="grid min-w-0 gap-6 xl:grid-cols-[minmax(0,1fr)_21.25rem]">
      <main className="min-w-0">
        <Link
          to="/apps"
          className="inline-flex items-center gap-2 text-muted-foreground text-sm transition-colors hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          Apps
        </Link>
        <header className="mt-6">
          <p className="font-medium text-muted-foreground text-sm">New application</p>
          <h1 className="mt-1 font-semibold text-2xl tracking-tight">Choose a project</h1>
          <p className="mt-1 max-w-2xl text-muted-foreground text-sm">
            Applications are deployed inside a project so they share environments, domains, and
            access rules.
          </p>
        </header>

        {isLoading ? (
          <div className="mt-8 grid gap-3 md:grid-cols-2">
            {[0, 1, 2, 3].map((index) => (
              <div
                key={index}
                className="h-30 animate-pulse rounded-2xl border border-border bg-card"
              />
            ))}
          </div>
        ) : isError ? (
          <QueryErrorState
            className="mt-8"
            title="Projects are unavailable"
            description="Codedock could not load projects for the active workspace."
            onRetry={() => void refetch()}
          />
        ) : projects.length === 0 ? (
          <ProjectEmptyState />
        ) : (
          <section className="mt-8 grid gap-3 md:grid-cols-2">
            {projects.map((project) => (
              <Link
                key={project.id}
                to="/projects/$projectId/new"
                params={{ projectId: project.id }}
                className="group rounded-2xl bg-card p-5 transition-colors hover:bg-muted/60"
              >
                <div className="flex items-start gap-3">
                  <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-muted text-muted-foreground transition-colors group-hover:bg-primary/12 group-hover:text-primary">
                    <Box className="h-4 w-4" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <h2 className="truncate font-semibold text-sm">{project.name}</h2>
                    <p className="mt-1 line-clamp-2 text-muted-foreground text-xs leading-5">
                      {project.description || 'No description yet'}
                    </p>
                  </div>
                  <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground transition-colors group-hover:text-primary" />
                </div>
                <div className="mt-4 flex items-center gap-4 border-border/70 border-t pt-3 text-muted-foreground text-xs">
                  <span>{project.totalServices} services</span>
                  <span>{project.defaultEnvironment?.name || 'No environment'}</span>
                </div>
              </Link>
            ))}
          </section>
        )}
      </main>

      <aside className="hidden xl:sticky xl:top-6 xl:block xl:self-start">
        <section className="rounded-2xl bg-card p-5">
          <div className="flex items-center gap-2">
            <Box className="h-4 w-4 text-primary" />
            <h2 className="font-semibold text-sm">Application setup</h2>
          </div>
          <p className="mt-3 text-muted-foreground text-sm leading-6">
            The next screen lets you deploy from Git, Docker, a database template, or Compose.
          </p>
        </section>
      </aside>
    </div>
  );
}
