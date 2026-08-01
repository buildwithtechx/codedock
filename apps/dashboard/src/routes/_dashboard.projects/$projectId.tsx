import { createFileRoute, Outlet } from '@tanstack/react-router';
import { ProjectContextSidebar } from '#/features/projects/project-context-sidebar';

export const Route = createFileRoute('/_dashboard/projects/$projectId')({
  component: ProjectRouteLayout,
});

function ProjectRouteLayout() {
  const { projectId } = Route.useParams();

  return (
    <div className="grid min-w-0 gap-6 xl:grid-cols-[minmax(0,1fr)_18rem]">
      <div className="min-w-0">
        <Outlet />
      </div>
      <div className="hidden xl:block">
        <ProjectContextSidebar projectId={projectId} />
      </div>
    </div>
  );
}
