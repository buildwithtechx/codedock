import { createFileRoute } from '@tanstack/react-router';
import { DeploymentDirectory } from '#/features/dashboard/deployment-directory';

export const Route = createFileRoute('/_dashboard/deployments')({
  component: DeploymentDirectory,
});
