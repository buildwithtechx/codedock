import { Link } from '@tanstack/react-router';
import { FileQuestion, LayoutDashboard } from 'lucide-react';
import { Button } from './button';
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from './empty';

export function NotFoundComponent() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Empty className="border-none">
        <EmptyHeader>
          <EmptyMedia
            variant="icon"
            className="size-16 rounded-2xl [&_svg:not([class*='size-'])]:size-7"
          >
            <FileQuestion />
          </EmptyMedia>
          <EmptyTitle className="text-xl">Page not found</EmptyTitle>
          <EmptyDescription>
            Sorry, we couldn't find the page you're looking for. The link might be broken, or the
            page may have been removed.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <Button asChild>
            <Link to="/">
              <LayoutDashboard className="mr-2 h-4 w-4" />
              Go to Dashboard
            </Link>
          </Button>
        </EmptyContent>
      </Empty>
    </div>
  );
}
