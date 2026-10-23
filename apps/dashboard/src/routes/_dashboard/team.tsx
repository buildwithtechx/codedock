import { createFileRoute } from '@tanstack/react-router';
import { TeamLayout } from '#/features/teams/team-page';

export const Route = createFileRoute('/_dashboard/team')({
  component: TeamPage,
});

export function TeamPage() {
  return <TeamLayout />;
}
