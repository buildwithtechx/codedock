import { MonitorCog, ShieldCheck } from 'lucide-react';
import { useAuthStore } from '#/stores/auth-store';
import { useGetPublicSettings } from './hooks';
import { SettingsSection } from './settings-section';

export function InstanceInfo() {
  const user = useAuthStore((state) => state.user);
  const { data: publicRes } = useGetPublicSettings();
  const publicSettings = publicRes?.data;
  const isCloud = publicSettings?.cloudMode ?? false;
  const version = publicSettings?.version || '0.1.0';

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<MonitorCog className="size-4 text-violet-500" />}
        iconBg="bg-violet-500/10"
        iconColor="text-violet-500"
        title="Instance Overview"
        description={
          isCloud ? 'Managed Codedock Cloud workspace.' : 'This self-hosted Codedock installation.'
        }
      >
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <div className="flex items-center justify-between gap-3 rounded-xl border border-border/60 bg-muted/20 p-4">
            <div className="flex min-w-0 items-center gap-3">
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-muted">
                <MonitorCog className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="min-w-0">
                <p className="truncate font-medium text-foreground text-sm">
                  {isCloud ? 'Codedock Cloud' : 'Self-hosted'}
                </p>
                <p className="text-muted-foreground text-xs">Operating mode</p>
              </div>
            </div>
            <span className="shrink-0 rounded-md bg-muted px-2 py-0.5 font-mono text-muted-foreground text-xs">
              v{version}
            </span>
          </div>

          <div className="flex items-center justify-between gap-3 rounded-xl border border-border/60 bg-muted/20 p-4">
            <div className="flex min-w-0 items-center gap-3">
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-muted">
                <ShieldCheck className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="min-w-0">
                <p className="font-medium text-foreground text-sm">Authentication</p>
                <p className="truncate text-muted-foreground text-xs">
                  {user?.email || 'Instance owner'}
                </p>
              </div>
            </div>
            <span className="inline-flex shrink-0 items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[11px] text-emerald-600 dark:text-emerald-400">
              Active
            </span>
          </div>
        </div>
      </SettingsSection>
    </div>
  );
}
