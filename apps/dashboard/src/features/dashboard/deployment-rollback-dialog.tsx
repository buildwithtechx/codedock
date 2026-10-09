import { RotateCcw } from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';

export function DeploymentRollbackDialog({
  open,
  releaseLabel,
  busy,
  onOpenChange,
  onConfirm,
}: {
  open: boolean;
  releaseLabel: string;
  busy: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void;
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <div className="flex items-start gap-3">
            <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10">
              <RotateCcw className="h-5 w-5 text-primary" />
            </span>
            <div className="space-y-1">
              <DialogTitle>Roll back to this release?</DialogTitle>
              <DialogDescription>{releaseLabel}</DialogDescription>
            </div>
          </div>
        </DialogHeader>
        <div className="space-y-3 text-sm">
          <p className="text-muted-foreground leading-6">
            Codedock will restore this release as the active version of the service. A rollback
            replays this release&apos;s frozen configuration over today&apos;s values, so settings
            changed since this release may be reverted.
          </p>
          <p className="text-muted-foreground leading-6">
            Attached data is left unchanged. This starts a new deployment and cannot be undone.
          </p>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={busy}>
            Cancel
          </Button>
          <Button onClick={onConfirm} disabled={busy}>
            {busy ? 'Rolling back...' : 'Roll back'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
