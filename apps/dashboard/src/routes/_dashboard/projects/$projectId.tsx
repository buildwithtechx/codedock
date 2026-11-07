import { createFileRoute, Outlet, useLocation } from '@tanstack/react-router';
import { MobileContextNav } from '#/components/layout/mobile-context-nav';
import { PageFrame } from '#/components/layout/page-frame';
import {
  getProjectContextNavigation,
  ProjectContextSidebar,
} from '#/features/projects/project-context-sidebar';

export const Route = createFileRoute('/_dashboard/projects/$projectId')({
  component: ProjectRouteLayout,
});

function ProjectRouteLayout() {
  const { projectId } = Route.useParams();
  const location = useLocation();
  const navigation = getProjectContextNavigation(projectId).map((item) => ({
    ...item,
    active: item.exact
      ? location.pathname === item.to || location.pathname === `${item.to}/`
      : location.pathname.startsWith(item.to),
  }));

  return (
    <div className="space-y-5">
      <MobileContextNav items={navigation} label="Project sections" />
      <PageFrame rail={<ProjectContextSidebar projectId={projectId} />}>
        <Outlet />
      </PageFrame>
    </div>
  );
}
