/**
 * Formats a byte count into a human-readable string (B, KB, MB, GB, TB, PB).
 */
export function formatBytes(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 B';
  if (!Number.isFinite(bytes) || bytes < 0) return '—';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const idx = Math.min(i, sizes.length - 1);
  return `${parseFloat((bytes / Math.pow(k, idx)).toFixed(dm))} ${sizes[idx]}`;
}

/**
 * Formats seconds into a concise human-readable duration (e.g. "2m 15s", "45s", "1h 3m").
 */
export function formatDuration(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds <= 0) return '0s';
  const sec = Math.floor(seconds);
  if (sec < 60) return `${sec}s`;
  const mins = Math.floor(sec / 60);
  const remainingSec = sec % 60;
  if (mins < 60) {
    return remainingSec > 0 ? `${mins}m ${remainingSec}s` : `${mins}m`;
  }
  const hours = Math.floor(mins / 60);
  const remainingMins = mins % 60;
  return remainingMins > 0 ? `${hours}h ${remainingMins}m` : `${hours}h`;
}
