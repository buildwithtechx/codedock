import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { AuditLogList } from '#/features/audit/audit-log-list';

export const Route = createFileRoute('/_dashboard/audit')({
  component: AuditPage,
  validateSearch: z.object({
    q: z.string().optional(),
    category: z.string().optional(),
    actor: z.string().optional(),
    source: z.string().optional(),
    period: z.string().optional(),
  }),
});

export function AuditPage() {
  return <AuditLogList />;
}
