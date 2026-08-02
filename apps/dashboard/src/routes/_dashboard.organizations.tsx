import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/organizations')({
  beforeLoad: () => {
    throw redirect({ to: '/settings', search: { tab: 'team' } as never });
  },
});
