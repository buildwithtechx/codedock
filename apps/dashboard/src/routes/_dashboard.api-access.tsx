import { createFileRoute, Navigate } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/api-access')({
  component: () => <Navigate to="/settings" search={{ tab: 'tokens' }} replace />,
});
