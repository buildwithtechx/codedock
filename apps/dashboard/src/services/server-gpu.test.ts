import { describe, expect, it } from 'vitest';
import type { ServerGPUInfo } from '#/interfaces/server';
import { gpuMismatch, serverHasGpu, sortForGpu } from './server-gpu';

describe('server-gpu helpers', () => {
  it('detects if server has GPU available', () => {
    expect(serverHasGpu(null)).toBe(false);
    expect(serverHasGpu({})).toBe(false);
    expect(serverHasGpu({ gpu: null })).toBe(false);
    expect(serverHasGpu({ gpu: { available: false } as ServerGPUInfo })).toBe(false);
    expect(serverHasGpu({ gpu: { available: true } as ServerGPUInfo })).toBe(true);
  });

  it('sorts GPU servers first when requested', () => {
    const list = [
      { id: '1', gpu: { available: false } as ServerGPUInfo },
      { id: '2', gpu: { available: true } as ServerGPUInfo },
      { id: '3', gpu: null },
    ];
    const sorted = sortForGpu(list, true);
    expect(sorted.map((s) => s.id)).toEqual(['2', '1', '3']);

    const unsorted = sortForGpu(list, false);
    expect(unsorted.map((s) => s.id)).toEqual(['1', '2', '3']);
  });

  it('detects GPU mismatch override vs probed state', () => {
    expect(gpuMismatch(null)).toBeNull();
    expect(gpuMismatch({ probed: false } as ServerGPUInfo)).toBeNull();
    expect(gpuMismatch({ probed: true, override: 'yes', detected: false } as ServerGPUInfo)).toBe(
      'yes-not-detected'
    );
    expect(gpuMismatch({ probed: true, override: 'no', detected: true } as ServerGPUInfo)).toBe(
      'no-but-detected'
    );
    expect(
      gpuMismatch({ probed: true, override: 'yes', detected: true } as ServerGPUInfo)
    ).toBeNull();
  });
});
