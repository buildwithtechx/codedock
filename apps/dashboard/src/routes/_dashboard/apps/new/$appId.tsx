import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { InstallWizard } from '#/features/apps';

export const Route = createFileRoute('/_dashboard/apps/new/$appId')({
  component: AppInstallPage,
  validateSearch: z.object({
    projectId: z.string().optional(),
  }),
});

function AppInstallPage() {
  const { appId } = Route.useParams();
  const { projectId } = Route.useSearch();
  return <InstallWizard appId={appId} projectId={projectId} />;
}
