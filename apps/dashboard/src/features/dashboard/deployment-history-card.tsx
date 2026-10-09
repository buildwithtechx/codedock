import { Link } from '@tanstack/react-router';
import { Ban, Clock, ExternalLink, GitBranch, MoreVertical, Rocket, RotateCcw } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import type { OrganizationDeployment } from '#/features/services';
import { useCancelDeployment, useRollback, useTrigger } from '#/hooks/use-deployments';

const activeStatuses = new Set(['PENDING', 'CLONING', 'PULLING', 'BUILDING']);
const successfulStatuses = new Set(['READY', 'ACTIVE', 'SUCCESS']);

function statusPill(status: string) {
  const normalized = status.toUpperCase();
  if (successfulStatuses.has(normalized))
    return 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400';
  if (normalized === 'FAILED') return 'bg-rose-500/10 text-rose-600 dark:text-rose-400';
  if (normalized === 'CANCELED' || normalized === 'REMOVED' || normalized === 'SLEPT')
    return 'bg-muted text-muted-foreground';
  return 'bg-amber-500/10 text-amber-600 dark:text-amber-400';
}

function statusLabel(status: string) {
  const normalized = status.toUpperCase();
  if (normalized === 'SUCCESS' || normalized === 'READY' || normalized === 'ACTIVE')
    return 'Deployed';
  return status
    .toLowerCase()
    .replace(/_/g, ' ')
    .replace(/\b\w/g, (letter) => letter.toUpperCase());
}

function timeAgo(value: string) {
  const diff = Date.now() - new Date(value).getTime();
  const minutes = Math.floor(diff / 60_000);
  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  return `${Math.floor(days / 30)}mo ago`;
}

function buildDuration(deployment: OrganizationDeployment) {
  if (!deployment.finishedAt) return null;
  const seconds = Math.max(
    1,
    Math.round(
      (new Date(deployment.finishedAt).getTime() - new Date(deployment.createdAt).getTime()) / 1000
    )
  );
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  return `${minutes}m ${seconds % 60}s`;
}

export function DeploymentHistoryCard({ deployment }: { deployment: OrganizationDeployment }) {
  const [menuOpen, setMenuOpen] = useState(false);
  const trigger = useTrigger();
  const rollback = useRollback();
  const cancel = useCancelDeployment();
  const normalized = deployment.status.toUpperCase();
  const isActive = activeStatuses.has(normalized);
  const canRollback = successfulStatuses.has(normalized) || normalized === 'FAILED';
  const duration = buildDuration(deployment);
  const commitHash = deployment.commitHash?.slice(0, 7);

  const redeploy = async () => {
    try {
      await trigger.mutateAsync({ serviceId: deployment.serviceId });
      toast.success('Redeploy started');
    } catch {
      toast.error('Failed to start redeploy');
    }
  };

  const rollbackTo = async () => {
    try {
      await rollback.mutateAsync({ deploymentId: deployment.id });
      toast.success('Rollback started');
    } catch {
      toast.error('Failed to start rollback');
    }
  };

  const cancelDeploy = async () => {
    try {
      await cancel.mutateAsync({ deploymentId: deployment.id });
      toast.success('Deployment canceled');
    } catch {
      toast.error('Failed to cancel deployment');
    }
  };

  return (
    <div className="group relative grid grid-cols-[auto_minmax(0,1fr)] items-start gap-x-3 gap-y-2 px-4 py-4 transition-colors first:rounded-t-2xl last:rounded-b-2xl hover:bg-muted/25 sm:flex sm:items-center sm:gap-4">
      <Link
        to="/deployments/$deploymentId"
        params={{ deploymentId: deployment.id }}
        aria-label={deployment.serviceName || deployment.projectName}
        className="absolute inset-0 z-0"
      />
      <div className="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-muted/45 transition-colors group-hover:bg-muted/65">
        <Rocket className="h-4 w-4 text-muted-foreground" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="mb-0.5 flex flex-wrap items-center gap-2">
          <p className="min-w-0 truncate font-medium text-foreground text-sm">
            {deployment.serviceName || 'Unknown app'}
          </p>
          {deployment.version != null && (
            <span className="shrink-0 rounded-md bg-muted/60 px-1.5 py-0.5 font-medium font-mono text-muted-foreground text-xs">
              v{deployment.version}
            </span>
          )}
          <span
            className={`inline-flex shrink-0 items-center rounded-full px-2 py-0.5 font-medium text-xs ${statusPill(deployment.status)}`}
          >
            {statusLabel(deployment.status)}
          </span>
        </div>
        <div className="mt-0.5 flex flex-wrap items-center gap-x-2 gap-y-1">
          <p className="min-w-0 max-w-full truncate text-muted-foreground text-xs sm:max-w-[320px]">
            {deployment.commitMessage || deployment.trigger || 'Manual deployment'}
          </p>
          <span className="text-muted-foreground/40">·</span>
          <span className="shrink-0 text-muted-foreground text-xs">{deployment.projectName}</span>
          <span className="text-muted-foreground/40">·</span>
          <span className="shrink-0 text-muted-foreground text-xs">
            {timeAgo(deployment.createdAt)}
          </span>
          {duration && (
            <>
              <span className="text-muted-foreground/40">·</span>
              <span className="flex shrink-0 items-center gap-1 text-muted-foreground text-xs">
                <Clock className="h-3 w-3" />
                {duration}
              </span>
            </>
          )}
          {deployment.branch && (
            <>
              <span className="hidden text-muted-foreground/40 sm:inline">·</span>
              <span className="hidden min-w-0 items-center gap-1 text-muted-foreground text-xs sm:flex">
                <GitBranch className="h-3 shrink-0" />
                <span className="truncate">{deployment.branch}</span>
              </span>
            </>
          )}
        </div>
      </div>
      <div
        className={`relative col-start-2 flex shrink-0 items-center gap-2 ${menuOpen ? 'z-30' : 'z-10'}`}
      >
        {commitHash && (
          <span className="inline-flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 font-mono text-muted-foreground text-xs">
            {commitHash}
          </span>
        )}
        <DropdownMenu onOpenChange={setMenuOpen}>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              aria-label={`Actions for ${deployment.serviceName}`}
              className="h-8 w-8 text-muted-foreground/50 hover:bg-muted/50 hover:text-foreground"
            >
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem asChild>
              <Link to="/deployments/$deploymentId" params={{ deploymentId: deployment.id }}>
                <ExternalLink className="h-4 w-4" />
                View details
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem onClick={redeploy} disabled={trigger.isPending}>
              <Rocket className="h-4 w-4" />
              Redeploy
            </DropdownMenuItem>
            {canRollback && (
              <DropdownMenuItem onClick={rollbackTo} disabled={rollback.isPending}>
                <RotateCcw className="h-4 w-4" />
                Roll back to this
              </DropdownMenuItem>
            )}
            {isActive && (
              <DropdownMenuItem
                variant="destructive"
                onClick={cancelDeploy}
                disabled={cancel.isPending}
              >
                <Ban className="h-4 w-4" />
                Cancel deployment
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}
