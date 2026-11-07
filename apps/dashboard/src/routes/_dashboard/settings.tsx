import { createFileRoute } from '@tanstack/react-router';
import { SettingsLayout } from '#/features/settings/settings-page';

export const Route = createFileRoute('/_dashboard/settings')({
  component: SettingsLayout,
});
