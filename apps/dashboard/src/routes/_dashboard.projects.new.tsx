import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { ProjectCreationPage } from '#/features/projects/project-creation-page';

export const Route = createFileRoute('/_dashboard/projects/new')({
  component: ProjectCreationPage,
  validateSearch: z.object({
    template: z.enum(['one-click', 'examples']).optional(),
  }),
});
