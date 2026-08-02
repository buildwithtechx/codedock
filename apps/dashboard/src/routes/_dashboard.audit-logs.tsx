import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/audit-logs')({
  beforeLoad: () => {
    throw redirect({ to: '/settings', search: { tab: 'audit' } as never });
  },
});
