import { createFileRoute, redirect } from '@tanstack/react-router';
import { z } from 'zod';

export const Route = createFileRoute('/_dashboard/projects/$projectId/new')({
  beforeLoad: ({ params }) => {
    throw redirect({
      to: '/projects/$projectId',
      params: { projectId: params.projectId },
    });
  },
  validateSearch: z.object({
    tab: z.enum(['resources', 'one-click', 'examples']).optional(),
    resource: z.enum(['git', 'database', 'docker']).optional(),
  }),
});
