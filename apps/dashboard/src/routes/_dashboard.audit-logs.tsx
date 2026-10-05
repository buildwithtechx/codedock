import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/audit-logs')({
  beforeLoad: () => {
    throw redirect({ to: '/audit' });
  },
});
