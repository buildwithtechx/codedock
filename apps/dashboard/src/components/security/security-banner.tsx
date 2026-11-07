import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { ArrowRight, ExternalLink, ShieldAlert, X } from 'lucide-react';
import { useState } from 'react';
import { useDismissedAdvisories, useSecurityAdvisories } from '#/hooks/use-security';
import { SecurityReviewDialog } from './security-review-dialog';

export function SecurityBanner() {
  const { data, refetch } = useSecurityAdvisories();
  const { dismissed, dismiss } = useDismissedAdvisories();
  const [reviewOpen, setReviewOpen] = useState(false);

  if (!data?.affectedProjects || data.affectedProjects.length === 0) {
    return null;
  }

  // Filter out advisories where all findings are for dismissed advisory IDs
  const activeFindings = data.affectedProjects.flatMap((p) => p.findings);
  const undismissedFinding = activeFindings.find((f) => !dismissed.includes(f.advisoryId));

  if (!undismissedFinding) {
    return null;
  }

  const activeAdvisory = data.advisories.find((a) => a.id === undismissedFinding.advisoryId);
  const affectedCount = data.affectedProjects.filter((p) =>
    p.findings.some((f) => !dismissed.includes(f.advisoryId))
  ).length;

  const isCritical =
    activeAdvisory?.severity === 'critical' ||
    undismissedFinding.severity.toLowerCase() === 'critical';

  return (
    <>
      <div
        role="alert"
        className={`mb-4 flex flex-col items-start justify-between gap-3 rounded-xl border p-4 sm:flex-row sm:items-center ${
          isCritical
            ? 'border-destructive/40 bg-destructive/10 text-destructive-foreground'
            : 'border-amber-500/40 bg-amber-500/10 text-foreground'
        }`}
      >
        <div className="flex items-start gap-3">
          <ShieldAlert
            className={`mt-0.5 size-5 shrink-0 ${
              isCritical ? 'text-destructive' : 'text-amber-500'
            }`}
          />
          <div className="space-y-0.5">
            <div className="flex flex-wrap items-center gap-2">
              <Badge
                variant={isCritical ? 'destructive' : 'secondary'}
                className="text-[10px] uppercase"
              >
                {activeAdvisory?.severity || undismissedFinding.severity}
              </Badge>
              <span className="font-semibold text-foreground text-sm">
                {activeAdvisory?.title || undismissedFinding.title}
              </span>
            </div>
            <p className="text-muted-foreground text-xs">
              {affectedCount} workload{affectedCount > 1 ? 's' : ''} affected by signed security
              advisories. Non-invasive checks detected vulnerable runtime or base images.
            </p>
          </div>
        </div>

        <div className="flex shrink-0 items-center gap-2 self-end sm:self-auto">
          <Button
            size="sm"
            variant="default"
            onClick={() => setReviewOpen(true)}
            className="h-8 gap-1.5 text-xs"
          >
            Review affected apps
            <ArrowRight className="size-3.5" />
          </Button>

          {activeAdvisory?.url && (
            <a
              href={activeAdvisory.url}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex h-8 items-center gap-1 rounded-md border border-border/50 px-2.5 text-muted-foreground text-xs hover:bg-accent hover:text-foreground"
            >
              Guidance
              <ExternalLink className="size-3" />
            </a>
          )}

          <button
            type="button"
            onClick={() => dismiss(undismissedFinding.advisoryId)}
            aria-label="Dismiss security notice"
            className="rounded-md p-1.5 text-muted-foreground hover:bg-accent hover:text-foreground"
          >
            <X className="size-4" />
          </button>
        </div>
      </div>

      <SecurityReviewDialog
        open={reviewOpen}
        onOpenChange={setReviewOpen}
        data={data}
        onRefresh={() => void refetch()}
      />
    </>
  );
}
