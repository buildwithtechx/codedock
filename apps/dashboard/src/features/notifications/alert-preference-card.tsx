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
    <div className="rounded-xl border border-border/50 bg-card p-4 transition-colors hover:border-border/80">
      <div className="mb-3.5">
        <p className="font-medium text-foreground text-sm">{event.title}</p>
        <p className="mt-1 text-muted-foreground text-xs leading-relaxed">{event.description}</p>
      </div>

      <div className="flex flex-wrap items-center gap-x-6 gap-y-3 border-border/30 border-t pt-3">
        <div className="flex items-center gap-2.5">
          <span className="font-medium text-foreground text-xs">Notify me</span>
          <Switch
            checked={state.notifyMe}
            onCheckedChange={onNotifyToggle}
            aria-label={`Notify me: ${event.title}`}
          />
        </div>

        <div className="flex items-center gap-2.5">
          <span className="text-muted-foreground text-xs">Organization default</span>
          <Switch
            checked={state.orgEnabled}
            onCheckedChange={onOrgToggle}
            aria-label={`Organization default: ${event.title}`}
          />
        </div>

        <div className="sm:ml-auto">
          <ChannelMultiSelect value={state.channels} onChange={onChannelChange} />
        </div>
      </div>
    </div>
  );
}
