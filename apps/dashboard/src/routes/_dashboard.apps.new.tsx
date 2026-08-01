import { createFileRoute } from '@tanstack/react-router';
import { AppCreationPage } from '#/features/dashboard/app-creation-page';

export const Route = createFileRoute('/_dashboard/apps/new')({
  component: AppCreationPage,
});
