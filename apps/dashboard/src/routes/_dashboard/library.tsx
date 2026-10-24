import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { LibraryPage } from '#/features/library';

export const Route = createFileRoute('/_dashboard/library')({
  component: LibraryRoutePage,
  validateSearch: z.object({
    tab: z.enum(['repositories', 'folder', 'url', 'apps', 'examples']).optional(),
  }),
});

function LibraryRoutePage() {
  const search = Route.useSearch();
  return <LibraryPage initialTab={search.tab} />;
}
