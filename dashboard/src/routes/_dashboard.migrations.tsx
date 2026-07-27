import { createFileRoute } from '@tanstack/react-router';
import { MigrationSettings } from '#/features/settings/migration-settings';

export const Route = createFileRoute('/_dashboard/migrations')({
  component: () => <MigrationSettings />,
});
