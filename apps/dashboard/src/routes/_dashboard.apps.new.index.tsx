import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { NewAppTabs } from '#/features/apps';

export const Route = createFileRoute('/_dashboard/apps/new/')({
  component: NewAppPage,
  validateSearch: z.object({
    app: z.string().optional(),
  }),
});

function NewAppPage() {
  const search = Route.useSearch();
  return <NewAppTabs deepLinkAppId={search.app} />;
}
