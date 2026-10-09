import { useQuery } from '@tanstack/react-query';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { jobsApi } from '#/features/jobs/api';
import { JobForm } from '#/features/jobs/job-form';

export const Route = createFileRoute('/_dashboard/jobs/$jobId/edit')({
  component: EditJobPage,
});

function EditJobPage() {
  const { jobId } = Route.useParams();
  const navigate = useNavigate();
  const backToDetail = () => navigate({ to: '/jobs/$jobId', params: { jobId } });

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['scheduled-task', jobId],
    queryFn: () => jobsApi.get(jobId),
  });
  const job = data?.data;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={backToDetail}
          className="flex size-8 items-center justify-center rounded-lg transition-colors hover:bg-muted"
          aria-label="Back to job"
        >
          <ArrowLeft className="size-4 text-muted-foreground" />
        </button>
        <h1 className="truncate font-medium text-2xl text-foreground/80">
          Edit job{job ? ` · ${job.name}` : ''}
        </h1>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
        </div>
      ) : isError || !job ? (
        <QueryErrorState
          title="Job is unavailable"
          description="Could not load this job for editing."
          onRetry={() => refetch()}
        />
      ) : (
        <JobForm
          job={job}
          onCancel={backToDetail}
          onSaved={() => {
            toast.success(`Job ${job.name} saved`);
            backToDetail();
          }}
        />
      )}
    </div>
  );
}
