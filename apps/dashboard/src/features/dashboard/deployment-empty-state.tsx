import { Link } from '@tanstack/react-router';
import { LayoutTemplate, Plus, Rocket, SearchX } from 'lucide-react';
import { Button } from '#/components/ui/button';

export function DeploymentEmptyState({
  hasFilters,
  onClear,
}: {
  hasFilters: boolean;
  onClear: () => void;
}) {
  if (hasFilters) {
    return (
      <section className="rounded-2xl bg-card px-6 py-16 text-center">
        <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-muted/70 text-muted-foreground">
          <SearchX className="h-6 w-6" />
        </span>
        <h2 className="mt-5 font-medium text-foreground/90 text-lg">
          No releases match these filters
        </h2>
        <p className="mx-auto mt-2 max-w-sm text-muted-foreground/75 text-sm leading-6">
          Change or clear a filter to look across a different set of releases.
        </p>
        <Button variant="outline" className="mt-6" onClick={onClear}>
          Clear filters
        </Button>
      </section>
    );
  }

  return (
    <section className="rounded-2xl bg-card px-6 pb-10 text-center">
      <div className="relative mx-auto h-44 w-64" aria-hidden="true">
        <span className="absolute top-12 left-10 h-21 w-30 rounded-2xl bg-muted/65" />
        <span className="absolute top-7 left-16 h-21 w-30 rounded-2xl bg-muted/80" />
        <span className="absolute top-2 left-22 flex h-21 w-30 flex-col gap-2 rounded-2xl bg-background p-4 text-left">
          <span className="flex gap-1.5">
            <i className="h-2 w-2 rounded-full bg-rose-400/70" />
            <i className="h-2 w-2 rounded-full bg-amber-400/70" />
            <i className="h-2 w-2 rounded-full bg-emerald-400/70" />
          </span>
          <i className="h-1.5 w-13 rounded-full bg-muted" />
          <i className="h-1.5 w-18 rounded-full bg-muted/75" />
          <span className="mt-auto flex h-8 w-8 items-center justify-center rounded-lg bg-muted text-muted-foreground">
            <Rocket className="h-4 w-4" />
          </span>
        </span>
        <span className="absolute top-14 right-3 flex h-12 w-12 items-center justify-center rounded-full border border-primary/45 border-dashed bg-card text-primary">
          <Plus className="h-5 w-5" />
        </span>
      </div>
      <h2 className="font-medium text-foreground/90 text-xl tracking-[-0.02em]">
        No deployments yet
      </h2>
      <p className="mx-auto mt-2 max-w-sm text-muted-foreground/75 text-sm leading-6">
        Deploy an app and its build status, release history, and commit details will appear here.
      </p>
      <div className="mt-7 flex flex-col items-center justify-center gap-2 sm:flex-row">
        <Link to="/apps/new">
          <Button className="gap-2 px-5">
            <Plus className="h-4 w-4" />
            Deploy app
          </Button>
        </Link>
        <Link to="/projects/new" search={{ template: 'one-click' }}>
          <Button variant="secondary" className="gap-2 px-5">
            <LayoutTemplate className="h-4 w-4" />
            Browse templates
          </Button>
        </Link>
      </div>
    </section>
  );
}
