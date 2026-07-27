import { createFileRoute } from '@tanstack/react-router';
import { UpdatesPage } from '#/features/settings/update-settings';

export const Route = createFileRoute('/_dashboard/updates')({
  component: UpdatesPage,
});
