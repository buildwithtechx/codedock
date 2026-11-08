import type { AdvisoryFinding, ProjectSecurityResult } from '@codedock/contracts';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Switch } from '@codedock/ui/switch';
import { ExternalLink, Loader2, RefreshCw, ShieldAlert, ShieldCheck } from 'lucide-react';
import {
  useApplySecurityRemediation,
  useCheckProjectSecurity,
  useProjectSecurity,
  useSetProjectSecurityPolicy,
} from '#/hooks/use-security';

export function ProjectSecurityCard({
  projectId,
  projectName,
  initialResult,
  onRefresh,
}: {
  projectId: string;
  projectName?: string;
  initialResult?: ProjectSecurityResult;
  onRefresh?: () => void;
}) {
  const { data: projectSecurity, isLoading } = useProjectSecurity(projectId, initialResult);
  const checkMutation = useCheckProjectSecurity();
  const policyMutation = useSetProjectSecurityPolicy();
  const applyMutation = useApplySecurityRemediation();

  const data = projectSecurity || initialResult;
  const findings = data?.findings || [];
  const autoRefresh = data?.autoRefresh ?? false;

  const handleCheck = async () => {
    await checkMutation.mutateAsync(projectId);
    onRefresh?.();
  };

  const handleTogglePolicy = (checked: boolean) => {
    policyMutation.mutate({ projectId, autoRefresh: checked });
  };

  const handleApply = (finding: AdvisoryFinding) => {
    if (!finding.canRefresh) return;
    applyMutation.mutate({
      projectId,
      payload: {
        advisoryId: finding.advisoryId,
        deploymentId: data?.deploymentId,
        serviceId: finding.serviceId,
        runtimePatchId: finding.runtimePatch?.reviewId,
      },
    });
  };

  const severityBadgeVariant = (sev: string) => {
    switch (sev.toLowerCase()) {
      case 'critical':
      case 'high':
        return 'destructive';
      case 'medium':
        return 'secondary';
      default:
        return 'outline';
    }
  };

  return (
    <div className="space-y-4 rounded-xl border border-border/60 bg-card p-5">
      <div className="flex flex-wrap items-center justify-between gap-3 border-border/40 border-b pb-3">
        <div className="flex items-center gap-2.5">
          {findings.length > 0 ? (
            <ShieldAlert className="size-5 text-destructive" />
          ) : (
            <ShieldCheck className="size-5 text-emerald-500" />
          )}
          <div>
            <h3 className="font-semibold text-foreground text-sm">
              {projectName || `Project ${projectId}`}
            </h3>
            {data?.checkedAt && (
              <p className="text-muted-foreground text-xs">
                Audited: {new Date(data.checkedAt).toLocaleString()}
              </p>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          {findings.length > 0 ? (
            <Badge variant="destructive">
              {findings.length} advisory finding{findings.length > 1 ? 's' : ''}
            </Badge>
          ) : (
            <Badge variant="success">Up to date</Badge>
          )}
          <Button
            size="sm"
            variant="outline"
            disabled={checkMutation.isPending || (isLoading && !data)}
            onClick={handleCheck}
            className="h-8 gap-1.5 text-xs"
          >
            {checkMutation.isPending ? (
              <Loader2 className="size-3.5 animate-spin" />
            ) : (
              <RefreshCw className="size-3.5" />
            )}
            Check now
          </Button>
        </div>
      </div>

      <div className="flex items-center justify-between rounded-lg bg-muted/40 px-3.5 py-2.5">
        <div className="space-y-0.5">
          <p className="font-medium text-foreground text-xs">
            Auto-refresh verified security patches
          </p>
          <p className="text-[11px] text-muted-foreground">
            Automatically rebuild and redeploy when non-breaking security updates are signed
            upstream.
          </p>
        </div>
        <Switch
          checked={autoRefresh}
          disabled={policyMutation.isPending}
          onCheckedChange={handleTogglePolicy}
        />
      </div>

      {findings.length === 0 ? (
        <div className="rounded-lg border border-border/60 border-dashed p-4 text-center">
          <p className="text-muted-foreground text-xs">
            No known CVEs or runtime security advisories apply to this workload.
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {findings.map((finding: AdvisoryFinding) => (
            <div
              key={`${finding.advisoryId}-${finding.serviceName}`}
              className="space-y-2.5 rounded-lg border border-border/50 bg-background/50 p-3.5"
            >
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <Badge variant={severityBadgeVariant(finding.severity)}>
                      {finding.severity.toUpperCase()}
                    </Badge>
                    <span className="font-medium text-foreground text-xs">{finding.title}</span>
                  </div>
                  {finding.serviceName && (
                    <p className="text-[11px] text-muted-foreground">
                      Service:{' '}
                      <span className="font-mono text-foreground">{finding.serviceName}</span>
                    </p>
                  )}
                </div>
                <Badge variant="outline" className="text-[10px]">
                  {finding.confidence}
                </Badge>
              </div>

              {finding.evidence && (
                <div className="break-all rounded bg-muted/60 p-2 font-mono text-[11px] text-muted-foreground">
                  {finding.evidence}
                </div>
              )}

              {finding.runtimePatch && (
                <div className="flex items-center gap-2 rounded-lg border border-primary/20 bg-primary/5 px-2.5 py-1.5 text-xs">
                  <span className="font-medium text-primary">Runtime Patch:</span>
                  <span className="font-mono text-muted-foreground">
                    {finding.runtimePatch.currentVersion}
                  </span>
                  <span className="font-bold text-primary">→</span>
                  <span className="font-mono font-semibold text-foreground">
                    {finding.runtimePatch.targetVersion}
                  </span>
                  <span className="text-[11px] text-muted-foreground">
                    ({finding.runtimePatch.runtime})
                  </span>
                </div>
              )}

              <div className="flex flex-wrap items-center justify-between gap-2 border-border/30 border-t pt-1">
                <div className="flex items-center gap-2">
                  <a
                    href={finding.action || '#'}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-primary text-xs underline-offset-4 hover:underline"
                  >
                    Guidance
                    <ExternalLink className="size-3" />
                  </a>
                  {finding.runtimeUpgradeRequired && (
                    <span className="text-[11px] text-amber-500">
                      Major runtime upgrade required
                    </span>
                  )}
                  {finding.imageChangeRequired && (
                    <span className="text-[11px] text-amber-500">
                      Image manifest change required
                    </span>
                  )}
                </div>

                {finding.canRefresh && (
                  <Button
                    size="sm"
                    variant="default"
                    disabled={applyMutation.isPending}
                    onClick={() => handleApply(finding)}
                    className="h-7 text-xs"
                  >
                    {applyMutation.isPending ? (
                      <>
                        <Loader2 className="mr-1 size-3 animate-spin" />
                        Remediating...
                      </>
                    ) : (
                      'Refresh & Redeploy'
                    )}
                  </Button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
