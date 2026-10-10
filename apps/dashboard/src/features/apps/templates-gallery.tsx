import { useNavigate } from '@tanstack/react-router';
import { ArrowRight, Search } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Skeleton } from '#/components/ui/skeleton';
import { useListExampleApps } from '#/hooks/use-templates';
import type { ExampleApp } from '#/interfaces/templates';
import { encodeDeploySlug } from '#/lib/slug-utils';
import { ExampleLogo } from './example-logo';

export function TemplatesGallery() {
  const navigate = useNavigate();
  const { data: examplesResponse, isLoading, isError, refetch } = useListExampleApps();
  const [query, setQuery] = useState('');

  const examples = Array.isArray(examplesResponse) ? examplesResponse : [];

  const filtered = useMemo(() => {
    const needle = (query ?? '').trim().toLowerCase();
    if (!needle) return examples;
    return examples.filter((example) =>
      `${example.name} ${example.description}`.toLowerCase().includes(needle)
    );
  }, [examples, query]);

  const handleSelect = (example: ExampleApp) => {
    void navigate({
      to: '/deploy/$slug',
      params: {
        slug: encodeDeploySlug('buildwithtechx/codedock-examples'),
      },
      search: {
        name: example.id,
        dir: example.id,
        template: example.id,
      },
    });
  };

  return (
    <div>
      <div className="relative w-full max-w-md">
        <Search className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          type="search"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Search templates…"
          aria-label="Search templates"
          className="pr-4 pl-10"
        />
      </div>

      {isError ? (
        <QueryErrorState
          title="Examples are unavailable"
          description="Codedock could not load the example project catalogue."
          onRetry={() => void refetch()}
        />
      ) : isLoading ? (
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2 2xl:grid-cols-3">
          {Array.from({ length: 6 }, (_, index) => (
            <Skeleton key={index} className="h-32 rounded-2xl" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="mt-6 rounded-2xl border border-border/50 bg-card px-5 py-12 text-center text-muted-foreground text-sm">
          <p>
            {examples.length === 0
              ? 'No example projects are available on this instance yet.'
              : 'No templates match your search.'}
          </p>
        </div>
      ) : (
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2 2xl:grid-cols-3">
          {filtered.map((example) => (
            <button
              key={example.id}
              type="button"
              onClick={() => handleSelect(example)}
              className="group flex w-full cursor-pointer items-start gap-3 rounded-2xl border border-border/50 bg-card p-5 text-left transition-all hover:border-primary/40 hover:shadow-md"
            >
              <ExampleLogo example={example} />
              <div className="min-w-0 flex-1">
                <div className="flex items-center justify-between gap-2">
                  <p className="truncate font-medium text-foreground">{example.name}</p>
                  <ArrowRight className="size-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-foreground" />
                </div>
                <p className="mt-1 line-clamp-2 text-muted-foreground text-sm">
                  {example.description}
                </p>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
