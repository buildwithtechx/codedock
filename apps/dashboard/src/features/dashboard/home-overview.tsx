import { useNavigate } from '@tanstack/react-router';
import { PageHeader } from '#/components/layout/page-header';
import { QueryErrorState } from '#/components/query-error-state';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { useListByOrganization } from '#/hooks/use-deployments';
import { useGetPublicSettings } from '#/hooks/use-settings';
import { useAuthStore } from '#/stores/auth-store';
import { HomeAppInventory } from './home-app-inventory';
import { HomeModeCard } from './home-mode-card';
import { HomeNextStep } from './home-next-step';
import { HomeProjectList } from './home-project-list';
import { HomeRuntimeSummary } from './home-runtime-summary';
import { HomeShortcuts } from './home-shortcuts';
import { HomeSidebar } from './home-sidebar';
import { HomeUpdatesBlock, useHomeAttention } from './home-updates-block';

export function HomeOverview() {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const { data: publicSettings } = useGetPublicSettings();
  const { data, isLoading, isError, refetch } = useListCanvasSummaries();
  const deploymentsQuery = useListByOrganization({ page: 1, limit: 1 });
  const attention = useHomeAttention();

  const isCloud = Boolean(publicSettings?.data?.cloudMode);
  const projects = data?.data || [];
  const totalDeployments = deploymentsQuery.data?.data?.total ?? 0;
  const firstName = user?.name?.trim().split(/\s+/)[0] || '';
  const hour = new Date().getHours();
  const greeting = hour < 12 ? 'Good morning' : hour < 18 ? 'Good afternoon' : 'Good evening';

  return (
    <div className="space-y-6">
      <PageHeader
        title={firstName ? `${greeting}, ${firstName}` : greeting}
        description="Here is what is happening across your projects."
      />

      {isError ? (
        <QueryErrorState
          title="Workspace overview is unavailable"
          description="Codedock could not load the active workspace."
          onRetry={() => void refetch()}
        />
      ) : (
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_340px]">
          <div className="min-w-0 space-y-6">
            <HomeProjectList
              projects={projects}
              isLoading={isLoading}
              onCreateProject={() => void navigate({ to: '/library' })}
            />
            <HomeShortcuts />
          </div>

          <HomeSidebar
            tip={
              attention.cards === 0 ? (
                <HomeNextStep hasProjects={projects.length > 0} isLoading={isLoading} />
              ) : null
            }
          >
            <HomeModeCard />

            {attention.cards === 0 && (!isCloud || projects.length > 0) && (
              <HomeRuntimeSummary
                projects={projects}
                totalDeployments={totalDeployments}
                brokenCount={attention.broken.length}
                issuesLoaded={attention.loaded}
                isLoading={isLoading}
                isUnavailable={isError}
              />
            )}

            <HomeUpdatesBlock attention={attention} />
            <HomeAppInventory isLoading={isLoading} />
          </HomeSidebar>
        </div>
      )}
    </div>
  );
}
