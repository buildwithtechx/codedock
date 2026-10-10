import { useNavigate } from '@tanstack/react-router';
import { BookOpen, ExternalLink, Plus, Search } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Skeleton } from '#/components/ui/skeleton';
import { AddCustomAppModal } from './add-custom-app-modal';
import { CatalogCard } from './catalog-card';
import { useAppCatalog } from './hooks';

const CATEGORY_ORDER = [
  'backend',
  'database',
  'cms',
  'analytics',
  'automation',
  'mail',
  'other',
] as const;

const STANDARD_CATEGORIES = new Set<string>(CATEGORY_ORDER);

function normalizeCategory(cat?: string): string {
  if (!cat) return 'other';
  const lower = cat.toLowerCase();
  return STANDARD_CATEGORIES.has(lower) && lower !== 'other' ? lower : 'other';
}

const CATEGORY_LABELS: Record<string, string> = {
  all: 'All',
  backend: 'Backend',
  database: 'Database',
  cms: 'CMS',
  analytics: 'Analytics',
  automation: 'Automation',
  mail: 'Mail',
  other: 'Other',
};

const DOCS_URL = 'https://docs.codedock.run';

export function AppCatalog({
  embedded = false,
  deepLinkAppId,
}: {
  embedded?: boolean;
  deepLinkAppId?: string;
}) {
  const navigate = useNavigate();
  const { data, isLoading, isError, refetch } = useAppCatalog();
  const [query, setQuery] = useState('');
  const [category, setCategory] = useState('all');
  const [addCustomOpen, setAddCustomOpen] = useState(false);
  const autoStarted = useRef(false);
  const catalog = Array.isArray(data) ? data : [];

  const categories = useMemo(() => {
    const present = new Set(catalog.map((app) => normalizeCategory(app.category)));
    return ['all', ...CATEGORY_ORDER.filter((item) => present.has(item))];
  }, [catalog]);

  const filtered = useMemo(() => {
    const needle = (query ?? '').trim().toLowerCase();
    return catalog.filter((app) => {
      const appCat = normalizeCategory(app?.category);
      if (category !== 'all' && appCat !== category) return false;
      if (!needle) return true;
      return `${app?.name ?? ''} ${app?.description ?? ''} ${appCat}`
        .toLowerCase()
        .includes(needle);
    });
  }, [catalog, category, query]);

  useEffect(() => {
    if (embedded || autoStarted.current || isLoading || !deepLinkAppId) return;
    const entry = catalog.find((app) => app.id === deepLinkAppId);
    if (entry) {
      autoStarted.current = true;
      void navigate({ to: '/apps/new/$appId', params: { appId: entry.id } });
    }
  }, [embedded, isLoading, catalog, deepLinkAppId, navigate]);

  if (isError) {
    return (
      <QueryErrorState
        title="Catalog is unavailable"
        description="Codedock could not load the one-click application catalog."
        onRetry={() => void refetch()}
      />
    );
  }

  return (
    <div>
      {!embedded && (
        <PageHeader
          title="Explore & Deploy Apps"
          description="Deploy open-source templates, databases, and services in one click."
          action={
            <Button variant="outline" onClick={() => setAddCustomOpen(true)} className="gap-2">
              <Plus className="size-4" />
              Add custom app
            </Button>
          }
        />
      )}
      <div className={embedded ? 'space-y-4' : 'mt-6 space-y-4'}>
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="relative w-full max-w-md flex-1">
            <Search className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="search"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Search apps…"
              aria-label="Search apps"
              className="pr-4 pl-10"
            />
          </div>
          {embedded && (
            <Button variant="outline" onClick={() => setAddCustomOpen(true)} className="gap-2">
              <Plus className="size-4" />
              Add custom app
            </Button>
          )}
        </div>
        {categories.length > 1 && (
          <div className="flex flex-wrap items-center gap-1">
            {categories.map((item) => {
              const active = category === item;
              return (
                <button
                  key={item}
                  type="button"
                  aria-pressed={active}
                  onClick={() => setCategory(item)}
                  className={
                    active
                      ? 'rounded-lg bg-foreground px-4 py-2 font-medium text-background text-sm transition-colors'
                      : 'rounded-lg px-4 py-2 font-medium text-muted-foreground text-sm transition-colors hover:bg-muted/50 hover:text-foreground'
                  }
                >
                  {CATEGORY_LABELS[item] ?? item.charAt(0).toUpperCase() + item.slice(1)}
                </button>
              );
            })}
          </div>
        )}
      </div>

      {isLoading ? (
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2 2xl:grid-cols-3">
          {Array.from({ length: 6 }, (_, index) => (
            <Skeleton key={index} className="h-28 rounded-2xl" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="mt-6 rounded-2xl border border-border/50 bg-card px-5 py-12 text-center text-muted-foreground text-sm">
          <p>
            {catalog.length === 0
              ? 'No one-click apps are available on this instance yet.'
              : 'No apps match your search.'}
          </p>
          {catalog.length === 0 && (
            <Button asChild size="sm" variant="secondary" className="mt-4">
              <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
                <BookOpen className="size-3.5" />
                Docs
                <ExternalLink className="size-3 opacity-60" />
              </a>
            </Button>
          )}
        </div>
      ) : (
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2 2xl:grid-cols-3">
          {filtered.map((app) => (
            <CatalogCard key={app.id} app={app} />
          ))}
        </div>
      )}

      {!embedded && !isLoading && catalog.length > 0 && (
        <div className="mt-6 flex justify-end">
          <Button variant="outline" onClick={() => void refetch()}>
            Refresh catalog
          </Button>
        </div>
      )}

      <AddCustomAppModal
        open={addCustomOpen}
        onOpenChange={setAddCustomOpen}
        onAdded={() => void refetch()}
      />
    </div>
  );
}
