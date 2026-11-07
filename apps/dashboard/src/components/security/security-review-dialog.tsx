import type { AdvisoriesResponse } from '@codedock/contracts';
import { Badge } from '@codedock/ui/badge';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { ExternalLink, ShieldCheck } from 'lucide-react';
import { ProjectSecurityCard } from './project-security-card';

export function SecurityReviewDialog({
  open,
  onOpenChange,
  data,
  onRefresh,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  data?: AdvisoriesResponse;
  onRefresh?: () => void;
}) {
  const affectedProjects = data?.affectedProjects || [];
  const advisories = data?.advisories || [];

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[85vh] max-w-3xl overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <span>Runtime Security Advisories & Audit</span>
            {data?.stale && (
              <Badge variant="outline" className="border-amber-500/30 text-amber-500">
                Feed Stale
              </Badge>
            )}
          </DialogTitle>
          <DialogDescription>
            Signed upstream advisories verified cryptographically. Non-invasive checks evaluate
            container digests, stack versions, and safe probe outputs.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6 pt-2">
          {data?.checkedAt && (
            <div className="flex items-center justify-between rounded-lg bg-muted/30 px-3 py-2 text-muted-foreground text-xs">
              <span>Feed Revision: #{data.revision}</span>
              <span>Last Checked: {new Date(data.checkedAt).toLocaleString()}</span>
            </div>
          )}

          <div className="space-y-4">
            <h4 className="font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              Affected Workloads ({affectedProjects.length})
            </h4>

            {affectedProjects.length === 0 ? (
              <div className="flex flex-col items-center justify-center rounded-xl border border-border/70 border-dashed p-6 text-center">
                <ShieldCheck className="mb-2 size-8 text-emerald-500" />
                <p className="font-medium text-foreground text-sm">All projects are secure</p>
                <p className="mt-1 text-muted-foreground text-xs">
                  No active workloads match known vulnerabilities in the signed advisory feed.
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {affectedProjects.map((proj) => (
                  <ProjectSecurityCard
                    key={proj.projectId}
                    projectId={proj.projectId}
                    initialResult={proj}
                    onRefresh={onRefresh}
                  />
                ))}
              </div>
            )}
          </div>

          {advisories.length > 0 && (
            <div className="space-y-3 border-border/50 border-t pt-4">
              <h4 className="font-semibold text-muted-foreground text-xs uppercase tracking-wider">
                Active Signed Advisories ({advisories.length})
              </h4>
              <div className="space-y-2">
                {advisories.map((adv) => (
                  <div
                    key={adv.id}
                    className="flex flex-col justify-between gap-2 rounded-lg bg-muted/40 p-3 text-xs sm:flex-row sm:items-center"
                  >
                    <div className="space-y-0.5">
                      <div className="flex items-center gap-2">
                        <Badge
                          variant={
                            adv.severity === 'critical' || adv.severity === 'high'
                              ? 'destructive'
                              : 'secondary'
                          }
                          className="text-[10px]"
                        >
                          {adv.severity.toUpperCase()}
                        </Badge>
                        <span className="font-medium text-foreground">{adv.title}</span>
                      </div>
                      <p className="line-clamp-2 text-[11px] text-muted-foreground">
                        {adv.message}
                      </p>
                    </div>
                    {adv.url && (
                      <a
                        href={adv.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex shrink-0 items-center gap-1 text-primary underline-offset-4 hover:underline"
                      >
                        Details
                        <ExternalLink className="size-3" />
                      </a>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
