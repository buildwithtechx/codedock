import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/scheduled-tasks')({
  beforeLoad: () => {
    throw redirect({ to: '/jobs' });
  },
});
