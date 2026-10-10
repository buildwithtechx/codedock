import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { DeployWizard } from '#/features/deploy/deploy-wizard';

export const Route = createFileRoute('/_dashboard/deploy/$slug')({
  component: DeploySlugPage,
  validateSearch: z.object({
    name: z.string().optional(),
    branch: z.string().optional(),
    dir: z.string().optional(),
    template: z.string().optional(),
  }),
});

function DeploySlugPage() {
  const { slug } = Route.useParams();
  const search = Route.useSearch();
  const decoded = decodeURIComponent(slug);
  const derivedName =
    decoded
      .split('/')
      .pop()
      ?.replace(/\.git$/, '') || 'my-project';
  const projectName = search.name || derivedName;
  const branch = search.branch || 'main';
  const repoUrl = decoded.startsWith('http') ? decoded : `https://github.com/${decoded}`;

  return (
    <div className="space-y-6">
      <DeployWizard
        initialProjectName={projectName}
        initialRepoUrl={repoUrl}
        initialBranch={branch}
      />
    </div>
  );
}
