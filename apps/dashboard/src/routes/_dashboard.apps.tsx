import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/apps')({
  beforeLoad: () => {
    throw redirect({ to: '/projects' });
  },
});
