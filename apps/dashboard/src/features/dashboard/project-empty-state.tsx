import { Link } from '@tanstack/react-router';
import {
  ArrowRight,
  BookOpen,
  ExternalLink,
  FolderKanban,
  LayoutTemplate,
  Plus,
} from 'lucide-react';
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

export function ProjectEmptyState() {
  return (
    <Empty className="border-none py-12">
      <EmptyHeader>
        <EmptyMedia
          variant="icon"
          className="size-16 rounded-2xl [&_svg:not([class*='size-'])]:size-7"
        >
          <FolderKanban />
        </EmptyMedia>
        <EmptyTitle className="text-xl">Create your first project</EmptyTitle>
        <EmptyDescription>
          Bring together the services, environments, and deployment activity that belong to one
          product.
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent className="flex-row flex-wrap justify-center">
        <Button asChild className="gap-2 px-5">
          <Link to="/projects/new">
            <Plus className="h-4 w-4" />
            New project
            <ArrowRight className="h-4 w-4" />
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
