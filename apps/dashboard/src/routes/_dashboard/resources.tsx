import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/resources')({
  beforeLoad: () => {
    throw redirect({ to: '/servers' });
  },
  component: () => null,
});
