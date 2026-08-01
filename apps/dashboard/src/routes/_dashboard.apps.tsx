import { createFileRoute } from '@tanstack/react-router';
import { AppDirectory } from '#/features/dashboard/app-directory';

export const Route = createFileRoute('/_dashboard/apps')({
  component: AppDirectory,
});
