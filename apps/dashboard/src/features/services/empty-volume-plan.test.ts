import { describe, expect, it } from 'vitest';
import { emptyVolumePlan } from './empty-volume-plan';

describe('emptyVolumePlan', () => {
  it('returns null when selected volume is not in list', () => {
    expect(emptyVolumePlan(['vol-1:/data'], 'nonexistent')).toBeNull();
  });

  it('returns null for read-only volumes', () => {
    expect(emptyVolumePlan(['vol-1:/data:ro'], 'vol-1:/data:ro')).toBeNull();
  });

  it('generates a fresh empty volume replacement while retaining previous', () => {
    const volumes = ['postgres-data:/var/lib/postgresql/data', 'app-logs:/var/log/app'];
    const plan = emptyVolumePlan(volumes, 'postgres-data:/var/lib/postgresql/data');

    expect(plan).not.toBeNull();
    expect(plan?.previous).toBe('postgres-data');
    expect(plan?.target).toBe('/var/lib/postgresql/data');
    expect(plan?.source).toMatch(/^codedock-fresh-/);
    expect(plan?.volumes).toHaveLength(2);
    expect(plan?.volumes[0]).toBe(`${plan?.source}:/var/lib/postgresql/data`);
    expect(plan?.volumes[1]).toBe('app-logs:/var/log/app');
  });
});
