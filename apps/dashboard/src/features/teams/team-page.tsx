import { useNavigate, useRouterState } from '@tanstack/react-router';
import { AlertTriangle, Building2, ShieldCheck, Users } from 'lucide-react';
import type * as React from 'react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { AuditLogList } from '#/features/audit/audit-log-list';
import { TeamSettings } from '#/features/settings/team-settings';
import { useListOrganizations } from '#/hooks/use-organizations';
import { useOrganizationStore } from '#/stores/organization-store';
import { TeamDangerZoneTab } from './team-danger-zone-tab';
import { TeamGeneralTab } from './team-general-tab';

export type TeamTabId = 'general' | 'members' | 'audit' | 'danger';

interface TeamTabItem {
  id: TeamTabId;
  label: string;
  icon: React.ReactNode;
}

const TEAM_TABS: TeamTabItem[] = [
  { id: 'general', label: 'General', icon: <Building2 className="h-4 w-4" /> },
  { id: 'members', label: 'Members', icon: <Users className="h-4 w-4" /> },
  { id: 'audit', label: 'Audit Logs', icon: <ShieldCheck className="h-4 w-4" /> },
  { id: 'danger', label: 'Danger Zone', icon: <AlertTriangle className="h-4 w-4" /> },
];

export function TeamLayout() {
  const navigate = useNavigate();
  const search = useRouterState({
    select: (state) => state.location.search as { tab?: string },
  });

  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);
  const { data: orgs, isLoading } = useListOrganizations();
  const organization = orgs?.find((o) => o.id === activeOrgId) || orgs?.[0];

  const rawTab = search.tab as TeamTabId;
  const activeTab: TeamTabId = TEAM_TABS.some((t) => t.id === rawTab) ? rawTab : 'general';

  const setActiveTab = (tab: TeamTabId) => {
    void navigate({
      to: '/team',
      search: { tab },
      replace: true,
    });
  };

  if (isLoading) {
    return <div className="h-64 animate-pulse rounded-2xl border border-border/80 bg-card" />;
  }

  if (!organization) {
    return (
      <div className="flex min-h-64 flex-col items-center justify-center rounded-2xl border border-border border-dashed bg-card px-6 text-center">
        <Building2 className="size-8 text-muted-foreground" />
        <h3 className="mt-4 font-semibold text-base">No workspace available</h3>
        <p className="mt-1 text-muted-foreground text-sm">
          Please create a workspace to view and manage team settings.
        </p>
      </div>
    );
  }

  const content: Record<TeamTabId, React.ReactNode> = {
    general: <TeamGeneralTab organization={organization} />,
    members: <TeamSettings />,
    audit: <AuditLogList />,
    danger: <TeamDangerZoneTab organization={organization} />,
  };

  return (
    <div className="space-y-5">
      <PageHeader
        title="Team & Workspace"
        description="Manage workspace organization, team member permissions, audit logs, and security controls."
      />

      {/* Mobile Tab Scroller */}
      <div className="flex gap-1 overflow-x-auto xl:hidden">
        {TEAM_TABS.map((tab) => {
          const isActive = tab.id === activeTab;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              type="button"
              aria-current={isActive ? 'page' : undefined}
              className={`flex shrink-0 items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors ${
                isActive
                  ? 'bg-primary/12 font-medium text-foreground'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground'
              }`}
            >
              {tab.icon}
              {tab.label}
            </button>
          );
        })}
      </div>

      <PageFrame
        className="lg:grid-cols-[minmax(0,1fr)_280px]"
        rail={
          <div className="space-y-3">
            <section className="rounded-2xl border border-border/50 bg-card p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary font-bold text-sm">
                  {organization.name.charAt(0).toUpperCase()}
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-foreground text-sm">
                    {organization.name}
                  </p>
                  <p className="truncate text-muted-foreground text-xs">Team Settings</p>
                </div>
              </div>
            </section>

            <div className="space-y-1 rounded-2xl border border-border/50 bg-card p-3">
              <p className="px-3 pt-2 pb-1 font-semibold text-[10px] text-muted-foreground/80 uppercase tracking-wider">
                Team Workspace
              </p>
              {TEAM_TABS.map((tab) => {
                const isActive = tab.id === activeTab;
                return (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    type="button"
                    aria-current={isActive ? 'page' : undefined}
                    className={`flex w-full items-center gap-2.5 rounded-xl px-3 py-1.5 text-left font-medium text-[13px] transition-colors ${
                      isActive
                        ? 'bg-foreground/[0.08] text-foreground'
                        : 'text-muted-foreground hover:bg-foreground/[0.04] hover:text-foreground'
                    }`}
                  >
                    <span className={isActive ? 'text-primary' : 'text-muted-foreground'}>
                      {tab.icon}
                    </span>
                    <span>{tab.label}</span>
                  </button>
                );
              })}
            </div>
          </div>
        }
      >
        {content[activeTab]}
      </PageFrame>
    </div>
  );
}
