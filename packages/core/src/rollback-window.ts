export const DEFAULT_ROLLBACK_WINDOW = 5;
export const MAX_ROLLBACK_WINDOW = 20;

/**
 * Clamps a rollback window value between 0 and maxWindow (default MAX_ROLLBACK_WINDOW).
 */
export function clampRollbackWindow(val: number, maxWindow: number = MAX_ROLLBACK_WINDOW): number {
  if (!Number.isFinite(val) || val <= 0) return 0;
  if (val > maxWindow) return maxWindow;
  return Math.trunc(val);
}

/**
 * Normalizes an unknown input into a valid integer rollback window.
 * If value is missing or invalid, returns fallback (default DEFAULT_ROLLBACK_WINDOW).
 */
export function normalizeRollbackWindow(
  value: unknown,
  fallback: number = DEFAULT_ROLLBACK_WINDOW
): number {
  const parsed =
    typeof value === 'number'
      ? value
      : typeof value === 'string' && value.trim() !== ''
        ? Number(value)
        : Number.NaN;
  if (!Number.isFinite(parsed)) return fallback;
  const whole = Math.trunc(parsed);
  if (whole <= 0) return 0;
  if (whole > MAX_ROLLBACK_WINDOW) return MAX_ROLLBACK_WINDOW;
  return whole;
}
