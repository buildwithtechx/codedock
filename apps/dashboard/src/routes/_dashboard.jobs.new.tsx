import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { ArrowLeft } from 'lucide-react';
import { JobForm } from '#/features/jobs/job-form';

export const Route = createFileRoute('/_dashboard/jobs/new')({
  component: NewJobPage,
});

function NewJobPage() {
  const navigate = useNavigate();

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={() => navigate({ to: '/jobs' })}
          className="flex size-8 items-center justify-center rounded-lg transition-colors hover:bg-muted"
          aria-label="Back to jobs"
        >
          <ArrowLeft className="h-4 w-4 text-muted-foreground" />
        </button>
        <div>
          <h1 className="font-medium text-2xl text-foreground/80">Create job</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Run a command on a service on a recurring schedule.
          </p>
        </div>
      </div>

      <JobForm
        onCancel={() => navigate({ to: '/jobs' })}
        onSaved={(job) => navigate({ to: '/jobs/$jobId', params: { jobId: job.id } })}
      />
    </div>
  );
}
