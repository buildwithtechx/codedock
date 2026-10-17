import { GitPullRequest, Loader2, ShieldCheck } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { checksService, type ProjectChecksConfig } from '#/services/checks';

interface ProjectChecksSettingsProps {
  projectId: string;
  initialConfig?: ProjectChecksConfig | null;
  onSaved?: (config: ProjectChecksConfig) => void;
}

export function ProjectChecksSettings({
  projectId,
  initialConfig,
  onSaved,
}: ProjectChecksSettingsProps) {
  const [config, setConfig] = useState<ProjectChecksConfig>(
    initialConfig ?? {
      projectId,
      enabled: true,
      reportCommitStatus: true,
      requirePassingChecks: false,
      targetBranches: ['main', 'master'],
      contextName: 'codedock/deploy',
      updatedAt: new Date().toISOString(),
    }
  );

  const [reportCommitStatus, setReportCommitStatus] = useState(config.reportCommitStatus);
  const [requirePassingChecks, setRequirePassingChecks] = useState(config.requirePassingChecks);
  const [targetBranches, setTargetBranches] = useState(config.targetBranches.join(', '));
  const [contextName, setContextName] = useState(config.contextName);
  const [isSaving, setIsSaving] = useState(false);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      const branches = targetBranches
        .split(',')
        .map((b) => b.trim())
        .filter(Boolean);

      const updated = await checksService.updateChecks(projectId, {
        enabled: true,
        reportCommitStatus,
        requirePassingChecks,
        targetBranches: branches.length > 0 ? branches : ['main', 'master'],
        contextName: contextName.trim() || 'codedock/deploy',
      });

      setConfig(updated);
      onSaved?.(updated);
      toast.success('GitHub checks settings saved successfully');
    } catch {
      toast.error('Failed to update checks settings');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Card className="border-border/70 shadow-sm">
      <CardHeader>
        <div className="flex items-center gap-2">
          <GitPullRequest className="size-4 text-primary" />
          <CardTitle className="font-semibold text-base">
            GitHub Checks & Deployment Admission
          </CardTitle>
        </div>
        <CardDescription className="text-xs">
          Report commit and pull request build statuses back to GitHub and gate production releases.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleSave} className="space-y-4">
          <div className="flex items-center justify-between rounded-lg border border-border p-3">
            <div className="space-y-0.5">
              <Label htmlFor="report-status" className="font-medium text-sm">
                Report Commit Status to GitHub
              </Label>
              <p className="text-muted-foreground text-xs">
                Posts pending, success, and failure statuses directly on commits and pull requests.
              </p>
            </div>
            <Switch
              id="report-status"
              checked={reportCommitStatus}
              onCheckedChange={setReportCommitStatus}
            />
          </div>

          <div className="flex items-center justify-between rounded-lg border border-border p-3">
            <div className="space-y-0.5">
              <Label htmlFor="require-checks" className="font-medium text-sm">
                Require Passing Checks Before Deploy
              </Label>
              <p className="text-muted-foreground text-xs">
                Prevents promotion to production if external checks or tests are failing or pending.
              </p>
            </div>
            <Switch
              id="require-checks"
              checked={requirePassingChecks}
              onCheckedChange={setRequirePassingChecks}
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="target-branches">Monitored Branches (comma-separated)</Label>
            <Input
              id="target-branches"
              value={targetBranches}
              onChange={(e) => setTargetBranches(e.target.value)}
              placeholder="e.g. main, master, staging"
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="context-name">GitHub Status Context</Label>
            <Input
              id="context-name"
              value={contextName}
              onChange={(e) => setContextName(e.target.value)}
              placeholder="e.g. codedock/deploy"
            />
          </div>

          <div className="flex justify-end pt-2">
            <Button type="submit" size="sm" disabled={isSaving} className="gap-1.5">
              {isSaving ? (
                <Loader2 className="size-3.5 animate-spin" />
              ) : (
                <ShieldCheck className="size-3.5" />
              )}
              <span>{isSaving ? 'Saving...' : 'Save Checks Configuration'}</span>
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
