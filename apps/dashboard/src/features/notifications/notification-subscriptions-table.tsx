import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useMemo, useState } from 'react';
import { useListOrganizations } from '#/features/organizations';
import { settingsService } from '#/features/settings/api';
import { useOrganizationStore } from '#/stores/organization-store';
import { AlertPreferenceCard, type AlertPreferenceState } from './alert-preference-card';
import { NOTIFICATION_EVENTS, NOTIFICATION_GROUPS } from './notification-events-meta';

export function NotificationSubscriptionsTable() {
  const queryClient = useQueryClient();
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const { data: organizations } = useListOrganizations();
  const organizationId = activeOrganizationId || organizations?.[0]?.id;

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
    queryKey: ['notification-subscriptions', organizationId],
    queryFn: () => settingsService.getSubscriptions(organizationId),
  });

  const { data: defaultsData } = useQuery({
    queryKey: ['notification-defaults', organizationId],
    queryFn: () => settingsService.getDefaults(organizationId),
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
    mutationFn: (payload: {
      category: string;
      channels: string[];
      enabled: boolean;
      orgId?: string;
    }) => settingsService.upsertSubscription(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notification-subscriptions'] });
    },
  });

  const defaultMutation = useMutation({
    mutationFn: (payload: {
      category: string;
      channels: string[];
      enabled: boolean;
      orgId?: string;
    }) => settingsService.upsertDefault(payload),
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
    subMutation.mutate({
      category: id,
      channels,
      enabled: current.notifyMe,
      orgId: organizationId,
    });
  };

  const handleOrgToggle = (id: string, orgEnabled: boolean) => {
    const current = states[id] || { channels: [], orgEnabled: false, notifyMe: false };
    setStates((prev) => ({
      ...prev,
      [id]: { ...current, orgEnabled },
    }));
    defaultMutation.mutate({
      category: id,
      channels: current.channels,
      enabled: orgEnabled,
      orgId: organizationId,
    });
  };

  const handleNotifyToggle = (id: string) => {
    const current = states[id] || { channels: [], orgEnabled: false, notifyMe: false };
    const nextNotify = !current.notifyMe;
    setStates((prev) => ({
      ...prev,
      [id]: { ...current, notifyMe: nextNotify },
    }));
    subMutation.mutate({
      category: id,
      channels: current.channels,
      enabled: nextNotify,
      orgId: organizationId,
    });
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
    <section
      aria-label="Alert preferences"
      className="space-y-5 rounded-2xl border border-border/50 bg-card p-5"
    >
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="min-w-0">
          <h2 className="font-medium text-base text-foreground">Alert preferences</h2>
          <p className="mt-1 text-muted-foreground text-sm leading-relaxed">
            Organization defaults apply when a member has not made an explicit choice for an alert.
            Members can override them with their own preferences.
          </p>
        </div>
      </div>

      <div className="scrollbar-hide flex items-center gap-1 overflow-x-auto border-border/50 border-b">
        <button
          type="button"
          onClick={() => setSelectedGroup('all')}
          className={`relative inline-flex min-w-0 shrink-0 items-center gap-1.5 whitespace-nowrap px-2.5 py-2 font-medium text-xs transition-colors sm:px-3.5 sm:text-sm ${
            selectedGroup === 'all'
              ? 'text-foreground'
              : 'text-muted-foreground hover:text-foreground/80'
          }`}
        >
          All
          <span
            className={`rounded-full px-1.5 py-0.5 font-medium text-[11px] tabular-nums sm:text-xs ${
              selectedGroup === 'all'
                ? 'bg-muted text-foreground'
                : 'bg-muted/60 text-muted-foreground'
            }`}
          >
            {groupCounts.all}
          </span>
          {selectedGroup === 'all' && (
            <span className="absolute start-2.5 end-2.5 bottom-0 h-0.5 rounded-full bg-primary sm:start-3.5 sm:end-3.5" />
          )}
        </button>

        {NOTIFICATION_GROUPS.map((group) => {
          const count = groupCounts[group.id] || 0;
          if (count === 0) return null;
          const active = selectedGroup === group.id;
          return (
            <button
              key={group.id}
              type="button"
              onClick={() => setSelectedGroup(group.id)}
              className={`relative inline-flex min-w-0 shrink-0 items-center gap-1.5 whitespace-nowrap px-2.5 py-2 font-medium text-xs transition-colors sm:px-3.5 sm:text-sm ${
                active ? 'text-foreground' : 'text-muted-foreground hover:text-foreground/80'
              }`}
            >
              {group.label}
              <span
                className={`rounded-full px-1.5 py-0.5 font-medium text-[11px] tabular-nums sm:text-xs ${
                  active ? 'bg-muted text-foreground' : 'bg-muted/60 text-muted-foreground'
                }`}
              >
                {count}
              </span>
              {active && (
                <span className="absolute start-2.5 end-2.5 bottom-0 h-0.5 rounded-full bg-primary sm:start-3.5 sm:end-3.5" />
              )}
            </button>
          );
        })}
      </div>

      <div className="space-y-5">
        {sections.map((section) => (
          <section key={section.id} className="space-y-3">
            {section.label && (
              <h3 className="font-medium text-foreground text-sm">{section.label}</h3>
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
          </section>
        ))}
      </div>
    </section>
  );
}
