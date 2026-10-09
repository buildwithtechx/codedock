import { createFileRoute } from '@tanstack/react-router';
import { DeploymentDetail } from '#/features/dashboard/deployment-detail';

export const Route = createFileRoute('/_dashboard/deployments/$deploymentId')({
  component: DeploymentDetailPage,
});

function DeploymentDetailPage() {
  const { deploymentId } = Route.useParams();
  return <DeploymentDetail deploymentId={deploymentId} />;
}
