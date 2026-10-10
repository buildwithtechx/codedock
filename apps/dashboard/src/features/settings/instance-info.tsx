import { MonitorCog, ShieldCheck } from 'lucide-react';
import { useAuthStore } from '#/stores/auth-store';
import { CloudConnectionCard } from './cloud-connection-card';
import { useGetPublicSettings } from './hooks';
import { PlatformModeCard } from './platform-mode-card';
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
          <div className="flex items-center gap-3 rounded-xl border border-border/60 bg-muted/20 p-4">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted">
              <MonitorCog className="h-4 w-4 text-muted-foreground" />
            </div>
            <div>
              <p className="font-medium text-foreground text-sm">
                {isCloud ? 'Codedock Cloud' : 'Self-hosted'}
              </p>
              <p className="text-muted-foreground text-xs">Operating mode · v{version}</p>
            </div>
          </div>
          <div className="flex items-center gap-3 rounded-xl border border-border/60 bg-muted/20 p-4">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted">
              <ShieldCheck className="h-4 w-4 text-muted-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-foreground text-sm">Authentication</p>
              <p className="truncate text-muted-foreground text-xs">
                {user?.email || 'Instance owner'}
              </p>
            </div>
          </div>
        </div>
      </SettingsSection>

      <PlatformModeCard initialMode={isCloud ? 'cloud' : 'self-hosted'} />
      {!isCloud && <CloudConnectionCard isCloudMode={false} />}
    </div>
  );
}
