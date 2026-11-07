import { createFileRoute } from '@tanstack/react-router';
import { ProjectDirectory } from '#/features/dashboard/project-directory';

export const Route = createFileRoute('/_dashboard/projects/')({
  component: ProjectDirectory,
});
