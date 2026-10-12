import { describe, expect, it } from 'vitest';
import type { ManagedOperation } from '#/features/servers/managed-types';
import { hetznerFailed, hetznerStageIndex, hetznerStages } from './hetzner-progress-stages';

function snapshot(partial: Partial<ManagedOperation>): ManagedOperation {
  return {
    id: 'op-1',
    status: 'RUNNING',
    phase: '',
    effects: '',
    logs: '',
    error: '',
    expiresAt: 0,
    ...partial,
  };
}

describe('hetzner progress stages', () => {
  it('walks creating, boot, bootstrap and verified in order', () => {
    expect(hetznerStageIndex(null)).toBe(-1);
    expect(hetznerStageIndex(snapshot({ phase: 'PROVISIONING', logs: 'Reserving edge-1' }))).toBe(
      0
    );
    expect(
      hetznerStageIndex(snapshot({ phase: 'BOOTSTRAP', logs: 'Waiting for 123 to boot' }))
    ).toBe(1);
    expect(
      hetznerStageIndex(snapshot({ phase: 'BOOTSTRAP', logs: 'Waiting for ssh on 203.0.113.10' }))
    ).toBe(2);
    expect(
      hetznerStageIndex(
        snapshot({ phase: 'BOOTSTRAP', logs: 'Waiting for docker on 203.0.113.10' })
      )
    ).toBe(2);
    expect(
      hetznerStageIndex(
        snapshot({ status: 'COMPLETED', phase: 'READY', logs: 'Managed server ready' })
      )
    ).toBe(hetznerStages.length);
  });

  it('flags failed and interrupted operations', () => {
    expect(hetznerFailed(snapshot({ status: 'FAILED' }))).toBe(true);
    expect(hetznerFailed(snapshot({ status: 'INTERRUPTED' }))).toBe(true);
    expect(hetznerFailed(snapshot({ status: 'RUNNING' }))).toBe(false);
    expect(hetznerFailed(null)).toBe(false);
  });
});
