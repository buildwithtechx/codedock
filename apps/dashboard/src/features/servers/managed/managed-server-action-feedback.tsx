import { AlertCircle, ArrowUpRight } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import type { CloudWorkspaceSummary } from '#/interfaces/server';
import { cn } from '#/lib/utils';

export interface ServerResizePreview {
  before: {
    cpuCores: number;
    memoryMb: number;
    diskMb: number;
  };
  after: {
    cpuCores: number;
    memoryMb: number;
    diskMb: number;
  };
}

export interface ManagedServerActions {
  busy?: boolean;
  error?: string | null;
  preview?: ServerResizePreview | null;
  remove: () => Promise<void> | void;
  resize?: () => Promise<void> | void;
  closePreview?: () => void;
}

export interface ManagedServerActionFeedbackProps {
  server: CloudWorkspaceSummary;
  actions: ManagedServerActions;
  deleting: boolean;
  onCancelDelete: () => void;
  className?: string;
}

export function ManagedServerActionFeedback({
  server,
  actions,
  deleting,
  onCancelDelete,
  className,
}: ManagedServerActionFeedbackProps) {
  const pending = ['queued', 'running'].includes(server.operation?.status ?? '');
  const preview = actions.preview;
  const sameSize =
    preview &&
    preview.before.cpuCores === preview.after.cpuCores &&
    preview.before.memoryMb === preview.after.memoryMb &&
    preview.before.diskMb === preview.after.diskMb;

  if (!actions.error && !deleting && !preview) return null;

  return (
    <div className={cn('space-y-4 text-left', className)}>
      {actions.error && (
        <div
          role="alert"
          className="flex items-center gap-2 rounded-xl border border-destructive/30 bg-destructive/10 p-4 text-destructive text-sm"
        >
          <AlertCircle className="size-4 shrink-0" />
          <span>{actions.error}</span>
        </div>
      )}

      {deleting && (
        <section className="space-y-4 rounded-2xl border border-destructive/20 bg-card p-5">
          <h3 className="font-semibold text-base text-destructive">Delete Server</h3>
          <p className="text-muted-foreground text-sm leading-relaxed">
            Deleting this server immediately deprovisions its cloud compute. If this server has an
            active subscription, manage or cancel it in billing to stop recurring charges.
          </p>

          {server.projectCount && server.projectCount > 0 ? (
            <p className="font-medium text-amber-600 text-xs dark:text-amber-400">
              Cannot delete server while it contains {server.projectCount} active project(s). Remove
              or migrate all projects first.
            </p>
          ) : null}

          <div className="flex flex-wrap items-center gap-2 pt-2">
            <Button asChild variant="secondary" size="sm">
              <a href={`/billing/overview?workspaceId=${encodeURIComponent(server.id)}`}>
                <span>Manage subscription</span>
                <ArrowUpRight className="size-3.5" />
              </a>
            </Button>
            <Button
              variant="destructive"
              size="sm"
              disabled={actions.busy || pending || (server.projectCount ?? 0) > 0}
              onClick={() => void actions.remove()}
            >
              Confirm deletion
            </Button>
            <Button variant="outline" size="sm" disabled={actions.busy} onClick={onCancelDelete}>
              Cancel
            </Button>
          </div>
        </section>
      )}

      {preview && (
        <section className="space-y-4 rounded-2xl border border-border/50 bg-card p-5">
          <h3 className="font-semibold text-base text-foreground">Confirm Server Resize</h3>
          <div className="grid grid-cols-2 gap-4 rounded-xl bg-muted/40 p-4 text-xs">
            <div>
              <p className="font-semibold text-muted-foreground uppercase">Current</p>
              <p className="mt-1 font-medium text-foreground">
                {preview.before.cpuCores} vCPU / {Math.round(preview.before.memoryMb / 1024)} GB RAM
                / {Math.round(preview.before.diskMb / 1024)} GB Disk
              </p>
            </div>
            <div>
              <p className="font-semibold text-muted-foreground uppercase">New</p>
              <p className="mt-1 font-medium text-foreground">
                {preview.after.cpuCores} vCPU / {Math.round(preview.after.memoryMb / 1024)} GB RAM /{' '}
                {Math.round(preview.after.diskMb / 1024)} GB Disk
              </p>
            </div>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            {!sameSize && (
              <Button
                size="sm"
                disabled={actions.busy || pending}
                onClick={() => void actions.resize?.()}
              >
                Confirm resize
              </Button>
            )}
            <Button
              variant="outline"
              size="sm"
              disabled={actions.busy}
              onClick={actions.closePreview}
            >
              Cancel
            </Button>
          </div>
        </section>
      )}
    </div>
  );
}
