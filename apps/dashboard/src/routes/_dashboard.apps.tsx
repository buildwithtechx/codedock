import { createFileRoute, Link } from '@tanstack/react-router';
import { ArrowRight, BookOpen, ExternalLink, Plus, Search } from 'lucide-react';
import { useMemo, useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Skeleton } from '#/components/ui/skeleton';
import {
  AppLogo,
  CatalogShortcut,
  InstalledAppRow,
  useAppCatalog,
  useInstalledApps,
} from '#/features/apps';

export const Route = createFileRoute('/_dashboard/apps')({
  component: AppsPage,
});

const DOCS_URL = 'https://docs.codedock.run';

function AppsPage() {
  const installedQuery = useInstalledApps();
  const catalogQuery = useAppCatalog();
  const [search, setSearch] = useState('');
  const apps = useMemo(() => installedQuery.data ?? [], [installedQuery.data]);
  const catalog = useMemo(
    () => (Array.isArray(catalogQuery.data) ? catalogQuery.data : []),
    [catalogQuery.data]
  );

  const installedIds = useMemo(() => new Set(apps.map((app) => app.appId)), [apps]);
  const suggestions = useMemo(
    () => catalog.filter((app) => !installedIds.has(app.id)).slice(0, 5),
    [catalog, installedIds]
  );
  const query = search.trim().toLowerCase();
  const filtered = apps.filter((app) =>
    [app.name, app.branch, app.domain].some((value) => value?.toLowerCase().includes(query))
  );

  if (installedQuery.isError) {
    return (
      <QueryErrorState
        title="Apps are unavailable"
        description="Codedock could not load your installed apps."
        onRetry={() => void installedQuery.refetch()}
      />
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Apps"
        description={
          installedQuery.isLoading
            ? 'Loading…'
            : `${apps.length} installed app${apps.length === 1 ? '' : 's'}`
        }
        action={
          apps.length > 0 ? (
            <Button asChild>
              <Link to="/apps/new">
                <Plus className="size-4" />
                New app
              </Link>
            </Button>
          ) : undefined
        }
      />

      {installedQuery.isLoading ? (
        <div className="divide-y divide-border/50 rounded-2xl bg-card" aria-busy="true">
          {Array.from({ length: 3 }, (_, index) => (
            <div key={index} className="flex items-center gap-4 px-5 py-4">
              <Skeleton className="size-10 shrink-0 rounded-xl" />
              <div className="min-w-0 flex-1 space-y-2">
                <Skeleton className="h-4 w-32 max-w-full" />
                <Skeleton className="h-3 w-48 max-w-full" />
              </div>
            </div>
          ))}
        </div>
      ) : apps.length === 0 ? (
        <div className="py-8 sm:py-12">
          <div className="flex items-center justify-center" aria-hidden="true">
            {['A', 'B'].map((letter, index) => (
              <div key={letter} className="flex items-center">
                {index > 0 && (
                  <span className="mx-1.5 w-4 border-border border-t-2 border-dashed sm:w-8" />
                )}
                <div className="flex size-14 shrink-0 items-center justify-center rounded-2xl bg-card font-semibold text-lg text-primary">
                  {letter}
                </div>
              </div>
            ))}
            <span className="mx-1.5 w-4 border-border border-t-2 border-dashed sm:w-8" />
            <div className="flex size-14 shrink-0 items-center justify-center rounded-2xl border border-primary/40 border-dashed bg-primary/5">
              <Plus className="size-6 text-primary" />
            </div>
          </div>
          <div className="mt-8 text-center">
            <h2 className="font-medium text-2xl tracking-tight">No apps installed yet</h2>
            <p className="mx-auto mt-2 max-w-md text-muted-foreground text-sm leading-relaxed">
              Install a one-click app from the catalog, or import your own code from the library.
            </p>
            <div className="mt-6 flex flex-wrap items-center justify-center gap-3">
              <Button asChild>
                <Link to="/apps/new">
                  <Plus className="size-4" />
                  New app
                </Link>
              </Button>
              <Button asChild variant="outline">
                <Link to="/library">Open library</Link>
              </Button>
              <Button asChild variant="secondary">
                <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
                  <BookOpen className="size-4" />
                  Docs
                  <ExternalLink className="size-3.5 opacity-60" />
                </a>
              </Button>
            </div>
          </div>

          {(catalogQuery.isLoading || suggestions.length > 0) && (
            <section className="mx-auto mt-10 max-w-2xl" aria-labelledby="popular-apps">
              <h2 id="popular-apps" className="mb-4 font-medium text-muted-foreground text-sm">
                Popular apps
              </h2>
              <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                {catalogQuery.isLoading
                  ? Array.from({ length: 4 }, (_, index) => (
                      <div
                        key={index}
                        className="flex animate-pulse items-center gap-3 rounded-xl bg-card p-4"
                        aria-hidden="true"
                      >
                        <div className="size-10 shrink-0 rounded-lg bg-muted" />
                        <div className="min-w-0 flex-1 space-y-2">
                          <div className="h-4 w-24 max-w-full rounded bg-muted" />
                          <div className="h-3 w-40 max-w-full rounded bg-muted/60" />
                        </div>
                      </div>
                    ))
                  : suggestions.map((app) => <CatalogShortcut key={app.id} app={app} />)}
                <Link
                  to="/apps/new"
                  className="group flex items-center gap-3 rounded-xl bg-muted/40 p-4 transition-colors hover:bg-muted/60"
                >
                  <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted/60">
                    <Plus className="size-5 text-muted-foreground" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="font-medium text-foreground text-sm">Browse all apps</p>
                    {!catalogQuery.isLoading && catalog.length - suggestions.length > 0 && (
                      <p className="mt-0.5 text-muted-foreground text-xs">
                        {catalog.length - suggestions.length} more in the catalog
                      </p>
                    )}
                  </div>
                  <ArrowRight className="size-4 shrink-0 text-muted-foreground" />
                </Link>
              </div>
            </section>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
          <section className="min-w-0 space-y-4" aria-label="Installed apps">
            <div className="relative">
              <Search className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                type="search"
                value={search}
                onChange={(event) => setSearch(event.target.value)}
                placeholder="Search installed apps…"
                aria-label="Search installed apps"
                className="h-10 bg-muted/60 pr-4 pl-10"
              />
            </div>
            <div className="divide-y divide-border/50 rounded-2xl bg-card">
              {filtered.length > 0 ? (
                filtered.map((app) => <InstalledAppRow key={app.id} app={app} />)
              ) : (
                <p className="px-5 py-12 text-center text-muted-foreground text-sm">
                  No apps match your search.
                </p>
              )}
            </div>
          </section>

          <aside
            className="rounded-2xl bg-card p-5 lg:sticky lg:top-6"
            aria-labelledby="suggested-apps"
          >
            <h2 id="suggested-apps" className="mb-4 font-medium text-foreground text-sm">
              You can also deploy
            </h2>
            {catalogQuery.isLoading ? (
              <div className="flex items-center gap-2" aria-busy="true">
                {[0, 1, 2].map((index) => (
                  <Skeleton key={index} className="size-10 rounded-xl" />
                ))}
              </div>
            ) : suggestions.length > 0 ? (
              <div className="flex items-center gap-2">
                {suggestions.slice(0, 3).map((app) => (
                  <AppLogo key={app.id} icon={app.icon} name={app.name} className="size-10" />
                ))}
              </div>
            ) : (
              <p className="text-muted-foreground text-sm">
                Every catalog app is already installed.
              </p>
            )}
            <Button asChild variant="secondary" className="mt-4 w-full">
              <Link to="/apps/new">
                Browse all apps
                <ArrowRight className="size-4" />
              </Link>
            </Button>
          </aside>
        </div>
      )}
    </div>
  );
}
