import { Link } from '@tanstack/react-router';
import { Activity, ChevronLeft, ChevronRight, Plus } from 'lucide-react';
import { useDeferredValue, useMemo, useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useListProjects } from '#/features/projects';
import { useListByOrganization } from '#/hooks/use-deployments';
import { DeploymentEmptyState } from './deployment-empty-state';
import { DeploymentFilters } from './deployment-filters';
import { DeploymentHistoryCard } from './deployment-history-card';
import { DeploymentOverview } from './deployment-overview';
import { computeDeploymentSummary } from './deployment-summary-compute';

export function DeploymentDirectory({ hideHeader = false }: { hideHeader?: boolean } = {}) {
  const [projectId, setProjectId] = useState('all');
  const [status, setStatus] = useState('all');
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [days, setDays] = useState(30);
  const pageSize = 20;

  const deferredSearch = useDeferredValue(search);
  const filters = {
    projectId: projectId === 'all' ? undefined : projectId,
    status: status === 'all' ? undefined : status,
    search: deferredSearch.trim() || undefined,
    days,
    page,
    limit: pageSize,
  };

  const { data: deploymentsResponse, isLoading, isError, refetch } = useListByOrganization(filters);
  const { data: projectsResponse } = useListProjects({ limit: 100 });

  const deployments = deploymentsResponse?.data?.records || [];
  const totalDeployments = deploymentsResponse?.data?.total ?? deployments.length;
  const projects = projectsResponse?.data?.records || [];
  const totalProjects = projectsResponse?.data?.total ?? projects.length;
  const totalPages = Math.max(1, Math.ceil(totalDeployments / pageSize));
  const hasFilters = projectId !== 'all' || status !== 'all' || search.trim() !== '';

  const summary = useMemo(() => computeDeploymentSummary(deployments, days), [deployments, days]);

  const handleClearFilters = () => {
    setProjectId('all');
    setStatus('all');
    setSearch('');
    setPage(1);
  };

  return (
    <div className="space-y-6">
      {!hideHeader && (
        <PageHeader
          title="Deployments"
          description={
            isLoading
              ? 'Loading release activity...'
              : `${totalDeployments} total across ${totalProjects} project${
                  totalProjects === 1 ? '' : 's'
                }`
          }
          action={
            <Button asChild className="gap-2">
              <Link to="/library">
                <Plus className="h-4 w-4" />
                Deploy
              </Link>
            </Button>
          }
        />
      )}

      <DeploymentFilters
        search={search}
        projectId={projectId}
        status={status}
        projects={projects}
        onSearchChange={(val) => {
          setSearch(val);
          setPage(1);
        }}
        onProjectChange={(val) => {
          setProjectId(val);
          setPage(1);
        }}
        onStatusChange={(val) => {
          setStatus(val);
          setPage(1);
        }}
      />

      <DeploymentOverview
        summary={summary}
        loading={isLoading}
        days={days}
        onDaysChange={(d) => {
          setDays(d);
          setPage(1);
        }}
        onFilterChange={(f) => {
          setStatus(f);
          setPage(1);
        }}
        onProjectChange={(pId) => {
          setProjectId(pId);
          setStatus('all');
          setPage(1);
        }}
        onRefresh={() => void refetch()}
      />

      <div className="space-y-3">
        <h2 className="font-medium text-base text-foreground">Deployment history</h2>

        {isError ? (
          <QueryErrorState
            title="Deployment history is unavailable"
            description="Codedock could not load releases for the active workspace."
            onRetry={() => void refetch()}
          />
        ) : isLoading ? (
          <div className="flex min-h-[20rem] items-center justify-center">
            <Activity className="h-5 w-5 animate-pulse text-muted-foreground" />
          </div>
        ) : deployments.length === 0 ? (
          <DeploymentEmptyState hasFilters={hasFilters} onClear={handleClearFilters} />
        ) : (
          <div className="space-y-4">
            <div className="divide-y divide-border/50 rounded-2xl border border-border/50 bg-card">
              {deployments.map((deployment) => (
                <DeploymentHistoryCard key={deployment.id} deployment={deployment} />
              ))}
            </div>

            {totalDeployments > 0 && (
              <nav
                aria-label="Deployment pages"
                className="flex flex-wrap items-center justify-between gap-3 pt-2 text-muted-foreground text-xs"
              >
                <span>
                  {(page - 1) * pageSize + 1}–{Math.min(page * pageSize, totalDeployments)} of{' '}
                  {totalDeployments} deployments
                </span>
                <div className="flex items-center gap-3">
                  <Button
                    type="button"
                    variant="secondary"
                    size="icon"
                    disabled={isLoading || page <= 1}
                    onClick={() => setPage((p) => p - 1)}
                    aria-label="Previous page"
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <span>
                    Page {page} of {totalPages}
                  </span>
                  <Button
                    type="button"
                    variant="secondary"
                    size="icon"
                    disabled={isLoading || page >= totalPages}
                    onClick={() => setPage((p) => p + 1)}
                    aria-label="Next page"
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              </nav>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
