import { Link } from '@tanstack/react-router';
import { Check, Clock3, FileText, Loader2, Pencil, Play, Trash2, X } from 'lucide-react';
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
import { Switch } from '#/components/ui/switch';
import type { Job } from '#/features/services';
import { jobsApi } from './api';
import { formatTime, isFailedStatus, relativeTime, statusDot, statusLabel } from './job-format';

export function JobCard({
  job,
  onChanged,
  onViewOutput,
}: {
  job: Job;
  onChanged: () => void;
  onViewOutput: (job: Job) => void;
}) {
  const [busy, setBusy] = useState(false);
  const [editingSchedule, setEditingSchedule] = useState(false);
  const [scheduleDraft, setScheduleDraft] = useState('');
  const [confirmDelete, setConfirmDelete] = useState(false);
  const enabled = job.status !== 'inactive';

  const handleRun = async () => {
    if (busy) return;
    setBusy(true);
    try {
      await jobsApi.trigger(job.id);
      toast.success(`Job ${job.name} started`);
      onChanged();
    } catch {
      toast.error(`Could not run ${job.name}`);
    } finally {
      setBusy(false);
    }
  };

  const handleToggle = async (next: boolean) => {
    if (busy) return;
    setBusy(true);
    try {
      await jobsApi.update(job.id, { status: next ? 'active' : 'inactive' });
      toast.success(next ? `Job ${job.name} enabled` : `Job ${job.name} disabled`);
      onChanged();
    } catch {
      toast.error(`Could not ${next ? 'enable' : 'disable'} ${job.name}`);
    } finally {
      setBusy(false);
    }
  };

  const handleSaveSchedule = async () => {
    if (busy || !scheduleDraft.trim()) return;
    setBusy(true);
    try {
      await jobsApi.update(job.id, { schedule: scheduleDraft.trim() });
      toast.success(`Schedule saved for ${job.name}`);
      setEditingSchedule(false);
      onChanged();
    } catch {
      toast.error(`Could not save schedule for ${job.name}`);
    } finally {
      setBusy(false);
    }
  };

  const handleDelete = async () => {
    try {
      await jobsApi.remove(job.id);
      toast.success(`Job ${job.name} deleted`);
      setConfirmDelete(false);
      onChanged();
    } catch {
      toast.error(`Could not delete ${job.name}`);
    }
  };

  return (
    <div className="rounded-2xl border border-border/60 bg-card px-5 py-4">
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <Clock3 className="size-4 shrink-0 text-muted-foreground/70" />
            <Link
              to="/jobs/$jobId"
              params={{ jobId: job.id }}
              className="truncate font-medium text-[14px] text-foreground hover:text-primary"
            >
              {job.name}
            </Link>
          </div>
          <div className="mt-1.5 flex items-center gap-1.5">
            <span className={`size-1.5 rounded-full ${statusDot(job.status)}`} />
            <span className="font-medium text-[12px]">{statusLabel(job.status)}</span>
            <span className="text-[12px] text-muted-foreground/60">
              · last run {relativeTime(job.lastRunAt)}
            </span>
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          {job.lastOutput && (
            <Button variant="ghost" size="sm" onClick={() => onViewOutput(job)}>
              <FileText className="size-3.5" />
              Output
            </Button>
          )}
          <Button variant="ghost" size="sm" onClick={handleRun} disabled={busy}>
            {busy ? <Loader2 className="size-3.5 animate-spin" /> : <Play className="size-3.5" />}
            Run
          </Button>
          <Switch
            checked={enabled}
            onCheckedChange={(next) => void handleToggle(next)}
            disabled={busy}
            aria-label={enabled ? `Disable ${job.name}` : `Enable ${job.name}`}
          />
          <Button variant="ghost" size="sm" asChild>
            <Link to="/jobs/$jobId/edit" params={{ jobId: job.id }} aria-label={`Edit ${job.name}`}>
              <Pencil className="size-3.5" />
            </Link>
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setConfirmDelete(true)}
            aria-label={`Delete ${job.name}`}
            className="text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
          >
            <Trash2 className="size-3.5" />
          </Button>
        </div>
      </div>

      <div className="mt-3 grid grid-cols-1 gap-x-8 gap-y-2 border-border/40 border-t pt-3 sm:grid-cols-2">
        <div className="flex items-center gap-2">
          <span className="w-20 shrink-0 text-[12px] text-muted-foreground/70">Schedule</span>
          {editingSchedule ? (
            <div className="flex items-center gap-1.5">
              <input
                value={scheduleDraft}
                onChange={(event) => setScheduleDraft(event.target.value)}
                spellCheck={false}
                className="w-40 rounded-md border border-border/60 bg-background px-2 py-1 font-mono text-[12px] text-foreground outline-none focus:border-primary"
              />
              <button
                type="button"
                onClick={() => void handleSaveSchedule()}
                disabled={busy}
                className="text-success disabled:opacity-50"
                aria-label="Save schedule"
              >
                <Check className="size-4" />
              </button>
              <button
                type="button"
                onClick={() => setEditingSchedule(false)}
                disabled={busy}
                className="text-muted-foreground hover:text-foreground disabled:opacity-50"
                aria-label="Cancel editing schedule"
              >
                <X className="size-4" />
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <code className="font-mono text-[12px] text-foreground">{job.schedule}</code>
              <button
                type="button"
                onClick={() => {
                  setEditingSchedule(true);
                  setScheduleDraft(job.schedule);
                }}
                className="text-muted-foreground/60 hover:text-foreground"
                aria-label={`Edit schedule for ${job.name}`}
              >
                <Pencil className="size-3" />
              </button>
            </div>
          )}
        </div>
        <div className="flex items-center gap-2">
          <span className="w-20 shrink-0 text-[12px] text-muted-foreground/70">Last run</span>
          <span className="text-[12px] text-foreground">
            {enabled ? formatTime(job.lastRunAt) : 'Not scheduled'}
          </span>
        </div>
      </div>

      <p
        className="mt-2 truncate font-mono text-[12px] text-muted-foreground/70"
        title={job.command}
      >
        {job.command}
      </p>

      {isFailedStatus(job.status) && job.lastOutput && (
        <p className="mt-2 truncate text-[12px] text-destructive" title={job.lastOutput}>
          {job.lastOutput}
        </p>
      )}

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
