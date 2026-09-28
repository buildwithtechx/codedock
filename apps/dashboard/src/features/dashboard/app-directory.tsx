import { Link } from '@tanstack/react-router';
import { ArrowUpRight, Box, Loader2, Plus } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useListAppsByOrganization } from '#/hooks/use-apps';
import { AppEmptyState } from './app-empty-state';

const statusClasses = {
  running: 'bg-emerald-500',
  building: 'bg-amber-400',
  stopped: 'bg-muted-foreground/50',
  error: 'bg-rose-500',
  created: 'bg-sky-500',
};

export function AppDirectory() {
  const { data, isLoading, isError, refetch } = useListAppsByOrganization();
  const apps = data?.data || [];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Apps"
        description={
          isLoading
            ? 'Loading applications...'
            : `${apps.length} application${apps.length === 1 ? '' : 's'} deployed`
        }
        action={
          apps.length > 0 ? (
            <Link to="/apps/new">
              <Button className="gap-2">
                <Plus className="h-4 w-4" />
                New app
              </Button>
            </Link>
          ) : undefined
        }
      />

      {isLoading ? (
        <div className="flex min-h-[25rem] items-center justify-center">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
        </div>
      ) : isError ? (
        <QueryErrorState
          title="Apps are unavailable"
          description="Codedock could not load apps for the active workspace."
          onRetry={() => void refetch()}
        />
      ) : apps.length === 0 ? (
        <AppEmptyState />
      ) : (
        <div className="grid gap-3 lg:grid-cols-2 xl:grid-cols-3">
          {apps.map((app) => (
            <Link
              key={app.id}
              to="/services/$serviceId"
              params={{ serviceId: app.id }}
              className="group rounded-xl bg-card p-4 transition-colors hover:bg-muted/60"
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
