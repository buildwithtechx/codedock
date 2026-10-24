import { createFileRoute } from '@tanstack/react-router';
import { ProjectDomainsTab } from '#/features/projects/project-domains-tab';

export const Route = createFileRoute('/_dashboard/projects/$projectId/domains')({
  component: ProjectDomainsRouteComponent,
});

function ProjectDomainsRouteComponent() {
  const { projectId } = Route.useParams();

  return <ProjectDomainsTab projectId={projectId} />;
}
