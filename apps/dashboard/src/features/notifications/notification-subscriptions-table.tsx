import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Bell, Check } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Switch } from '#/components/ui/switch';
import { settingsService } from '#/features/settings/api';
import { ChannelMultiSelect } from './channel-multi-select';
import { NOTIFICATION_EVENTS, NOTIFICATION_GROUPS } from './notification-events-meta';

type EventState = {
  channels: string[];
  orgEnabled: boolean;
  notifyMe: boolean;
};

export function NotificationSubscriptionsTable() {
  const queryClient = useQueryClient();
  const [selectedGroup, setSelectedGroup] = useState<string>('all');
  const [states, setStates] = useState<Record<string, EventState>>(() =>
    Object.fromEntries(
      NOTIFICATION_EVENTS.map((e) => [
        e.id,
        {
          channels: e.defaultChannels,
          orgEnabled: e.defaultOrgEnabled,
          notifyMe: e.defaultNotifyMe,
        },
      ])
    )
  );

  const { data: subsData } = useQuery({
    queryKey: ['notification-subscriptions'],
    queryFn: () => settingsService.getSubscriptions(),
  });

  const { data: defaultsData } = useQuery({
    queryKey: ['notification-defaults'],
    queryFn: () => settingsService.getDefaults(),
  });

  useEffect(() => {
    if (!subsData?.data && !defaultsData?.data) return;
    setStates((prev) => {
      const next = { ...prev };
      if (subsData?.data) {
        for (const sub of subsData.data) {
          if (next[sub.category]) {
            next[sub.category] = {
              ...next[sub.category],
              channels: sub.channels?.length ? sub.channels : next[sub.category].channels,
              notifyMe: sub.enabled,
            };
          }
        }
      }
      if (defaultsData?.data) {
        for (const def of defaultsData.data) {
          if (next[def.category]) {
            next[def.category] = {
              ...next[def.category],
              orgEnabled: def.enabled,
            };
          }
        }
      }
      return next;
    });
  }, [subsData, defaultsData]);

  const subMutation = useMutation({
    mutationFn: (payload: { category: string; channels: string[]; enabled: boolean }) =>
      settingsService.upsertSubscription(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notification-subscriptions'] });
    },
  });

  const defaultMutation = useMutation({
    mutationFn: (payload: { category: string; channels: string[]; enabled: boolean }) =>
      settingsService.upsertDefault(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notification-defaults'] });
    },
  });

  const filteredEvents = useMemo(() => {
    if (selectedGroup === 'all') return NOTIFICATION_EVENTS;
    return NOTIFICATION_EVENTS.filter((e) => e.group === selectedGroup);
  }, [selectedGroup]);

  const handleChannelChange = (id: string, channels: string[]) => {
    const current = states[id] || { channels: [], orgEnabled: false, notifyMe: false };
    setStates((prev) => ({
      ...prev,
      [id]: { ...current, channels },
    }));
    subMutation.mutate({ category: id, channels, enabled: current.notifyMe });
  };

  const handleOrgToggle = (id: string, orgEnabled: boolean) => {
    const current = states[id] || { channels: [], orgEnabled: false, notifyMe: false };
    setStates((prev) => ({
      ...prev,
      [id]: { ...current, orgEnabled },
    }));
    defaultMutation.mutate({ category: id, channels: current.channels, enabled: orgEnabled });
  };

  const handleNotifyToggle = (id: string) => {
    const current = states[id] || { channels: [], orgEnabled: false, notifyMe: false };
    const nextNotify = !current.notifyMe;
    setStates((prev) => ({
      ...prev,
      [id]: { ...current, notifyMe: nextNotify },
    }));
    subMutation.mutate({ category: id, channels: current.channels, enabled: nextNotify });
  };

  const groupCounts = useMemo(() => {
    const counts: Record<string, number> = { all: NOTIFICATION_EVENTS.length };
    for (const group of NOTIFICATION_GROUPS) {
      counts[group.id] = NOTIFICATION_EVENTS.filter((e) => e.group === group.id).length;
    }
    return counts;
  }, []);

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="mb-5 flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <Bell className="h-4 w-4" />
        </div>
        <div>
          <h2 className="font-semibold text-foreground text-sm">What to be notified about</h2>
          <p className="text-muted-foreground text-xs">
            Pick which events go to which channel. Each cell is one subscription.
          </p>
        </div>
      </div>

      <div className="mb-6 flex flex-wrap gap-1.5 border-border/60 border-b pb-3">
        <button
          type="button"
          onClick={() => setSelectedGroup('all')}
          className={`inline-flex items-center gap-1.5 rounded-lg px-2.5 py-1 font-medium text-xs transition-colors ${
            selectedGroup === 'all'
              ? 'bg-primary/10 text-primary'
              : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
          }`}
        >
          All
          <span className="rounded-full bg-muted px-1.5 py-0.2 font-mono text-[10px]">
            {groupCounts.all}
          </span>
        </button>
        {NOTIFICATION_GROUPS.map((group) => {
          const count = groupCounts[group.id] || 0;
          if (count === 0) return null;
          return (
            <button
              key={group.id}
              type="button"
              onClick={() => setSelectedGroup(group.id)}
              className={`inline-flex items-center gap-1.5 rounded-lg px-2.5 py-1 font-medium text-xs transition-colors ${
                selectedGroup === group.id
                  ? 'bg-primary/10 text-primary'
                  : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
              }`}
            >
              {group.label}
              <span className="rounded-full bg-muted px-1.5 py-0.2 font-mono text-[10px]">
                {count}
              </span>
            </button>
          );
        })}
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left">
          <thead>
            <tr className="border-border/60 border-b text-[11px] text-muted-foreground uppercase tracking-wider">
              <th className="pb-3 font-semibold">Event</th>
              <th className="pb-3 text-right font-semibold">Delivery Channels</th>
              <th className="pb-3 text-center font-semibold">Organization Defaults</th>
              <th className="pb-3 text-center font-semibold">Notify Me</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/40">
            {filteredEvents.map((event) => {
              const state = states[event.id] || {
                channels: event.defaultChannels,
                orgEnabled: event.defaultOrgEnabled,
                notifyMe: event.defaultNotifyMe,
              };
              return (
                <tr key={event.id} className="group hover:bg-muted/15">
                  <td className="py-3.5 pr-4">
                    <p className="font-medium text-foreground text-sm">{event.title}</p>
                    <p className="mt-0.5 max-w-xl text-muted-foreground text-xs leading-relaxed">
                      {event.description}
                    </p>
                  </td>
                  <td className="py-3.5 pr-6 text-right">
                    <ChannelMultiSelect
                      value={state.channels}
                      onChange={(channels) => handleChannelChange(event.id, channels)}
                    />
                  </td>
                  <td className="py-3.5 text-center">
                    <div className="inline-flex justify-center">
                      <Switch
                        checked={state.orgEnabled}
                        onCheckedChange={(orgEnabled) => handleOrgToggle(event.id, orgEnabled)}
                        aria-label={`Organization default for ${event.title}`}
                      />
                    </div>
                  </td>
                  <td className="py-3.5 text-center">
                    <button
                      type="button"
                      onClick={() => handleNotifyToggle(event.id)}
                      className={`inline-flex size-4.5 items-center justify-center rounded border transition-colors ${
                        state.notifyMe
                          ? 'border-primary bg-primary text-primary-foreground'
                          : 'border-border/70 hover:border-foreground/50'
                      }`}
                      aria-label={`Notify me for ${event.title}`}
                    >
                      {state.notifyMe && <Check className="size-3" />}
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </section>
  );
}
