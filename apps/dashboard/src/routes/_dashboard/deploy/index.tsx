import { createFileRoute } from '@tanstack/react-router';
import { DeployWizard } from '#/features/deploy/deploy-wizard';

export const Route = createFileRoute('/_dashboard/deploy/')({
  component: DeployIndexPage,
});

function DeployIndexPage() {
  return (
    <div className="space-y-6">
      <DeployWizard />
    </div>
  );
}
