import { CheckCircle2, Download, ExternalLink, RefreshCw, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { useCheckUpdate, useDeployUpdate, useGetUpdateStatus } from '#/features/settings';
import { SettingsSection } from './settings-section';

export const UpdatesPage = () => {
  const { data, isLoading, refetch } = useGetUpdateStatus();
  const { mutateAsync: checkUpdate, isPending: checking } = useCheckUpdate();
  const { mutateAsync: deployUpdate, isPending: deploying } = useDeployUpdate();

  const info = data?.data;
  const hasUpdate = info?.hasUpdate ?? false;
  const upToDate = !isLoading && !hasUpdate;

  const handleCheck = async () => {
    try {
      await checkUpdate();
      await refetch();
      toast.success('Update check complete');
    } catch {
      toast.error('Failed to check for updates');
    }
  };

  const handleDeploy = async () => {
    try {
      await deployUpdate();
      toast.success('Update started — daemon will restart shortly');
    } catch {
      toast.error('Failed to deploy update');
    }
  };

  return (
    <div className="space-y-6">
      <SettingsSection
        collapsible
        icon={<RefreshCw className="size-4 text-primary" />}
        title="Instance Updates"
        description="Update status, release advisories, and version management for this Codedock install."
      >
        <div className="space-y-4">
          <div className="flex flex-col gap-4 rounded-xl border border-border/50 bg-background/50 p-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-3.5">
              <div
                className={`flex size-10 shrink-0 items-center justify-center rounded-xl border ${
                  upToDate
                    ? 'border-emerald-500/20 bg-emerald-500/10 text-emerald-500'
                    : 'border-primary/20 bg-primary/10 text-primary'
                }`}
              >
                {upToDate ? <CheckCircle2 className="size-5" /> : <Download className="size-5" />}
              </div>
              <div>
                <p className="font-medium text-foreground text-sm">
                  {isLoading
                    ? 'Checking for updates...'
                    : hasUpdate
                      ? `Version ${info?.latestVersion} available`
                      : 'Codedock is up to date'}
                </p>
                <p className="mt-0.5 text-muted-foreground text-xs">
                  {isLoading ? (
                    <Skeleton className="h-3.5 w-32" />
                  ) : (
                    `Current version: v${info?.currentVersion || '0.1.0'}`
                  )}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              {hasUpdate && (
                <Button
                  size="sm"
                  onClick={handleDeploy}
                  disabled={deploying}
                  className="gap-1.5 text-xs"
                >
                  <Download className="size-3.5" />
                  {deploying ? 'Updating...' : 'Update Now'}
                </Button>
              )}
              <Button
                variant="outline"
                size="sm"
                onClick={handleCheck}
                disabled={checking || deploying}
                className="gap-1.5 text-xs"
              >
                <RefreshCw className={`size-3.5 ${checking ? 'animate-spin' : ''}`} />
                {checking ? 'Checking...' : 'Check now'}
              </Button>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <div className="rounded-xl border border-border/50 bg-muted/10 p-3.5">
              <p className="text-[11px] text-muted-foreground uppercase tracking-wider">
                Release Repository
              </p>
              <p className="mt-1 font-mono text-foreground text-xs">buildwithtechx/codedock</p>
            </div>
            <div className="rounded-xl border border-border/50 bg-muted/10 p-3.5">
              <p className="text-[11px] text-muted-foreground uppercase tracking-wider">
                Latest Remote Tag
              </p>
              <p className="mt-1 font-mono text-foreground text-xs">
                {isLoading ? '...' : `v${info?.latestVersion || '0.1.0'}`}
              </p>
            </div>
          </div>

          {info?.releaseNotes && hasUpdate && (
            <div className="space-y-2 rounded-xl border border-border/50 bg-muted/10 p-4">
              <p className="font-semibold text-foreground text-xs">Release Notes</p>
              <pre className="max-h-48 overflow-y-auto whitespace-pre-wrap rounded-lg bg-background/50 p-3 font-mono text-[11px] text-muted-foreground leading-relaxed">
                {info.releaseNotes}
              </pre>
            </div>
          )}

          <div className="border-border/40 border-t pt-3">
            <a
              href="https://github.com/buildwithtechx/codedock/releases"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 text-foreground text-xs underline-offset-4 hover:underline"
            >
              <ExternalLink className="size-3.5" />
              View full changelog on GitHub
            </a>
          </div>
        </div>
      </SettingsSection>

      <SettingsSection
        collapsible
        icon={<ShieldCheck className="size-4 text-emerald-500" />}
        iconBg="bg-emerald-500/10"
        iconColor="text-emerald-500"
        title="Security & Updates Posture"
        description="All upgrade manifests are pulled strictly from authenticated release repositories."
      >
        <p className="text-muted-foreground text-xs leading-relaxed">
          Codedock adheres to a pull-only architecture. System binaries and container images are
          verified against signed tags on GitHub Releases. No telemetry or automated remote
          executions are forced upon your self-hosted instance.
        </p>
      </SettingsSection>
    </div>
  );
};
