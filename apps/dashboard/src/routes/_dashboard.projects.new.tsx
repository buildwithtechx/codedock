import { createFileRoute } from '@tanstack/react-router';
import { ProjectCreationPage } from '#/features/projects/project-creation-page';

export const Route = createFileRoute('/_dashboard/projects/new')({
  component: ProjectCreationPage,
});
