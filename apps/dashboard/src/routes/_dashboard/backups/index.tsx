import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { BackupsList } from '#/features/backups/backups-list';

export const Route = createFileRoute('/_dashboard/backups/')({
  validateSearch: z.object({
    add: z.string().optional(),
  }),
  component: BackupsList,
});
