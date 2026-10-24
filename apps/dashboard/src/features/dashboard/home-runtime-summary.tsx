import { Link } from '@tanstack/react-router';
import {
  Activity,
  AlertTriangle,
  ArrowRight,
  CheckCircle2,
  FolderKanban,
  Rocket,
} from 'lucide-react';
import type { CanvasSummary } from '#/interfaces/projects';

export function HomeRuntimeSummary({
  projects,
  totalDeployments = 0,
  brokenCount = 0,
  issuesLoaded = true,
  isLoading,
  isUnavailable = false,
}: {
  projects: CanvasSummary[];
  totalDeployments?: number;
  brokenCount?: number;
  issuesLoaded?: boolean;
  isLoading: boolean;
  isUnavailable?: boolean;
}) {
  const busy = isLoading || isUnavailable;
  const hasIssues = issuesLoaded && brokenCount > 0;
  const statusTone = !issuesLoaded
    ? 'text-muted-foreground'
    : brokenCount === 0
      ? 'text-emerald-500'
      : 'text-amber-500';
  const statusValue = !issuesLoaded
    ? '–'
    : brokenCount === 0
      ? 'Operational'
      : brokenCount === 1
        ? '1 issue'
        : `${brokenCount} issues`;

  const statusRow = (
    <>
      <div className="flex min-w-0 items-center gap-2">
        {hasIssues ? (
          <AlertTriangle className={`size-4 shrink-0 ${statusTone}`} />
        ) : (
          <CheckCircle2 className={`size-4 shrink-0 ${statusTone}`} />
        )}
        <span className="truncate text-muted-foreground text-sm">System status</span>
      </div>
      <span className={`inline-flex shrink-0 items-center gap-1 font-medium text-sm ${statusTone}`}>
        {statusValue}
        {hasIssues && <ArrowRight className="size-3.5 rtl:rotate-180" />}
      </span>
    </>
  );

  return (
    <div className="rounded-2xl border border-border/50 bg-card p-5 shadow-xs">
      <div className="mb-4 flex items-center gap-2">
        <Activity className="size-4 text-muted-foreground" />
        <h3 className="font-semibold text-foreground text-sm">Activity</h3>
      </div>

      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <FolderKanban className="size-4 text-primary" />
            </div>
            <span className="text-muted-foreground text-sm">Projects</span>
          </div>
          <span className="font-semibold text-foreground text-lg tabular-nums">
            {busy ? '–' : projects.length}
          </span>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-500/10">
              <Rocket className="size-4 text-orange-500" />
            </div>
            <span className="text-muted-foreground text-sm">Deployments</span>
          </div>
          <span className="font-semibold text-foreground text-lg tabular-nums">
            {busy ? '–' : totalDeployments}
          </span>
        </div>

        <div className="my-2 h-px bg-border/60" />

        {hasIssues ? (
          <Link
            to="/monitoring"
            className="-mx-2 flex items-center justify-between rounded-lg px-2 py-1 transition-colors hover:bg-muted/40"
          >
            {statusRow}
          </Link>
        ) : (
          <div className="flex items-center justify-between">{statusRow}</div>
        )}
      </div>
    </div>
  );
}
