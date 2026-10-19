import { createFileRoute } from '@tanstack/react-router';
import { InstallWizard } from '#/features/apps';

export const Route = createFileRoute('/_dashboard/apps/new/$appId')({
  component: AppInstallPage,
});

function AppInstallPage() {
  const { appId } = Route.useParams();
  return <InstallWizard appId={appId} />;
}
