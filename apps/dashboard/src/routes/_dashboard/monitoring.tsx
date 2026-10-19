import { createFileRoute } from '@tanstack/react-router';
import { IssuesView } from '#/features/issues';

export const Route = createFileRoute('/_dashboard/monitoring')({
  component: IssuesView,
});
