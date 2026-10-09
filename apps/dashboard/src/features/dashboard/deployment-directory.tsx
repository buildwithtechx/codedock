import { Link } from '@tanstack/react-router';
import { Activity, ArrowRight, CheckCircle2, CircleX, Rocket, Search, Zap } from 'lucide-react';
import { useDeferredValue, useState } from 'react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useListProjects } from '#/features/projects';
import type { OrganizationDeployment } from '#/features/services';
import { useListByOrganization } from '#/hooks/use-deployments';
import { DeploymentEmptyState } from './deployment-empty-state';
import { DeploymentHistoryCard } from './deployment-history-card';

const statusOptions = [
  ['all', 'All'],
  ['READY', 'Ready'],
  ['ACTIVE', 'Active'],
  ['BUILDING', 'Building'],
  ['FAILED', 'Failed'],
] as const;

const activeStatuses = new Set(['PENDING', 'CLONING', 'PULLING', 'BUILDING']);
const successfulStatuses = new Set(['READY', 'ACTIVE', 'SUCCESS']);

export function DeploymentDirectory() {
  const [projectId, setProjectId] = useState('all');
  const [status, setStatus] = useState('all');
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = 20;
  const deferredSearch = useDeferredValue(search);
  const filters = {
    projectId: projectId === 'all' ? undefined : projectId,
    status: status === 'all' ? undefined : status,
    search: deferredSearch.trim() || undefined,
    page,
    limit: pageSize,
  };
  const { data: deploymentsResponse, isLoading, isError, refetch } = useListByOrganization(filters);
  const { data: projectsResponse } = useListProjects({ limit: 100 });
  const deployments = deploymentsResponse?.data?.records || [];
  const totalDeployments = deploymentsResponse?.data?.total ?? deployments.length;
  const projects = projectsResponse?.data?.records || [];
  const totalProjects = projectsResponse?.data?.total ?? projects.length;
  const totalPages = Math.ceil(totalDeployments / pageSize);
  const hasFilters = projectId !== 'all' || status !== 'all' || search.trim() !== '';

  return (
    <div className="space-y-6">
      <PageHeader
        title="Deployments"
        description={
          isLoading
            ? 'Loading release activity...'
            : `${totalDeployments} total across ${totalProjects} project${
                totalProjects === 1 ? '' : 's'
              }`
        }
      />
      <PageFrame rail={<DeploymentSummary deployments={deployments} isLoading={isLoading} />}>
        <div className="space-y-4">
          <DeploymentFilters
            projectId={projectId}
            projects={projects}
            search={search}
            status={status}
            onProjectChange={(val) => {
              setProjectId(val);
              setPage(1);
            }}
            onSearchChange={(val) => {
              setSearch(val);
              setPage(1);
            }}
            onStatusChange={(val) => {
              setStatus(val);
              setPage(1);
            }}
          />
          {isError ? (
            <QueryErrorState
              title="Deployment history is unavailable"
              description="Codedock could not load releases for the active workspace."
              onRetry={() => void refetch()}
            />
          ) : isLoading ? (
            <div className="flex min-h-[25rem] items-center justify-center">
              <Activity className="h-5 w-5 animate-pulse text-muted-foreground" />
            </div>
          ) : deployments.length === 0 ? (
            <DeploymentEmptyState
              hasFilters={hasFilters}
              onClear={() => {
                setProjectId('all');
                setStatus('all');
                setSearch('');
                setPage(1);
              }}
            />
          ) : (
            <>
              <div className="divide-y divide-border/50 rounded-2xl bg-card">
                {deployments.map((deployment) => (
                  <DeploymentHistoryCard key={deployment.id} deployment={deployment} />
                ))}
              </div>
              {totalPages > 1 && (
                <div className="flex items-center justify-between border-border/40 border-t pt-4">
                  <p className="text-muted-foreground text-sm">
                    Showing {(page - 1) * pageSize + 1} to{' '}
                    {Math.min(page * pageSize, totalDeployments)} of {totalDeployments}
                  </p>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      disabled={page <= 1}
                    >
                      Previous
                    </Button>
                    <span className="text-muted-foreground text-sm">
                      Page {page} of {totalPages}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                      disabled={page >= totalPages}
                    >
                      Next
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </PageFrame>
    </div>
  );
}

function DeploymentFilters({
  projectId,
  projects,
  search,
  status,
  onProjectChange,
  onSearchChange,
  onStatusChange,
}: {
  projectId: string;
  projects: { id: string; name: string }[];
  search: string;
  status: string;
  onProjectChange: (value: string) => void;
  onSearchChange: (value: string) => void;
  onStatusChange: (value: string) => void;
}) {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-center">
      <div className="relative w-full sm:min-w-55 sm:flex-1">
        <Search className="pointer-events-none absolute top-1/2 left-3.5 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          value={search}
          onChange={(event) => onSearchChange(event.target.value)}
          placeholder="Search app, project, branch, or commit"
          className="h-10 bg-muted/60 ps-10 pe-4"
        />
      </div>
      <Select value={projectId} onValueChange={onProjectChange}>
        <SelectTrigger className="w-full bg-card sm:w-44">
          <SelectValue placeholder="All projects" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All projects</SelectItem>
          {projects.map((project) => (
            <SelectItem key={project.id} value={project.id}>
              {project.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <div className="max-w-full shrink-0 overflow-x-auto [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        <div className="inline-flex items-center gap-1">
          {statusOptions.map(([value, label]) => (
            <button
              key={value}
              type="button"
              aria-pressed={status === value}
              onClick={() => onStatusChange(value)}
              className={`inline-flex h-10 shrink-0 items-center whitespace-nowrap rounded-lg px-4 font-medium text-sm transition-colors ${
                status === value
                  ? 'bg-foreground text-background'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground'
              }`}
            >
              {label}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}

function DeploymentSummary({
  deployments,
  isLoading,
}: {
  deployments: OrganizationDeployment[];
  isLoading: boolean;
}) {
  const successful = deployments.filter((deployment) =>
    successfulStatuses.has(deployment.status.toUpperCase())
  ).length;
  const active = deployments.filter((deployment) =>
    activeStatuses.has(deployment.status.toUpperCase())
  ).length;
  const failed = deployments.filter(
    (deployment) => deployment.status.toUpperCase() === 'FAILED'
  ).length;
  const items = [
    { label: 'Total', value: deployments.length, icon: Rocket, tone: 'text-primary bg-primary/12' },
    {
      label: 'Successful',
      value: successful,
      icon: CheckCircle2,
      tone: 'text-emerald-500 bg-emerald-500/12',
    },
    { label: 'Failed', value: failed, icon: CircleX, tone: 'text-rose-500 bg-rose-500/12' },
  ];

  return (
    <div className="space-y-4">
      <aside className="rounded-2xl bg-card p-5">
        <div className="flex items-center gap-2">
          <Activity className="h-4 w-4 text-muted-foreground" />
          <h2 className="font-semibold text-sm">Overview</h2>
        </div>
        <div className="mt-5 space-y-3.5">
          {items.map((item) => (
            <div key={item.label} className="flex items-center justify-between">
              <span className="flex items-center gap-2.5 text-muted-foreground text-sm">
                <span
                  className={`flex h-8 w-8 items-center justify-center rounded-lg ${item.tone}`}
                >
                  <item.icon className="h-4 w-4" />
                </span>
                {item.label}
              </span>
              <span className="font-semibold">{isLoading ? '–' : item.value}</span>
            </div>
          ))}
          {active > 0 && (
            <div className="flex items-center justify-between">
              <span className="flex items-center gap-2.5 text-muted-foreground text-sm">
                <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-500/12 text-amber-500">
                  <Activity className="h-4 w-4" />
                </span>
                In progress
              </span>
              <span className="font-semibold">{isLoading ? '–' : active}</span>
            </div>
          )}
        </div>
        {!isLoading && deployments.length > 0 && (
          <div className="mt-4 flex h-1.5 overflow-hidden rounded-full bg-muted/40">
            {successful > 0 && (
              <div
                className="bg-emerald-500"
                style={{ width: `${(successful / deployments.length) * 100}%` }}
              />
            )}
            {failed > 0 && (
              <div
                className="bg-rose-500"
                style={{ width: `${(failed / deployments.length) * 100}%` }}
              />
            )}
            {active > 0 && (
              <div
                className="bg-amber-400"
                style={{ width: `${(active / deployments.length) * 100}%` }}
              />
            )}
          </div>
        )}
      </aside>
      <DeploymentTip hasDeployments={deployments.length > 0} />
    </div>
  );
}

function DeploymentTip({ hasDeployments }: { hasDeployments: boolean }) {
  return (
    <section className="rounded-2xl border border-primary/10 bg-gradient-to-br from-primary/10 via-primary/5 to-transparent p-5">
      <div className="flex items-center gap-2">
        <Zap className="h-4 w-4 text-primary" />
        <h2 className="font-semibold text-sm">{hasDeployments ? 'Auto deploy' : 'Get started'}</h2>
      </div>
      <p className="mt-3 text-muted-foreground text-sm leading-6">
        {hasDeployments
          ? 'Push to a connected branch and Codedock builds and ships the new release automatically.'
          : 'Connect a repository or choose a template to create the first deployment.'}
      </p>
      <Link
        to={hasDeployments ? '/projects' : '/projects/new'}
        className="mt-4 inline-flex items-center gap-1.5 font-medium text-sm transition-colors hover:text-primary"
      >
        {hasDeployments ? 'View projects' : 'Create project'}
        <ArrowRight className="h-3.5 w-3.5" />
      </Link>
    </section>
  );
}
