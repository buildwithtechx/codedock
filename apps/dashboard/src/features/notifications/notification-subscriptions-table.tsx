import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Bell } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { settingsService } from '#/features/settings/api';
import { AlertPreferenceCard, type AlertPreferenceState } from './alert-preference-card';
import { NOTIFICATION_EVENTS, NOTIFICATION_GROUPS } from './notification-events-meta';

export function NotificationSubscriptionsTable() {
  const queryClient = useQueryClient();
  const [selectedGroup, setSelectedGroup] = useState<string>('all');
  const [states, setStates] = useState<Record<string, AlertPreferenceState>>(() =>
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

  const sections = useMemo(() => {
    if (selectedGroup !== 'all') {
      const groupDef = NOTIFICATION_GROUPS.find((g) => g.id === selectedGroup);
      return [
        {
          id: selectedGroup,
          label: null,
          events: NOTIFICATION_EVENTS.filter((e) => e.group === selectedGroup),
        },
      ];
    }
    return NOTIFICATION_GROUPS.map((g) => ({
      id: g.id,
      label: g.label,
      events: NOTIFICATION_EVENTS.filter((e) => e.group === g.id),
    })).filter((s) => s.events.length > 0);
  }, [selectedGroup]);

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="mb-5 flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <Bell className="h-4 w-4" />
        </div>
        <div>
          <h2 className="font-semibold text-foreground text-sm">Alert preferences</h2>
          <p className="mt-0.5 text-muted-foreground text-xs">
            Organization defaults apply when a member has not made an explicit choice for an alert. Members can override them with their own preferences.
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

      <div className="space-y-6">
        {sections.map((section) => (
          <div key={section.id} className="space-y-3">
            {section.label && (
              <h3 className="font-semibold text-foreground text-sm tracking-tight">
                {section.label}
              </h3>
            )}
            <div className="space-y-3">
              {section.events.map((event) => {
                const state = states[event.id] || {
                  channels: event.defaultChannels,
                  orgEnabled: event.defaultOrgEnabled,
                  notifyMe: event.defaultNotifyMe,
                };
                return (
                  <AlertPreferenceCard
                    key={event.id}
                    event={event}
                    state={state}
                    onChannelChange={(channels) => handleChannelChange(event.id, channels)}
                    onOrgToggle={(orgEnabled) => handleOrgToggle(event.id, orgEnabled)}
                    onNotifyToggle={() => handleNotifyToggle(event.id)}
                  />
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
