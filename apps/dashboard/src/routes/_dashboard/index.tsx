import { createFileRoute } from '@tanstack/react-router';
import { HomeOverview } from '#/features/dashboard/home-overview';

export const Route = createFileRoute('/_dashboard/')({
  component: DashboardIndex,
});

function DashboardIndex() {
  return <HomeOverview />;
}
