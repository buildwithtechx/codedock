import { Link } from '@tanstack/react-router';
import { ArrowUpRight, Box, CircleAlert, Loader2, Plus } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { useListAppsByOrganization } from '#/hooks/use-apps';

const statusClasses = {
  running: 'bg-emerald-500',
  building: 'bg-amber-400',
  stopped: 'bg-muted-foreground/50',
  error: 'bg-rose-500',
  created: 'bg-sky-500',
};

export function AppDirectory() {
  const { data, isLoading, isError } = useListAppsByOrganization();
  const apps = data?.data || [];

  return (
    <div className="space-y-6">
      <header className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
        <div>
          <p className="font-medium text-muted-foreground text-sm">Running workloads</p>
          <h1 className="mt-1 font-semibold text-2xl tracking-tight">Apps</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Every application deployed in the active organization.
          </p>
        </div>
        <Link to="/projects" className="self-start sm:self-auto">
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            New app
          </Button>
        </Link>
      </header>

      {isLoading ? (
        <div className="flex min-h-80 items-center justify-center rounded-2xl border border-border/80 bg-card">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
        </div>
      ) : isError ? (
        <div className="flex min-h-80 flex-col items-center justify-center rounded-2xl border border-border/80 bg-card px-6 text-center">
          <CircleAlert className="h-6 w-6 text-destructive" />
          <h2 className="mt-3 font-semibold">Apps are unavailable</h2>
          <p className="mt-1 text-muted-foreground text-sm">
            Refresh the page or check your organization access.
          </p>
        </div>
      ) : apps.length === 0 ? (
        <div className="flex min-h-80 flex-col items-center justify-center rounded-2xl border border-border border-dashed bg-card px-6 text-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
            <Box className="h-5 w-5" />
          </div>
          <h2 className="mt-4 font-semibold text-lg">No apps deployed</h2>
          <p className="mt-1 max-w-sm text-muted-foreground text-sm">
            Choose a project, then add an application to start a deployment.
          </p>
          <Link to="/projects" className="mt-5">
            <Button>Browse projects</Button>
          </Link>
        </div>
      ) : (
        <div className="grid gap-3 lg:grid-cols-2 xl:grid-cols-3">
          {apps.map((app) => (
            <Link
              key={app.id}
              to="/services/$serviceId"
              params={{ serviceId: app.id }}
              className="group rounded-xl border border-border/80 bg-card p-4 shadow-sm transition-colors hover:border-primary/35 hover:bg-primary/4"
            >
              <div className="flex items-start gap-3">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-muted text-muted-foreground transition-colors group-hover:bg-primary/12 group-hover:text-primary">
                  <Box className="h-4 w-4" />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <h2 className="truncate font-semibold text-sm">{app.name}</h2>
                    <span
                      className={`h-1.5 w-1.5 shrink-0 rounded-full ${statusClasses[app.status]}`}
                    />
                  </div>
                  <p className="mt-1 truncate text-muted-foreground text-xs">
                    {app.domain || app.repositoryUrl || 'No public endpoint'}
                  </p>
                </div>
                <ArrowUpRight className="h-3.5 w-3.5 shrink-0 text-muted-foreground transition-colors group-hover:text-primary" />
              </div>
              <div className="mt-4 flex items-center justify-between border-border/65 border-t pt-3 text-muted-foreground text-xs">
                <span className="capitalize">{app.runtimeMode}</span>
                <span className="capitalize">{app.status}</span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
