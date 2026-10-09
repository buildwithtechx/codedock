import { createFileRoute } from '@tanstack/react-router';
import { DeployWizard } from '#/features/deploy/deploy-wizard';

export const Route = createFileRoute('/_dashboard/deploy/$slug')({
  component: DeploySlugPage,
});

function DeploySlugPage() {
  const { slug } = Route.useParams();
  const decoded = decodeURIComponent(slug);
  const projectName =
    decoded
      .split('/')
      .pop()
      ?.replace(/\.git$/, '') || 'my-project';
  const repoUrl = decoded.startsWith('http') ? decoded : `https://github.com/${decoded}`;

  return (
    <div className="space-y-6">
      <DeployWizard initialProjectName={projectName} initialRepoUrl={repoUrl} />
    </div>
  );
}
