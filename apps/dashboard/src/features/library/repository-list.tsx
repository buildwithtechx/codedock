import {
  ArrowRight,
  BookOpen,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  Globe,
  Lock,
  Search,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { cn } from '#/lib/utils';
import { useLibraryRepos } from './hooks';
import { ProviderAccounts } from './provider-accounts';
import type { LibraryRepo, ProviderConnection, SortBy, VisibilityFilter } from './types';

const DOCS_URL = 'https://docs.codedock.run';

function timeAgo(dateStr?: string): string | null {
  if (!dateStr) return null;
  const time = Date.parse(dateStr);
  if (!Number.isFinite(time)) return null;
  const mins = Math.floor((Date.now() - time) / 60000);
  if (mins < 1) return 'just now';
  if (mins < 60) return `${mins}m ago`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  return `${Math.floor(days / 30)}mo ago`;
}

export function RepositoryList({
  connections,
  provider,
  onProviderChange,
  onImport,
  onImportUrl,
}: {
  connections: ProviderConnection[];
  provider: string;
  onProviderChange: (provider: string) => void;
  onImport: (repo: LibraryRepo) => void;
  onImportUrl: () => void;
}) {
  const repos = useLibraryRepos(provider, true);

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="border-border/50 border-b px-5 py-4">
        <div className="mb-4">
          <ProviderAccounts
            connections={connections}
            selected={provider}
            onSelect={onProviderChange}
            onImportUrl={onImportUrl}
          />
        </div>
        <div className="flex items-center gap-2">
          <div className="relative flex-1">
            <Search
              aria-hidden="true"
              className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground"
            />
            <Input
              type="text"
              value={repos.search}
              onChange={(event) => repos.setSearch(event.target.value)}
              placeholder="Search repositories…"
              aria-label="Search repositories"
              className="pr-4 pl-10"
            />
          </div>
          <div className="hidden items-center rounded-xl border border-border/50 bg-muted/40 p-0.5 sm:flex">
            {(['all', 'public', 'private'] as VisibilityFilter[]).map((value) => (
              <button
                key={value}
                type="button"
                onClick={() => repos.setVisibility(value)}
                className={cn(
                  'rounded-lg px-3.5 py-2 font-medium text-xs capitalize transition-all',
                  repos.visibility === value
                    ? 'bg-foreground text-background shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'
                )}
              >
                {value}
              </button>
            ))}
          </div>
          <div className="relative">
            <select
              value={repos.sort}
              onChange={(event) => repos.setSort(event.target.value as SortBy)}
              aria-label="Sort repositories"
              className="cursor-pointer appearance-none rounded-xl border border-border/50 bg-muted/40 py-2.5 pr-8 pl-3.5 font-medium text-foreground text-xs transition-all hover:bg-muted/70 focus:outline-none focus:ring-2 focus:ring-primary/20"
            >
              <option value="updated">Recent</option>
              <option value="name">Name</option>
            </select>
            <ChevronDown className="pointer-events-none absolute top-1/2 right-2.5 size-3.5 -translate-y-1/2 text-muted-foreground" />
          </div>
        </div>
      </div>

      {repos.loading ? (
        <div className="divide-y divide-border/50">
          {Array.from({ length: 5 }, (_, index) => (
            <div key={index} className="flex animate-pulse items-center gap-4 px-5 py-4">
              <div className="h-10 w-10 rounded-xl bg-muted" />
              <div className="flex-1 space-y-2">
                <div className="h-4 w-32 rounded bg-muted" />
                <div className="h-3 w-48 rounded bg-muted" />
              </div>
            </div>
          ))}
        </div>
      ) : repos.isError ? (
        <div className="px-6 py-12 text-center">
          <p className="font-medium text-foreground">Repositories could not be loaded</p>
          <p className="mx-auto mt-1 max-w-sm text-muted-foreground text-sm">
            The provider may be unreachable or the token may have expired.
          </p>
          <button
            type="button"
            onClick={() => void repos.refetch()}
            className="mt-4 font-medium text-primary text-sm hover:underline"
          >
            Try again
          </button>
        </div>
      ) : repos.repos.length === 0 ? (
        <div className="px-6 py-12 text-center">
          <h3 className="font-medium text-foreground/80 text-lg">
            {repos.search ? 'No matching repositories' : 'No repositories found'}
          </h3>
          <p className="mx-auto mt-1 max-w-sm text-muted-foreground text-sm leading-relaxed">
            {repos.search
              ? 'Try a different search term or visibility filter.'
              : 'This account has no repositories visible to Codedock yet.'}
          </p>
          {!repos.search && (
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
        <>
          <div className="divide-y divide-border/50">
            {repos.repos.map((repo) => {
              const updated = timeAgo(repo.updatedAt);
              return (
                <button
                  key={repo.id}
                  type="button"
                  onClick={() => onImport(repo)}
                  className="group flex w-full items-center justify-between gap-4 px-5 py-3.5 text-left transition-colors hover:bg-muted/40"
                >
                  <div className="flex min-w-0 items-center gap-4">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-muted/60 transition-colors group-hover:bg-muted">
                      {repo.private ? (
                        <Lock className="size-[18px] text-muted-foreground" />
                      ) : (
                        <Globe className="size-[18px] text-muted-foreground" />
                      )}
                    </div>
                    <div className="min-w-0">
                      <div className="flex items-center gap-2">
                        <p className="truncate font-medium text-foreground text-sm">{repo.name}</p>
                        {repo.private && (
                          <span className="rounded-md bg-muted px-1.5 py-0.5 font-medium text-[10px] text-muted-foreground">
                            Private
                          </span>
                        )}
                      </div>
                      <p className="mt-0.5 truncate font-mono text-muted-foreground text-xs">
                        {repo.fullName}
                      </p>
                    </div>
                  </div>
                  <div className="flex shrink-0 items-center gap-3">
                    {updated && <span className="text-muted-foreground text-xs">{updated}</span>}
                    <ArrowRight className="size-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-muted-foreground" />
                  </div>
                </button>
              );
            })}
          </div>
          <div className="flex items-center justify-between gap-3 border-border/30 border-t px-5 py-3">
            <span className="text-muted-foreground/50 text-xs">
              {repos.total} repositor{repos.total === 1 ? 'y' : 'ies'}
              {repos.search && ` matching "${repos.search}"`}
            </span>
            {repos.totalPages > 1 && (
              <div className="flex items-center gap-1.5">
                <button
                  type="button"
                  onClick={() => repos.setPage(repos.page - 1)}
                  disabled={repos.page <= 1}
                  aria-label="Previous page"
                  className="inline-flex size-7 items-center justify-center rounded-lg border border-border/50 text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
                >
                  <ChevronLeft className="size-4" />
                </button>
                <span className="px-1 text-muted-foreground text-xs tabular-nums">
                  {repos.page} / {repos.totalPages}
                </span>
                <button
                  type="button"
                  onClick={() => repos.setPage(repos.page + 1)}
                  disabled={repos.page >= repos.totalPages}
                  aria-label="Next page"
                  className="inline-flex size-7 items-center justify-center rounded-lg border border-border/50 text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
                >
                  <ChevronRight className="size-4" />
                </button>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
