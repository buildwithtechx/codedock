import { describe, expect, it } from 'vitest';
import type { ServerGPUInfo } from '#/interfaces/server';
import { gpuMismatch, serverHasGpu, sortForGpu } from './server-gpu';

const mockGpu = (over: Partial<ServerGPUInfo> = {}): ServerGPUInfo => ({
  available: true,
  detected: true,
  count: 1,
  override: null,
  probed: true,
  ...over,
});

const makeServer = (id: string, gpu?: ServerGPUInfo | null) => ({ id, gpu });

describe('serverHasGpu', () => {
  it('returns false for servers without GPU annotations', () => {
    expect(serverHasGpu(undefined)).toBe(false);
    expect(serverHasGpu(null)).toBe(false);
    expect(serverHasGpu({})).toBe(false);
    expect(serverHasGpu({ gpu: null })).toBe(false);
  });

  it('follows the available flag', () => {
    expect(serverHasGpu({ gpu: mockGpu() })).toBe(true);
    expect(serverHasGpu({ gpu: mockGpu({ available: false }) })).toBe(false);
  });
});

describe('sortForGpu', () => {
  const list = [
    makeServer('a'),
    makeServer('b', mockGpu()),
    makeServer('c', mockGpu({ available: false })),
    makeServer('d', mockGpu()),
  ];

  it('lists GPU servers first while keeping relative order inside each partition', () => {
    expect(sortForGpu(list, true).map((s) => s.id)).toEqual(['b', 'd', 'a', 'c']);
  });

  it('keeps original order when GPU is not needed', () => {
    expect(sortForGpu(list, false).map((s) => s.id)).toEqual(['a', 'b', 'c', 'd']);
  });

  it('never drops servers from the array', () => {
    expect(sortForGpu(list, true)).toHaveLength(4);
    expect(sortForGpu(list, false)).toHaveLength(4);
  });
});

describe('gpuMismatch', () => {
  it('detects when manual override is yes but hardware probe found none', () => {
    expect(gpuMismatch(mockGpu({ override: 'yes', detected: false, available: true }))).toBe(
      'yes-not-detected'
    );
  });

  it('detects when manual override is no but hardware probe found a GPU', () => {
    expect(gpuMismatch(mockGpu({ override: 'no', detected: true, available: false }))).toBe(
      'no-but-detected'
    );
  });

  it('returns null when override matches detection or is not set', () => {
    expect(gpuMismatch(mockGpu({ override: 'yes', detected: true }))).toBeNull();
    expect(gpuMismatch(mockGpu({ override: 'no', detected: false }))).toBeNull();
    expect(gpuMismatch(mockGpu({ override: null, detected: true }))).toBeNull();
  });

  it('returns null when the server was not probed', () => {
    expect(gpuMismatch(mockGpu({ override: 'yes', detected: false, probed: false }))).toBeNull();
  });
});
