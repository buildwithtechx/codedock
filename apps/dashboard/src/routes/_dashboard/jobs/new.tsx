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
      <div className="mb-6 flex items-center gap-3">
        <button
          type="button"
          onClick={() => navigate({ to: '/jobs' })}
          className="flex size-8 items-center justify-center rounded-lg transition-colors hover:bg-muted"
          aria-label="Back to jobs"
        >
          <ArrowLeft className="h-4 w-4 text-muted-foreground rtl:rotate-180" />
        </button>
        <div>
          <h1 className="font-medium text-2xl text-foreground/80 tracking-[-0.2px]">New job</h1>
          <p className="mt-1 text-muted-foreground/70 text-sm">
            Run a command on one of your servers — on a schedule, once, or on an event — with retry,
            secrets, live logs and notifications.
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
