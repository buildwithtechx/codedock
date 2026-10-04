import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/apps/new')({
  beforeLoad: () => {
    throw redirect({ to: '/projects/new' });
  },
});
