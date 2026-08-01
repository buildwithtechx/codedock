import { createFileRoute } from '@tanstack/react-router';
import { ServerCreationPage } from '#/features/servers/server-creation-page';

export const Route = createFileRoute('/_dashboard/servers/new')({
  component: ServerCreationPage,
});
