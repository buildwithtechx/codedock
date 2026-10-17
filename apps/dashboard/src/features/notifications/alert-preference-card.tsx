import { Switch } from '#/components/ui/switch';
import { ChannelMultiSelect } from './channel-multi-select';
import type { NotificationEventDef } from './notification-events-meta';

export interface AlertPreferenceState {
  channels: string[];
  orgEnabled: boolean;
  notifyMe: boolean;
}

interface AlertPreferenceCardProps {
  event: NotificationEventDef;
  state: AlertPreferenceState;
  onChannelChange: (channels: string[]) => void;
  onOrgToggle: (enabled: boolean) => void;
  onNotifyToggle: () => void;
}

export function AlertPreferenceCard({
  event,
  state,
  onChannelChange,
  onOrgToggle,
  onNotifyToggle,
}: AlertPreferenceCardProps) {
  return (
    <div className="space-y-4 rounded-xl border border-border/40 bg-card p-4">
      <div>
        <p className="font-medium text-foreground text-sm">{event.title}</p>
        <p className="mt-1 text-muted-foreground text-xs leading-relaxed">{event.description}</p>
      </div>

      <div className="flex flex-wrap items-center gap-x-6 gap-y-4">
        <div className="flex items-center gap-3">
          <span className="text-foreground text-sm">Notify me</span>
          <Switch
            checked={state.notifyMe}
            onCheckedChange={onNotifyToggle}
            aria-label={`Notify me: ${event.title}`}
          />
        </div>

        <div className="flex items-center gap-3">
          <span className="text-muted-foreground text-sm">Organization defaults</span>
          <Switch
            checked={state.orgEnabled}
            onCheckedChange={onOrgToggle}
            aria-label={`Organization defaults: ${event.title}`}
          />
        </div>

        <ChannelMultiSelect value={state.channels} onChange={onChannelChange} />
      </div>
    </div>
  );
}
