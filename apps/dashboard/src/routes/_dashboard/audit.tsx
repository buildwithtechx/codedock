import { createFileRoute, redirect } from '@tanstack/react-router';
import { z } from 'zod';
import { AuditLogList } from '#/features/audit/audit-log-list';

export const Route = createFileRoute('/_dashboard/audit')({
  beforeLoad: () => {
    throw redirect({
      to: '/team',
      search: { tab: 'audit' },
      replace: true,
    });
  },
  validateSearch: z.object({
    q: z.string().optional(),
    category: z.string().optional(),
    actor: z.string().optional(),
    source: z.string().optional(),
    period: z.string().optional(),
  }),
  component: AuditPage,
});

export function AuditPage() {
  return <AuditLogList />;
}
