import { Link } from '@tanstack/react-router';
import { Activity, ArrowRight, CheckCircle2, CircleX, Rocket, Search, Zap } from 'lucide-react';
import { useDeferredValue, useState } from 'react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useListProjects } from '#/features/projects';
import type { OrganizationDeployment } from '#/features/services';
import { useListByOrganization } from '#/hooks/use-deployments';
import { DeploymentEmptyState } from './deployment-empty-state';

const statusOptions = [
  ['all', 'All'],
  ['READY', 'Ready'],
  ['ACTIVE', 'Active'],
  ['BUILDING', 'Building'],
  ['FAILED', 'Failed'],
] as const;

const activeStatuses = new Set(['PENDING', 'CLONING', 'PULLING', 'BUILDING']);
const successfulStatuses = new Set(['READY', 'ACTIVE', 'SUCCESS']);

const statusTone = (status: string) => {
  const normalized = status.toUpperCase();
  if (successfulStatuses.has(normalized)) return 'bg-emerald-500';
  if (normalized === 'FAILED') return 'bg-rose-500';
  if (activeStatuses.has(normalized)) return 'bg-amber-400';
  return 'bg-muted-foreground/50';
};

const statusLabel = (status: string) =>
  status
    .toLowerCase()
    .replace(/_/g, ' ')
    .replace(/\b\w/g, (letter) => letter.toUpperCase());

export function DeploymentDirectory() {
  const [projectId, setProjectId] = useState('all');
  const [status, setStatus] = useState('all');
  const [search, setSearch] = useState('');
  const deferredSearch = useDeferredValue(search);
  const filters = {
    projectId: projectId === 'all' ? undefined : projectId,
    status: status === 'all' ? undefined : status,
    search: deferredSearch.trim() || undefined,
    limit: 50,
  };
  const { data: deploymentsResponse, isLoading, isError, refetch } = useListByOrganization(filters);
  const { data: projectsResponse } = useListProjects();
  const deployments = deploymentsResponse?.data?.records || [];
  const projects = projectsResponse?.data?.records || [];
  const hasFilters = projectId !== 'all' || status !== 'all' || search.trim() !== '';

  return (
    <div className="space-y-6">
      <PageHeader
        title="Deployments"
        description={
          isLoading
            ? 'Loading release activity...'
            : `${deployments.length} total across ${projects.length} project${
                projects.length === 1 ? '' : 's'
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
            onProjectChange={setProjectId}
            onSearchChange={setSearch}
            onStatusChange={setStatus}
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
              }}
            />
          ) : (
            <DeploymentTable deployments={deployments} />
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
          className="h-10 bg-card pl-10"
        />
      </div>
      <Select value={projectId} onValueChange={onProjectChange}>
        <SelectTrigger className="h-10 w-full bg-card sm:w-44">
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
      <div className="inline-flex max-w-full flex-wrap items-center gap-1 rounded-xl bg-muted/45 p-1">
        {statusOptions.map(([value, label]) => (
          <button
            key={value}
            type="button"
            onClick={() => onStatusChange(value)}
            className={`inline-flex h-8 items-center rounded-lg px-3 font-medium text-xs transition-colors ${
              status === value
                ? 'bg-card text-foreground shadow-sm'
                : 'text-muted-foreground hover:bg-background/70 hover:text-foreground'
            }`}
          >
            {label}
          </button>
        ))}
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
      </aside>
      <DeploymentGetStarted />
    </div>
  );
}

function DeploymentGetStarted() {
  return (
    <section className="rounded-2xl border border-primary/10 bg-gradient-to-br from-primary/10 via-primary/5 to-transparent p-5">
      <div className="flex items-center gap-2">
        <Zap className="h-4 w-4 text-primary" />
        <h2 className="font-semibold text-sm">Get started</h2>
      </div>
      <p className="mt-3 text-muted-foreground text-sm leading-6">
        Connect a repository or choose a template to create the first deployment.
      </p>
      <Link
        to="/projects/new"
        className="mt-4 inline-flex items-center gap-1.5 font-medium text-sm transition-colors hover:text-primary"
      >
        Create project
        <ArrowRight className="h-3.5 w-3.5" />
      </Link>
    </section>
  );
}

function DeploymentTable({ deployments }: { deployments: OrganizationDeployment[] }) {
  return (
    <section className="overflow-x-auto rounded-2xl bg-card">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Release</TableHead>
            <TableHead>Project</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Branch</TableHead>
            <TableHead>Trigger</TableHead>
            <TableHead>Created</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {deployments.map((deployment) => (
            <TableRow key={deployment.id}>
              <TableCell>
                <div className="min-w-36">
                  <p className="font-medium text-sm">{deployment.serviceName || 'Unknown app'}</p>
                  <p className="mt-0.5 truncate font-mono text-muted-foreground text-xs">
                    {deployment.commitHash?.slice(0, 7) || deployment.id.slice(0, 7)}
                  </p>
                </div>
              </TableCell>
              <TableCell className="font-medium text-sm">{deployment.projectName}</TableCell>
              <TableCell>
                <span className="inline-flex items-center gap-2 text-sm">
                  <span className={`h-1.5 w-1.5 rounded-full ${statusTone(deployment.status)}`} />
                  {statusLabel(deployment.status)}
                </span>
              </TableCell>
              <TableCell className="max-w-36 truncate font-mono text-xs">
                {deployment.branch || '–'}
              </TableCell>
              <TableCell className="max-w-32 truncate text-sm">
                {deployment.trigger || '–'}
              </TableCell>
              <TableCell className="whitespace-nowrap text-muted-foreground text-sm">
                {new Date(deployment.createdAt).toLocaleString()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </section>
  );
}
