import type { ServerGPUInfo } from '#/interfaces/server';

export function serverHasGpu(
  server: { gpu?: ServerGPUInfo | null } | null | undefined
): boolean {
  return server?.gpu?.available === true;
}

export function sortForGpu<T extends { gpu?: ServerGPUInfo | null }>(
  servers: readonly T[],
  needsGpu: boolean
): T[] {
  if (!needsGpu) return [...servers];
  const yes: T[] = [];
  const rest: T[] = [];
  for (const server of servers) {
    if (serverHasGpu(server)) {
      yes.push(server);
    } else {
      rest.push(server);
    }
  }
  return [...yes, ...rest];
}

export function gpuMismatch(
  gpu: ServerGPUInfo | null | undefined
): 'yes-not-detected' | 'no-but-detected' | null {
  if (!gpu || !gpu.probed) return null;
  if (gpu.override === 'yes' && !gpu.detected) return 'yes-not-detected';
  if (gpu.override === 'no' && gpu.detected) return 'no-but-detected';
  return null;
}
