import { createFileRoute } from '@tanstack/react-router';
import { MaintenancePage } from '#/features/settings/maintenance-settings';

export const Route = createFileRoute('/_dashboard/maintenance')({
  component: MaintenancePage,
});
