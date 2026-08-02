import { useNavigate } from '@tanstack/react-router';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { useAuthStore } from '#/stores/auth-store';
import { HomeAppInventory } from './home-app-inventory';
import { HomeNextStep } from './home-next-step';
import { HomeProjectList } from './home-project-list';
import { HomeRuntimeSummary } from './home-runtime-summary';
import { HomeShortcuts } from './home-shortcuts';

export function HomeOverview() {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const { data, isLoading, isError, refetch } = useListCanvasSummaries();
  const projects = data?.data || [];
  const firstName = user?.name?.trim().split(/\s+/)[0] || 'there';
  const hour = new Date().getHours();
  const greeting = hour < 12 ? 'Good morning' : hour < 18 ? 'Good afternoon' : 'Good evening';

  return (
    <div className="space-y-6">
      <PageHeader
        title={`${greeting}, ${firstName}`}
        description="Here is what is happening across your projects."
      />
      <PageFrame
        rail={
          <div className="space-y-4">
            <HomeRuntimeSummary projects={projects} isLoading={isLoading} isUnavailable={isError} />
            <HomeNextStep hasProjects={projects.length > 0} />
            <HomeAppInventory projects={projects} isLoading={isLoading} />
          </div>
        }
      >
        <div className="space-y-6">
          {isError ? (
            <QueryErrorState
              title="Workspace overview is unavailable"
              description="Codedock could not load the active workspace."
              onRetry={() => void refetch()}
            />
          ) : (
            <>
              <HomeProjectList
                projects={projects}
                isLoading={isLoading}
                onCreateProject={() => void navigate({ to: '/projects/new' })}
              />
              <HomeShortcuts />
            </>
          )}
        </div>
      </PageFrame>
    </div>
  );
}
