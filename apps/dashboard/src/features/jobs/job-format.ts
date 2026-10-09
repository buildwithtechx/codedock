import type { JobStatus } from '#/features/services';

export function formatTime(iso?: string | null): string {
  if (!iso) return 'Never';
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return 'Never';
  return date.toLocaleString();
}

export function formatDuration(ms?: number | null): string {
  if (ms == null) return '';
  return ms < 1000 ? `${ms}ms` : `${(ms / 1000).toFixed(1)}s`;
}

export function statusTone(status: JobStatus): string {
  if (status === 'failed' || status === 'error') return 'text-destructive';
  if (status === 'completed' || status === 'active') return 'text-success';
  if (status === 'running') return 'text-warning';
  return 'text-muted-foreground';
}

export function statusDot(status: JobStatus): string {
  if (status === 'failed' || status === 'error') return 'bg-destructive';
  if (status === 'completed' || status === 'active') return 'bg-success';
  if (status === 'running') return 'bg-warning';
  return 'bg-muted-foreground/40';
}

export function statusLabel(status: JobStatus): string {
  return status.charAt(0).toUpperCase() + status.slice(1);
}

export function isFailedStatus(status: JobStatus): boolean {
  return status === 'failed' || status === 'error';
}

export function relativeTime(iso?: string | null): string {
  if (!iso) return 'Never';
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return 'Never';
  const diff = Date.now() - date.getTime();
  const min = Math.floor(diff / 60_000);
  if (min < 1) return 'just now';
  if (min < 60) return `${min}m ago`;
  const hr = Math.floor(min / 60);
  if (hr < 24) return `${hr}h ago`;
  const days = Math.floor(hr / 24);
  if (days < 7) return `${days}d ago`;
  return date.toLocaleDateString();
}
