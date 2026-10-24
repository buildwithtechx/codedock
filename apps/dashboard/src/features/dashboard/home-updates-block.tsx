import { Link } from '@tanstack/react-router';
import { AlertTriangle, ArrowRight, ArrowUpCircle, X } from 'lucide-react';
import { useMemo, useState } from 'react';
import { openInstanceUpgradeModal } from '#/features/settings/instance-upgrade-modal';
import { useListAttentionIssues } from '#/hooks/use-issues';
import { useGetUpdateStatus } from '#/hooks/use-settings';
import type { AttentionIssue } from '#/interfaces/issues';

export interface HomeAttentionState {
  broken: AttentionIssue[];
  behind: AttentionIssue[];
  hasInstanceUpdate: boolean;
  currentVersion: string;
  latestVersion: string;
  loaded: boolean;
  showBroken: boolean;
  showBehind: boolean;
  cards: number;
  hide: (card: 'broken' | 'behind') => void;
}

export function useHomeAttention(): HomeAttentionState {
  const issuesQuery = useListAttentionIssues();
  const updateQuery = useGetUpdateStatus();
  const [hidden, setHidden] = useState<Array<'broken' | 'behind'>>([]);

  const rawIssues = issuesQuery.data?.data?.issues ?? [];
  const openIssues = useMemo(
    () => rawIssues.filter((issue) => issue.status !== 'resolved'),
    [rawIssues]
  );

  const broken = useMemo(
    () => openIssues.filter((issue) => issue.severity !== 'info'),
    [openIssues]
  );
  const behind = useMemo(
    () => openIssues.filter((issue) => issue.severity === 'info'),
    [openIssues]
  );

  const updateInfo = updateQuery.data?.data;
  const hasInstanceUpdate = Boolean(updateInfo?.hasUpdate);
  const currentVersion = (updateInfo?.currentVersion || '0.1.0').replace(/^v/, '');
  const latestVersion = (updateInfo?.latestVersion || currentVersion).replace(/^v/, '');

  const showBroken = broken.length > 0 && !hidden.includes('broken');
  const showBehind = (behind.length > 0 || hasInstanceUpdate) && !hidden.includes('behind');

  return {
    broken,
    behind,
    hasInstanceUpdate,
    currentVersion,
    latestVersion,
    loaded: !issuesQuery.isLoading,
    showBroken,
    showBehind,
    cards: (showBroken ? 1 : 0) + (showBehind ? 1 : 0),
    hide: (card) => setHidden((prev) => (prev.includes(card) ? prev : [...prev, card])),
  };
}

export function HomeUpdatesBlock({ attention }: { attention: HomeAttentionState }) {
  if (attention.cards === 0) return null;

  const shownBroken = attention.broken.slice(0, 3);
  const shownBehind = attention.behind.slice(0, 3);
  const totalUpdates = attention.behind.length + (attention.hasInstanceUpdate ? 1 : 0);

  return (
    <div className="space-y-3">
      {attention.showBroken && (
        <div className="overflow-hidden rounded-2xl border border-border/50 bg-card shadow-xs">
          <div className="flex items-start gap-3 border-border/50 border-b px-5 py-4">
            <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-destructive/10">
              <AlertTriangle className="size-4 text-destructive" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <h3 className="truncate font-semibold text-[14px] text-destructive">
                  Needs attention
                </h3>
                <span className="ms-auto shrink-0 rounded-full bg-muted px-1.5 py-px font-medium text-[11px] text-muted-foreground tabular-nums">
                  {attention.broken.length}
                </span>
              </div>
              <p className="mt-0.5 text-[12px] text-muted-foreground leading-snug">
                Services or infrastructure requiring remediation
              </p>
            </div>
            <button
              type="button"
              onClick={() => attention.hide('broken')}
              aria-label="Hide attention card"
              className="-me-2 -mt-1 inline-flex size-7 items-center justify-center rounded-lg text-muted-foreground/70 transition-colors hover:bg-muted hover:text-foreground"
            >
              <X className="size-3.5" />
            </button>
          </div>
          <ul className="divide-y divide-border/50 px-5 py-3.5">
            {shownBroken.map((issue) => (
              <li
                key={issue.id}
                className="flex items-center justify-between gap-2.5 py-2.5 first:pt-0 last:pb-0"
              >
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-[13px] text-foreground">{issue.title}</p>
                  <p className="truncate text-[11px] text-muted-foreground">{issue.subject}</p>
                </div>
                <Link
                  to="/monitoring"
                  className="inline-flex h-7 shrink-0 items-center gap-1 rounded-lg bg-destructive px-2.5 font-medium text-[12px] text-destructive-foreground transition-colors hover:bg-destructive/90"
                >
                  Inspect
                </Link>
              </li>
            ))}
          </ul>
          <div className="border-border/50 border-t px-5 py-2.5">
            <Link
              to="/monitoring"
              className="inline-flex items-center gap-1 font-medium text-[12px] text-muted-foreground transition-colors hover:text-foreground"
            >
              View all
              <ArrowRight className="size-3.5 rtl:rotate-180" />
            </Link>
          </div>
        </div>
      )}

      {attention.showBehind && (
        <div className="overflow-hidden rounded-2xl border border-border/50 bg-card shadow-xs">
          <div className="flex items-start gap-3 border-border/50 border-b px-5 py-4">
            <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-amber-500/10">
              <ArrowUpCircle className="size-4 text-amber-500" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <h3 className="truncate font-semibold text-[14px] text-amber-500">
                  Updates available
                </h3>
                <span className="ms-auto shrink-0 rounded-full bg-muted px-1.5 py-px font-medium text-[11px] text-muted-foreground tabular-nums">
                  {totalUpdates}
                </span>
              </div>
              <p className="mt-0.5 text-[12px] text-muted-foreground leading-snug">
                Newer releases ready to apply
              </p>
            </div>
            <button
              type="button"
              onClick={() => attention.hide('behind')}
              aria-label="Hide updates card"
              className="-me-2 -mt-1 inline-flex size-7 items-center justify-center rounded-lg text-muted-foreground/70 transition-colors hover:bg-muted hover:text-foreground"
            >
              <X className="size-3.5" />
            </button>
          </div>
          <ul className="divide-y divide-border/50 px-5 py-3.5">
            {attention.hasInstanceUpdate && (
              <li className="flex items-center justify-between gap-2.5 py-2.5 first:pt-0 last:pb-0">
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-[13px] text-foreground">
                    Codedock Control Plane
                  </p>
                  <p className="font-mono text-[11px] text-muted-foreground">
                    v{attention.currentVersion}{' '}
                    <span className="text-amber-500">→ v{attention.latestVersion}</span>
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => openInstanceUpgradeModal()}
                  className="inline-flex h-7 shrink-0 items-center gap-1.5 rounded-lg border border-border px-2.5 font-medium text-[12px] text-muted-foreground transition-colors hover:bg-muted/50 hover:text-foreground"
                >
                  Update
                </button>
              </li>
            )}
            {shownBehind.map((issue) => (
              <li
                key={issue.id}
                className="flex items-center justify-between gap-2.5 py-2.5 first:pt-0 last:pb-0"
              >
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-[13px] text-foreground">{issue.title}</p>
                  <p className="truncate text-[11px] text-muted-foreground">{issue.subject}</p>
                </div>
                <Link
                  to="/monitoring"
                  className="inline-flex h-7 shrink-0 items-center gap-1.5 rounded-lg border border-border px-2.5 font-medium text-[12px] text-muted-foreground transition-colors hover:bg-muted/50 hover:text-foreground"
                >
                  Update
                </Link>
              </li>
            ))}
          </ul>
          <div className="border-border/50 border-t px-5 py-2.5">
            <Link
              to="/monitoring"
              className="inline-flex items-center gap-1 font-medium text-[12px] text-muted-foreground transition-colors hover:text-foreground"
            >
              View all
              <ArrowRight className="size-3.5 rtl:rotate-180" />
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}
