import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import type { Job } from '#/features/services';
import { formatTime } from './job-format';

export function JobOutputDialog({ job, onClose }: { job: Job | null; onClose: () => void }) {
  return (
    <Dialog open={job !== null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Last output{job ? ` · ${job.name}` : ''}</DialogTitle>
          <DialogDescription>
            {job ? `Captured ${formatTime(job.lastRunAt).toLowerCase()}` : ''}
          </DialogDescription>
        </DialogHeader>
        <pre className="max-h-96 overflow-auto whitespace-pre-wrap rounded-lg bg-muted/40 p-4 font-mono text-xs leading-relaxed">
          {job?.lastOutput || 'No output recorded yet.'}
        </pre>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
