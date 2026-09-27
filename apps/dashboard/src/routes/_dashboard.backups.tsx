import { createFileRoute } from '@tanstack/react-router';
import { BackupsList } from '#/features/backups/backups-list';

export const Route = createFileRoute('/_dashboard/backups')({
  validateSearch: (search: Record<string, unknown>) => ({
    tab: search.tab as string | undefined,
    add: search.add as string | undefined,
  }),
  component: BackupsList,
});
