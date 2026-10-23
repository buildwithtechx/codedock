import { useNavigate, useRouterState } from '@tanstack/react-router';
import { Blocks, Cloud, KeyRound, ScrollText, ShieldCheck } from 'lucide-react';
import type * as React from 'react';
import { PageFrame } from '#/components/layout/page-frame';
import { PageHeader } from '#/components/layout/page-header';
import { ApiKeysList } from '#/features/users/api-keys-list';
import { CloudInitTab } from './cloud-init-tab';
import { CloudTokensTab } from './cloud-tokens-tab';
import { IntegrationTokensTab } from './integration-tokens-tab';
import { PrivateKeysTab } from './private-keys-tab';

export type TokensTabId =
  | 'private-keys'
  | 'cloud-tokens'
  | 'integration-tokens'
  | 'cloud-init'
  | 'api-tokens';

interface TokensTabItem {
  id: TokensTabId;
  label: string;
  icon: React.ReactNode;
}

const TOKEN_TABS: TokensTabItem[] = [
  { id: 'private-keys', label: 'Private Keys', icon: <KeyRound className="h-4 w-4" /> },
  { id: 'cloud-tokens', label: 'Cloud Tokens', icon: <Cloud className="h-4 w-4" /> },
  { id: 'integration-tokens', label: 'Integration Tokens', icon: <Blocks className="h-4 w-4" /> },
  { id: 'cloud-init', label: 'Cloud-Init Scripts', icon: <ScrollText className="h-4 w-4" /> },
  { id: 'api-tokens', label: 'API Tokens', icon: <ShieldCheck className="h-4 w-4" /> },
];

export function TokensLayout() {
  const navigate = useNavigate();
  const search = useRouterState({
    select: (state) => state.location.search as { tab?: string },
  });

  const rawTab = search.tab as TokensTabId;
  const activeTab: TokensTabId = TOKEN_TABS.some((t) => t.id === rawTab) ? rawTab : 'private-keys';

  const setActiveTab = (tab: TokensTabId) => {
    void navigate({
      to: '/tokens',
      search: { tab },
      replace: true,
    });
  };

  const content: Record<TokensTabId, React.ReactNode> = {
    'private-keys': <PrivateKeysTab />,
    'cloud-tokens': <CloudTokensTab />,
    'integration-tokens': <IntegrationTokensTab />,
    'cloud-init': <CloudInitTab />,
    'api-tokens': <ApiKeysList />,
  };

  return (
    <div className="space-y-5">
      <PageHeader
        title="Keys & Tokens"
        description="Manage SSH keys, cloud provider credentials, integration tokens, bootstrap scripts, and API access."
      />

      {/* Mobile Tab Scroller */}
      <div className="flex gap-1 overflow-x-auto xl:hidden">
        {TOKEN_TABS.map((tab) => {
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
                <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
                  <KeyRound className="size-4" />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-foreground text-sm">Security & Access</p>
                  <p className="truncate text-muted-foreground text-xs">Credentials & Auth</p>
                </div>
              </div>
            </section>

            <div className="space-y-1 rounded-2xl border border-border/50 bg-card p-3">
              <p className="px-3 pt-2 pb-1 font-semibold text-[10px] text-muted-foreground/80 uppercase tracking-wider">
                Keys & Credentials
              </p>
              {TOKEN_TABS.map((tab) => {
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
