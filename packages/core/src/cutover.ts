export type CutoverConfig = {
  ports: string[];
  volumes: string[];
  publicPort?: number;
  slug: string;
  serviceName: string;
  deploymentId: string;
};

export const CUTOVER_SHARED_VOLUMES_ENV = 'CODEDOCK_CUTOVER_SHARED_VOLUMES';

export function hasHostPortBindings(ports: string[]): boolean {
  return ports.some((p) => {
    const mapping = p.includes('/') ? p.slice(0, p.lastIndexOf('/')) : p;
    return mapping.includes(':');
  });
}

/** Temporary name for candidate while predecessor still serves. */
export function newcomerName(cfg: { deploymentId: string; serviceName: string }): string {
  const id = cfg.deploymentId.replace(/[^a-zA-Z0-9]/g, '').slice(0, 16);
  return `codedock-cutover-${id}-${cfg.serviceName}`;
}

export function stableName(cfg: { slug: string; serviceName: string }): string {
  return `codedock-${cfg.slug}-${cfg.serviceName}`;
}

function normalizeDest(dest: string): string {
  const cleaned = dest.replace(/\/+/g, '/').replace(/\/$/, '');
  return cleaned.startsWith('/') ? cleaned : `/${cleaned}`;
}

export interface CutoverMount {
  Type?: string;
  Source?: string;
  Name?: string;
  Destination: string;
  RW?: boolean;
}

export type CutoverDecision =
  | { mode: 'overlap' }
  | {
      mode: 'replace';
      reason:
        | 'no-incumbent'
        | 'host-ports'
        | 'shared-namespaces'
        | 'devices'
        | 'storage-unverified'
        | 'shared-rw-volume'
        | 'resources-unconfigured'
        | 'capacity-unavailable'
        | 'insufficient-capacity';
      clashes?: string[];
    };

export interface CutoverInput {
  incumbentRunning: boolean;
  ports: string[];
  networkMode?: string;
  pidMode?: string;
  hasDevices?: boolean;
  volumes?: string[];
  incumbentMounts?: readonly CutoverMount[];
  sharedOptIn?: ReadonlySet<string>;
  resources?: { cpuCores?: number; memoryMb?: number };
  capacity?: { availableCpuCores: number; availableMemoryMb: number };
}

export function parseVolumeDestination(spec: string): { target: string; isReadOnly: boolean } {
  const parts = spec.split(':');
  if (parts.length === 1) {
    return { target: normalizeDest(parts[0]), isReadOnly: false };
  }
  const target = parts[1];
  const mode = parts.length > 2 ? parts[2] : '';
  const isReadOnly = mode.includes('ro');
  return { target: normalizeDest(target), isReadOnly };
}

export function decideCutover(opts: CutoverInput): CutoverDecision {
  if (!opts.incumbentRunning) {
    return { mode: 'replace', reason: 'no-incumbent' };
  }
  if (hasHostPortBindings(opts.ports)) {
    return { mode: 'replace', reason: 'host-ports' };
  }
  if (opts.networkMode || opts.pidMode) {
    return { mode: 'replace', reason: 'shared-namespaces' };
  }
  if (opts.hasDevices) {
    return { mode: 'replace', reason: 'devices' };
  }

  if (opts.volumes && opts.incumbentMounts) {
    const liveRw = opts.incumbentMounts
      .filter((m) => m.RW !== false)
      .map((m) => normalizeDest(m.Destination));
    const clashes: string[] = [];
    for (const spec of opts.volumes) {
      const parsed = parseVolumeDestination(spec);
      if (parsed.isReadOnly) continue;
      if (liveRw.includes(parsed.target)) {
        if (!opts.sharedOptIn || !opts.sharedOptIn.has(parsed.target)) {
          clashes.push(parsed.target);
        }
      }
    }
    if (clashes.length > 0) {
      return { mode: 'replace', reason: 'shared-rw-volume', clashes };
    }
  }

  if (opts.resources && opts.capacity) {
    const cpu = opts.resources.cpuCores ?? 0;
    const mem = opts.resources.memoryMb ?? 0;
    if (cpu > 0 && cpu > opts.capacity.availableCpuCores) {
      return { mode: 'replace', reason: 'insufficient-capacity' };
    }
    if (mem > 0 && mem > opts.capacity.availableMemoryMb) {
      return { mode: 'replace', reason: 'insufficient-capacity' };
    }
  }

  return { mode: 'overlap' };
}
