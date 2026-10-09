import { Link } from '@tanstack/react-router';
import { BookOpen, ExternalLink, LayoutTemplate, Plus, Rocket, SearchX } from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '#/components/ui/empty';

const DOCS_URL = 'https://docs.codedock.run';

export function DeploymentEmptyState({
  hasFilters,
  onClear,
}: {
  hasFilters: boolean;
  onClear: () => void;
}) {
  if (hasFilters) {
    return (
      <Empty>
        <EmptyHeader>
          <EmptyMedia
            variant="icon"
            className="size-16 rounded-2xl [&_svg:not([class*='size-'])]:size-7"
          >
            <SearchX />
          </EmptyMedia>
          <EmptyTitle className="text-xl">No releases match these filters</EmptyTitle>
          <EmptyDescription>
            Change or clear a filter to look across a different set of releases.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent className="flex-row justify-center">
          <Button variant="outline" onClick={onClear}>
            Clear filters
          </Button>
          <Button asChild variant="secondary" className="gap-2 px-5">
            <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
              <BookOpen className="h-4 w-4" />
              Docs
              <ExternalLink className="h-3.5 w-3.5 opacity-60" />
            </a>
          </Button>
        </EmptyContent>
      </Empty>
    );
  }

  return (
    <Empty className="border-none py-12">
      <EmptyHeader>
        <EmptyMedia
          variant="icon"
          className="size-16 rounded-2xl [&_svg:not([class*='size-'])]:size-7"
        >
          <Rocket />
        </EmptyMedia>
        <EmptyTitle className="text-xl">No deployments yet</EmptyTitle>
        <EmptyDescription>
          Deploy an app and its build status, release history, and commit details will appear here.
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent className="flex-row flex-wrap justify-center">
        <Button asChild className="gap-2 px-5">
          <Link to="/projects/new">
            <Plus className="h-4 w-4" />
            Deploy app
          </Link>
        </Button>
        <Button asChild variant="secondary" className="gap-2 px-5">
          <Link to="/projects/new" search={{ template: 'one-click' }}>
            <LayoutTemplate className="h-4 w-4" />
            Browse templates
          </Link>
        </Button>
        <Button asChild variant="secondary" className="gap-2 px-5">
          <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
            <BookOpen className="h-4 w-4" />
            Docs
            <ExternalLink className="h-3.5 w-3.5 opacity-60" />
          </a>
        </Button>
      </EmptyContent>
    </Empty>
  );
}
