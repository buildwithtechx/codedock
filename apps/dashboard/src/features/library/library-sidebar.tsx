import { useNavigate } from '@tanstack/react-router';
import { BookOpen, GitBranch, Globe, Lock, Zap } from 'lucide-react';
import { GIT_PROVIDERS, type ProviderConnection } from './types';

export function LibrarySidebar({
  connections,
  counts,
}: {
  connections: ProviderConnection[];
  counts: { total: number; publicCount: number; privateCount: number };
}) {
  const navigate = useNavigate();
  const byProvider = new Map(connections.map((item) => [item.provider, item]));
  const displayProviders = GIT_PROVIDERS;

  return (
    <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
      <div className="rounded-2xl border border-border/50 bg-card p-5">
        <h3 className="mb-4 font-semibold text-foreground text-sm">Connection</h3>
        <div className="space-y-2.5">
          {displayProviders.map((provider) => {
            const status = byProvider.get(provider.id);
            const connected = status?.connected ?? false;
            return (
              <div key={provider.id} className="flex items-center justify-between gap-3">
                <div className="flex min-w-0 items-center gap-2.5">
                  <div
                    className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-lg ${
                      connected ? 'bg-emerald-500/10' : 'bg-muted/60'
                    }`}
                  >
                    <img
                      src={provider.icon}
                      alt=""
                      aria-hidden="true"
                      className="size-4 object-contain"
                    />
                  </div>
                  <div className="min-w-0">
                    <p className="truncate font-medium text-foreground text-sm">{provider.name}</p>
                    <p className="truncate text-muted-foreground text-xs">
                      {connected ? (status?.accountName ?? 'Connected') : 'Not connected'}
                    </p>
                  </div>
                </div>
                <span
                  className={`inline-flex shrink-0 items-center gap-1.5 rounded-full px-2 py-0.5 font-medium text-[10px] ${
                    connected
                      ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                      : 'bg-muted/60 text-muted-foreground'
                  }`}
                >
                  <span
                    className={`size-1.5 rounded-full ${
                      connected ? 'bg-emerald-500' : 'bg-muted-foreground/40'
                    }`}
                  />
                  {connected ? 'Connected' : '—'}
                </span>
              </div>
            );
          })}
        </div>
        <button
          type="button"
          onClick={() => void navigate({ to: '/settings', search: { tab: 'git' } as never })}
          className="mt-3 block font-medium text-muted-foreground text-xs transition-colors hover:text-foreground hover:underline"
        >
          Manage git providers in settings
        </button>
      </div>

      {counts.total > 0 && (
        <div className="rounded-2xl border border-border/50 bg-card p-5">
          <div className="mb-4 flex items-center gap-2">
            <BookOpen className="size-4 text-muted-foreground" />
            <h3 className="font-semibold text-foreground text-sm">Overview</h3>
          </div>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                  <GitBranch className="size-4 text-primary" />
                </div>
                <span className="text-muted-foreground text-sm">Total</span>
              </div>
              <span className="font-semibold text-foreground text-lg">{counts.total}</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500/10">
                  <Globe className="size-4 text-blue-500" />
                </div>
                <span className="text-muted-foreground text-sm">Public</span>
              </div>
              <span className="font-semibold text-foreground text-lg">{counts.publicCount}</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-500/10">
                  <Lock className="size-4 text-orange-500" />
                </div>
                <span className="text-muted-foreground text-sm">Private</span>
              </div>
              <span className="font-semibold text-foreground text-lg">{counts.privateCount}</span>
            </div>
          </div>
        </div>
      )}

      <div className="rounded-2xl border border-primary/10 bg-gradient-to-br from-primary/5 via-primary/3 to-transparent p-5">
        <div className="mb-3 flex items-center gap-2">
          <Zap className="size-4 text-primary" />
          <h3 className="font-semibold text-foreground text-sm">Quick tip</h3>
        </div>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Select any repository to deploy it instantly. Configure automatic deployments on every
          push.
        </p>
      </div>
    </div>
  );
}
