import { Link } from '@tanstack/react-router';
import { GitBranch, Plus, SearchX } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { EmptyIllustration } from './empty-illustration';

export function DeploymentEmptyState({
  hasFilters,
  onClear,
}: {
  hasFilters: boolean;
  onClear: () => void;
}) {
  if (hasFilters) {
    return (
      <div className="rounded-2xl border border-border/50 bg-card p-16 text-center">
        <div className="mx-auto max-w-md">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full border border-border/50 bg-muted/60">
            <SearchX className="h-7 w-7 text-muted-foreground/50" />
          </div>
          <h3 className="mb-2 font-medium text-foreground/80 text-lg">
            No releases match these filters
          </h3>
          <p className="text-muted-foreground text-sm leading-relaxed">
            Change or clear a filter to look across a different set of releases.
          </p>
          <Button variant="outline" size="sm" onClick={onClear} className="mt-6">
            Clear filters
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-2xl border border-border/50 bg-card px-6 pb-10 text-center">
      <EmptyIllustration className="relative mx-auto h-44 w-64" />
      <h3 className="mb-2 font-medium text-foreground/80 text-lg">No deployments yet</h3>
      <p className="mx-auto mb-8 max-w-sm text-muted-foreground text-sm leading-relaxed">
        Deploy an app and its build status, release history, and commit details will appear here.
      </p>
      <div className="flex flex-col items-center justify-center gap-3 sm:flex-row">
        <Button asChild size="lg" className="gap-2 px-6">
          <Link to="/library">
            <Plus className="h-4 w-4" />
            Create deployment
          </Link>
        </Button>
        <Button asChild size="lg" variant="secondary" className="gap-2 px-6">
          <Link to="/library" search={{ tab: 'examples' }}>
            <GitBranch className="h-4 w-4" />
            Browse templates
          </Link>
        </Button>
      </div>
    </div>
  );
}
