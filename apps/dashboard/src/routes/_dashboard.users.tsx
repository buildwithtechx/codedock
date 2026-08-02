import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/users')({
  beforeLoad: () => {
    throw redirect({ to: '/settings', search: { tab: 'team' } as never });
  },
});
