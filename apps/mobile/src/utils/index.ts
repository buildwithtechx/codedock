export function cleanUrl(url: string): string {
  return url.trim().replace(/\/+$/, '');
}

export function formatStatus(status?: string): {
  label: string;
  isSuccess: boolean;
  isWarning: boolean;
  isDanger: boolean;
} {
  const normalized = (status || 'stopped').toLowerCase();
  return {
    label: normalized,
    isSuccess: normalized === 'running' || normalized === 'succeeded' || normalized === 'online',
    isWarning:
      normalized === 'building' ||
      normalized === 'deploying' ||
      normalized === 'pending' ||
      normalized === 'starting',
    isDanger:
      normalized === 'failed' ||
      normalized === 'error' ||
      normalized === 'cancelled' ||
      normalized === 'unreachable',
  };
}

export function formatDuration(seconds?: number): string {
  if (!seconds || seconds <= 0) return '';
  if (seconds < 60) return `${seconds}s`;
  const mins = Math.floor(seconds / 60);
  const remaining = seconds % 60;
  return remaining > 0 ? `${mins}m ${remaining}s` : `${mins}m`;
}

export function truncate(str: string, maxLen: number = 32): string {
  if (!str) return '';
  return str.length > maxLen ? `${str.slice(0, maxLen - 1)}…` : str;
}
