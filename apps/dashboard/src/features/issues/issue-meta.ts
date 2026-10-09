import type {
  AttentionIssue,
  AttentionSeverity,
  DisplaySeverity,
  SeverityFilter,
} from './interfaces';

export const KIND_ORDER = [
  'deployment-failed',
  'service-unhealthy',
  'backup-failed',
  'migration-failed',
  'quota-pressure',
];

export const KIND_META: Record<string, { title: string; subtitle: string }> = {
  'deployment-failed': {
    title: 'Deployments',
    subtitle: 'Releases that failed in the last 24 hours',
  },
  'service-unhealthy': {
    title: 'Services',
    subtitle: 'Workloads that stopped running',
  },
  'backup-failed': {
    title: 'Backups',
    subtitle: 'Backup runs that did not complete',
  },
  'migration-failed': {
    title: 'Migrations',
    subtitle: 'Migration runs waiting on a fix',
  },
  'quota-pressure': {
    title: 'Capacity',
    subtitle: 'Managed quota nearing its limit',
  },
};

export const UNKNOWN_KIND_META = {
  title: 'Other',
  subtitle: 'Issues from unrecognized sources',
};

export function kindMeta(kind: string): { title: string; subtitle: string } {
  return KIND_META[kind] ?? UNKNOWN_KIND_META;
}

export const SEVERITY_FILTERS: SeverityFilter[] = ['all', 'outage', 'action_required', 'advisory'];

export const SEVERITY_LABELS: Record<DisplaySeverity, string> = {
  outage: 'Outage',
  action_required: 'Action required',
  advisory: 'Advisory',
};

export function toDisplaySeverity(severity: AttentionSeverity): DisplaySeverity {
  if (severity === 'critical') return 'outage';
  if (severity === 'warning') return 'action_required';
  return 'advisory';
}

export const SEVERITY_RANK: Record<DisplaySeverity, number> = {
  outage: 0,
  action_required: 1,
  advisory: 2,
};

export const SEVERITY_TONE: Record<
  DisplaySeverity,
  { text: string; soft: string; border: string; dot: string; bar: string }
> = {
  outage: {
    text: 'text-destructive',
    soft: 'bg-destructive/10',
    border: 'border-destructive/25',
    dot: 'bg-destructive',
    bar: 'bg-destructive',
  },
  action_required: {
    text: 'text-amber-500',
    soft: 'bg-amber-500/10',
    border: 'border-amber-500/25',
    dot: 'bg-amber-500',
    bar: 'bg-amber-500',
  },
  advisory: {
    text: 'text-muted-foreground',
    soft: 'bg-muted',
    border: 'border-border/60',
    dot: 'bg-muted-foreground/50',
    bar: 'bg-muted-foreground/40',
  },
};

export const ACTION_LABELS: Record<string, string> = {
  'restart-service': 'Restart',
  redeploy: 'Redeploy',
  'resume-migration': 'Resume',
  'retry-backup': 'Retry backup',
};

export function issueTargetHref(issue: AttentionIssue): string {
  if (issue.projectId) return `/projects/${issue.projectId}`;
  if (issue.kind === 'backup-failed') return '/backups';
  if (issue.kind === 'deployment-failed') return '/deployments';
  return '/organizations';
}

export function formatRelativeTime(value: string): string {
  const time = new Date(value).getTime();
  if (Number.isNaN(time)) return '';
  const diff = Date.now() - time;
  if (diff < 0) return 'just now';
  const minutes = Math.floor(diff / 60_000);
  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  return new Date(time).toLocaleDateString();
}
