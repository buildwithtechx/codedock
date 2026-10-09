import { useQuery, useQueryClient } from '@tanstack/react-query';
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router';
import {
  ArrowLeft,
  CalendarClock,
  Loader2,
  Play,
  Server,
  SlidersHorizontal,
  Terminal,
  Trash2,
} from 'lucide-react';
import type { ReactNode } from 'react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '#/components/ui/empty';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Switch } from '#/components/ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { jobsApi } from '#/features/jobs/api';
import { formatTime, statusDot, statusLabel } from '#/features/jobs/job-format';

export const Route = createFileRoute('/_dashboard/jobs/$jobId')({
  component: JobDetailsPage,
});

function InfoRow({ icon, label, value }: { icon: ReactNode; label: string; value: ReactNode }) {
  return (
    <div className="flex items-start gap-3 rounded-xl border border-border/50 bg-card px-4 py-3">
      <span className="mt-0.5 text-muted-foreground/60">{icon}</span>
      <span className="w-28 shrink-0 text-[12px] text-muted-foreground/70">{label}</span>
      <span className="min-w-0 flex-1 break-words text-[13px] text-foreground">{value}</span>
    </div>
  );
}

function JobDetailsPage() {
  const { jobId } = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [busy, setBusy] = useState(false);

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['scheduled-task', jobId],
    queryFn: () => jobsApi.get(jobId),
  });
  const job = data?.data;
  const enabled = job?.status !== 'inactive';

  const refresh = () => {
    void refetch();
    void queryClient.invalidateQueries({ queryKey: ['scheduled-tasks'] });
  };

  const handleRun = async () => {
    if (!job || busy) return;
    setBusy(true);
    try {
      await jobsApi.trigger(job.id);
      toast.success(`Job ${job.name} started`);
      refresh();
    } catch {
      toast.error(`Could not run ${job.name}`);
    } finally {
      setBusy(false);
    }
  };

  const handleToggle = async (next: boolean) => {
    if (!job || busy) return;
    setBusy(true);
    try {
      await jobsApi.update(job.id, { status: next ? 'active' : 'inactive' });
      toast.success(next ? `Job ${job.name} enabled` : `Job ${job.name} disabled`);
      refresh();
    } catch {
      toast.error(`Could not ${next ? 'enable' : 'disable'} ${job.name}`);
    } finally {
      setBusy(false);
    }
  };

  const handleDelete = async () => {
    if (!job) return;
    try {
      await jobsApi.remove(job.id);
      toast.success(`Job ${job.name} deleted`);
      void queryClient.invalidateQueries({ queryKey: ['scheduled-tasks'] });
      navigate({ to: '/jobs' });
    } catch {
      toast.error(`Could not delete ${job.name}`);
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="size-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError) {
    return (
      <QueryErrorState
        title="Job is unavailable"
        description="Could not load this job."
        onRetry={() => refetch()}
      />
    );
  }

  if (!job) {
    return (
      <Empty>
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <Server />
          </EmptyMedia>
          <EmptyTitle>Job not found</EmptyTitle>
          <EmptyDescription>
            This job no longer exists or you don&apos;t have access to it.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <Button asChild variant="outline">
            <Link to="/jobs">Back to Jobs</Link>
          </Button>
        </EmptyContent>
      </Empty>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <button
          type="button"
          onClick={() => navigate({ to: '/jobs' })}
          className="flex size-8 items-center justify-center rounded-lg transition-colors hover:bg-muted"
          aria-label="Back to jobs"
        >
          <ArrowLeft className="size-4 text-muted-foreground" />
        </button>
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <h1 className="truncate font-medium text-2xl text-foreground/80">{job.name}</h1>
            <span className="inline-flex shrink-0 items-center gap-1.5 rounded-full bg-muted px-2 py-0.5 font-medium text-[11px]">
              <span className={`size-1.5 rounded-full ${statusDot(job.status)}`} />
              {statusLabel(job.status)}
            </span>
          </div>
          <p className="mt-1 truncate font-mono text-[12px] text-muted-foreground/60">{job.id}</p>
        </div>
        <div className="flex shrink-0 items-center gap-1.5">
          <Button variant="outline" size="sm" onClick={() => void handleRun()} disabled={busy}>
            {busy ? <Loader2 className="size-4 animate-spin" /> : <Play className="size-4" />}
            Run
          </Button>
          <Button variant="outline" size="sm" asChild>
            <Link to="/jobs/$jobId/edit" params={{ jobId: job.id }}>
              <SlidersHorizontal className="size-4" />
              Edit
            </Link>
          </Button>
          <Switch
            checked={enabled}
            onCheckedChange={(next) => void handleToggle(next)}
            disabled={busy}
            aria-label={enabled ? `Disable ${job.name}` : `Enable ${job.name}`}
          />
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setConfirmDelete(true)}
            aria-label={`Delete ${job.name}`}
            className="text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
          >
            <Trash2 className="size-4" />
          </Button>
        </div>
      </div>

      <Tabs defaultValue="overview">
        <TabsList variant="line">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="output">Last output</TabsTrigger>
        </TabsList>
        <TabsContent value="overview" className="mt-6">
          <div className="max-w-3xl space-y-3">
            <InfoRow
              icon={<CalendarClock className="size-4" />}
              label="Schedule"
              value={<code className="font-mono text-[12px]">{job.schedule}</code>}
            />
            <InfoRow
              icon={<CalendarClock className="size-4" />}
              label="Last run"
              value={enabled ? formatTime(job.lastRunAt) : 'Not scheduled'}
            />
            <InfoRow
              icon={<Server className="size-4" />}
              label="Service"
              value={<code className="font-mono text-[12px]">{job.serviceId}</code>}
            />
            <div className="rounded-xl border border-border/50 bg-card p-4">
              <p className="mb-2 flex items-center gap-2 font-medium text-[12px] text-muted-foreground">
                <Terminal className="size-4" />
                Command
              </p>
              <pre className="overflow-x-auto rounded-lg bg-muted/40 px-3 py-2.5 font-mono text-[12px]">
                {job.command}
              </pre>
            </div>
            <InfoRow
              icon={<CalendarClock className="size-4" />}
              label="Created"
              value={formatTime(job.createdAt)}
            />
          </div>
        </TabsContent>
        <TabsContent value="output" className="mt-6">
          <div className="max-w-3xl">
            {job.lastOutput ? (
              <pre className="max-h-96 overflow-auto whitespace-pre-wrap rounded-xl border border-border/50 bg-card p-4 font-mono text-xs leading-relaxed">
                {job.lastOutput}
              </pre>
            ) : (
              <div className="rounded-2xl border border-border/60 border-dashed py-12 text-center text-muted-foreground/70 text-sm">
                No output recorded yet. Run the job to capture its first output.
              </div>
            )}
          </div>
        </TabsContent>
      </Tabs>

      <Dialog open={confirmDelete} onOpenChange={setConfirmDelete}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {job.name}?</DialogTitle>
            <DialogDescription>
              This removes the scheduled job. Past execution state is kept for reference.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setConfirmDelete(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={() => void handleDelete()}>
              Delete job
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
