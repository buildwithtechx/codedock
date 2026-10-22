/**
 * Configuration change only. Docker creates the new volume during normal deployment,
 * while the previous volume remains retained in the host volume pool.
 */
export function emptyVolumePlan(volumes: readonly string[], selected: string) {
  if (!volumes.includes(selected)) return null;

  const parts = selected.split(':');
  const host = parts[0] || '';
  const target = parts.length > 1 ? parts[1] : host;
  const isReadOnly = parts.length > 2 && parts[2].includes('ro');

  if (isReadOnly || !target.startsWith('/')) {
    return null;
  }

  const freshSource = `codedock-fresh-${crypto.randomUUID()}`;
  const mode = parts.length > 2 ? `:${parts[2]}` : '';
  const newSpec = `${freshSource}:${target}${mode}`;

  return {
    previous: host,
    source: freshSource,
    target,
    volumes: volumes.map((v) => (v === selected ? newSpec : v)),
  };
}

export type EmptyVolumePlan = NonNullable<ReturnType<typeof emptyVolumePlan>>;
