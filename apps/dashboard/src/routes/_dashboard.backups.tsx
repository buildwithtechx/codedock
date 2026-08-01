import { createFileRoute } from '@tanstack/react-router';
import { BackupsList } from '#/features/backups/backups-list';

export const Route = createFileRoute('/_dashboard/backups')({
  component: BackupsList,
});
