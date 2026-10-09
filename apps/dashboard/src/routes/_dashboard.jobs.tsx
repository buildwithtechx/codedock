import { useQueries, useQueryClient } from '@tanstack/react-query';
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router';
import { Loader2, Plus, Search } from 'lucide-react';
import { useMemo, useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { jobsApi } from '#/features/jobs/api';
import { JobCard } from '#/features/jobs/job-card';
import { isFailedStatus } from '#/features/jobs/job-format';
import { JobOutputDialog } from '#/features/jobs/job-output-dialog';
import { JobsEmptyState } from '#/features/jobs/jobs-empty-state';
import { type JobStatusFilter, JobsOverview } from '#/features/jobs/jobs-overview';
import { useListAllProjects } from '#/features/projects';
import type { Job } from '#/features/services';

export const Route = createFileRoute('/_dashboard/jobs')({
  component: JobsPage,
});

function matchesStatus(job: Job, filter: JobStatusFilter): boolean {
  if (filter === 'all') return true;
  if (filter === 'running') return job.status === 'running';
  if (filter === 'failed') return isFailedStatus(job.status);
  if (filter === 'scheduled') return job.status !== 'inactive';
  return job.status === 'inactive';
}

export function JobsPage() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [query, setQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<JobStatusFilter>('all');
  const [outputJob, setOutputJob] = useState<Job | null>(null);

  const {
    data: projectsResponse,
    isLoading: isLoadingProjects,
    isError: projectsError,
    refetch: refetchProjects,
  } = useListAllProjects();
  const projects = projectsResponse || [];
  const projectNames = new Map(projects.map((project) => [project.id, project.name]));

  const taskQueries = useQueries({
    queries: projects.map((project) => ({
      queryKey: ['scheduled-tasks', 'project', project.id],
      queryFn: () => jobsApi.listByProject(project.id),
    })),
  });
  const isLoadingTasks = taskQueries.some((query) => query.isLoading);
  const failedQueries = taskQueries.flatMap((query, index) =>
    query.isError ? [{ query, project: projects[index] }] : []
  );
  const tasks = taskQueries.flatMap((query) => query.data?.data || []);
  const loading = isLoadingProjects || (isLoadingTasks && tasks.length === 0);

  const refresh = () => {
    void queryClient.invalidateQueries({ queryKey: ['scheduled-tasks'] });
  };

  const q = query.trim().toLowerCase();
  const visible = tasks.filter(
    (job) => (!q || job.name.toLowerCase().includes(q)) && matchesStatus(job, statusFilter)
  );
  const counts: Record<JobStatusFilter, number> = {
    all: tasks.length,
    running: tasks.filter((job) => job.status === 'running').length,
    failed: tasks.filter((job) => isFailedStatus(job.status)).length,
    scheduled: tasks.filter((job) => job.status !== 'inactive').length,
    disabled: tasks.filter((job) => job.status === 'inactive').length,
  };
  const showEmpty = tasks.length === 0 && statusFilter === 'all' && !q;

  const grouped = useMemo(() => {
    const map = new Map<string, Job[]>();
    for (const job of visible) {
      const key = job.projectId ?? '';
      const list = map.get(key) ?? [];
      list.push(job);
      map.set(key, list);
    }
    return map;
  }, [visible]);

  if (projectsError) {
    return (
      <QueryErrorState
        title="Jobs are unavailable"
        description="Could not load scheduled tasks."
        onRetry={() => {
          void refetchProjects();
          for (const query of taskQueries) void query.refetch();
        }}
      />
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <PageHeader
          title="Jobs"
          description="Schedule recurring commands on your services and track every run."
        />
        <Button asChild>
          <Link to="/jobs/new">
            <Plus className="size-4" />
            New job
          </Link>
        </Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
        </div>
      ) : showEmpty ? (
        <JobsEmptyState onCreate={() => navigate({ to: '/jobs/new' })} />
      ) : (
        <>
          <div className="relative max-w-md">
            <Search className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground" />
            <input
              type="text"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Search jobs…"
              className="w-full rounded-xl border border-border/50 bg-card py-2.5 ps-10 pe-4 text-foreground text-sm transition-all placeholder:text-muted-foreground focus:border-primary/20 focus:outline-none focus:ring-2 focus:ring-primary/20"
            />
          </div>

          {failedQueries.length > 0 && (
            <div role="alert" className="rounded-lg border p-3 text-sm">
              <p>
                Could not load jobs for{' '}
                {failedQueries.map(({ project }) => project.name).join(', ')}.
              </p>
              <button
                type="button"
                className="mt-2 underline"
                onClick={() => {
                  for (const { query } of failedQueries) void query.refetch();
                }}
              >
                Retry failed projects
              </button>
            </div>
          )}

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_340px]">
            <div className="min-w-0 space-y-8">
              {visible.length === 0 ? (
                <div className="rounded-2xl border border-border/50 bg-card px-5 py-12 text-center text-muted-foreground text-sm">
                  No jobs match this filter.
                </div>
              ) : (
                [...grouped.entries()].map(([projectId, jobs]) => (
                  <section key={projectId || 'unknown'} className="space-y-3">
                    <h2 className="font-semibold text-[13px] text-muted-foreground/70 uppercase tracking-wide">
                      {projectNames.get(projectId) ?? 'Other jobs'}
                    </h2>
                    <div className="space-y-3">
                      {jobs.map((job) => (
                        <JobCard
                          key={job.id}
                          job={job}
                          onChanged={refresh}
                          onViewOutput={setOutputJob}
                        />
                      ))}
                    </div>
                  </section>
                ))
              )}
            </div>

            <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
              <JobsOverview counts={counts} active={statusFilter} onSelect={setStatusFilter} />
            </div>
          </div>
        </>
      )}

      <JobOutputDialog job={outputJob} onClose={() => setOutputJob(null)} />
    </div>
  );
}
