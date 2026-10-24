import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/dns')({
  beforeLoad: () => {
    throw redirect({ to: '/domains' });
  },
});
