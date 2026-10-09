import { Link } from '@tanstack/react-router';
import { ArrowRight, Check, Eye, Loader2 } from 'lucide-react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { cn } from '#/lib/utils';
import type { AttentionIssue } from './interfaces';
import {
  ACTION_LABELS,
  formatRelativeTime,
  issueTargetHref,
  kindMeta,
  SEVERITY_TONE,
  toDisplaySeverity,
} from './issue-meta';
import { KindIcon } from './kind-icon';

export function IssueRow({
  issue,
  busy,
  onAcknowledge,
  onResolve,
  onAct,
}: {
  issue: AttentionIssue;
  busy: boolean;
  onAcknowledge: (issue: AttentionIssue) => void;
  onResolve: (issue: AttentionIssue) => void;
  onAct: (issue: AttentionIssue) => void;
}) {
  const severity = toDisplaySeverity(issue.severity);
  const tone = SEVERITY_TONE[severity];
  const meta = kindMeta(issue.kind);
  const resolved = issue.status === 'resolved';
  const acked = issue.status === 'acked';
  const actLabel = issue.action ? (ACTION_LABELS[issue.action] ?? 'Fix') : null;
  const seen = formatRelativeTime(issue.firstSeen);
  const metaLine = [
    issue.occurrences > 1 ? `${issue.occurrences} occurrences` : null,
    seen ? `since ${seen}` : null,
  ]
    .filter(Boolean)
    .join(' · ');

  return (
    <li className="flex items-start gap-3 px-5 py-4">
      <span
        className={cn(
          'flex size-9 shrink-0 items-center justify-center rounded-xl',
          tone.soft,
          tone.text
        )}
      >
        <KindIcon kind={issue.kind} className="size-[18px]" />
      </span>
      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
          <Link
            to={issueTargetHref(issue)}
            className="truncate font-medium text-sm hover:text-primary hover:underline"
          >
            {issue.title}
          </Link>
          <span className="shrink-0 text-muted-foreground/60 text-xs">{meta.title}</span>
          {acked && (
            <Badge variant="secondary" className="shrink-0">
              Acked
            </Badge>
          )}
        </div>
        {issue.detail && (
          <p
            className={cn(
              'mt-1 line-clamp-2 text-[13px] leading-snug',
              severity === 'advisory' ? 'text-muted-foreground' : tone.text
            )}
            title={issue.detail}
          >
            {issue.detail}
          </p>
        )}
        {issue.remediation && !resolved && (
          <p className="mt-1 truncate text-muted-foreground text-xs" title={issue.remediation}>
            {issue.remediation}
          </p>
        )}
        {metaLine && <p className="mt-0.5 truncate text-muted-foreground/70 text-xs">{metaLine}</p>}
      </div>
      <div className="flex shrink-0 flex-wrap items-center justify-end gap-1.5">
        {resolved ? (
          <Button variant="ghost" size="sm" asChild>
            <Link to={issueTargetHref(issue)}>
              View
              <ArrowRight className="size-3.5" />
            </Link>
          </Button>
        ) : (
          <>
            {actLabel && issue.action && (
              <Button
                size="sm"
                disabled={busy}
                onClick={() => onAct(issue)}
                title={issue.remediation || actLabel}
              >
                {busy && <Loader2 className="size-3.5 animate-spin" />}
                {actLabel}
              </Button>
            )}
            {!actLabel && (
              <Button variant="ghost" size="sm" asChild>
                <Link to={issueTargetHref(issue)}>
                  View
                  <ArrowRight className="size-3.5" />
                </Link>
              </Button>
            )}
            {!acked && (
              <Button
                variant="ghost"
                size="sm"
                disabled={busy}
                onClick={() => onAcknowledge(issue)}
                title="Acknowledge this issue"
              >
                <Eye className="size-3.5" />
                Ack
              </Button>
            )}
            <Button
              variant="outline"
              size="sm"
              disabled={busy}
              onClick={() => onResolve(issue)}
              title="Mark resolved"
            >
              <Check className="size-3.5" />
              Resolve
            </Button>
          </>
        )}
      </div>
    </li>
  );
}
