import { describe, expect, it } from 'vitest';
import {
  decideCutover,
  hasHostPortBindings,
  newcomerName,
  parseVolumeDestination,
  stableName,
} from '../src/cutover';

describe('cutover helpers', () => {
  it('detects host port bindings', () => {
    expect(hasHostPortBindings([])).toBe(false);
    expect(hasHostPortBindings(['3000'])).toBe(false);
    expect(hasHostPortBindings(['3000/tcp', '8080/udp'])).toBe(false);
    expect(hasHostPortBindings(['8080:80'])).toBe(true);
    expect(hasHostPortBindings(['127.0.0.1:8080:80'])).toBe(true);
    expect(hasHostPortBindings(['0.0.0.0:8080:80/tcp'])).toBe(true);
  });

  it('generates consistent newcomer and stable names', () => {
    expect(newcomerName({ deploymentId: 'dep-1234-abcd', serviceName: 'web' })).toBe(
      'codedock-cutover-dep1234abcd-web'
    );
    expect(stableName({ slug: 'my-app', serviceName: 'api' })).toBe('codedock-my-app-api');
  });

  it('parses volume destinations and ro flags', () => {
    expect(parseVolumeDestination('/data')).toEqual({ target: '/data', isReadOnly: false });
    expect(parseVolumeDestination('/host/data:/app/data')).toEqual({
      target: '/app/data',
      isReadOnly: false,
    });
    expect(parseVolumeDestination('/host/data:/app/data:ro')).toEqual({
      target: '/app/data',
      isReadOnly: true,
    });
  });

  it('decides cutover mode correctly', () => {
    expect(decideCutover({ incumbentRunning: false, ports: [] })).toEqual({
      mode: 'replace',
      reason: 'no-incumbent',
    });

    expect(decideCutover({ incumbentRunning: true, ports: ['8080:80'] })).toEqual({
      mode: 'replace',
      reason: 'host-ports',
    });

    expect(decideCutover({ incumbentRunning: true, ports: ['3000'], networkMode: 'host' })).toEqual({
      mode: 'replace',
      reason: 'shared-namespaces',
    });

    expect(
      decideCutover({
        incumbentRunning: true,
        ports: ['3000'],
        volumes: ['/host/data:/app/data'],
        incumbentMounts: [{ Destination: '/app/data', RW: true }],
      })
    ).toEqual({
      mode: 'replace',
      reason: 'shared-rw-volume',
      clashes: ['/app/data'],
    });

    expect(
      decideCutover({
        incumbentRunning: true,
        ports: ['3000'],
        volumes: ['/host/data:/app/data'],
        incumbentMounts: [{ Destination: '/app/data', RW: true }],
        sharedOptIn: new Set(['/app/data']),
      })
    ).toEqual({
      mode: 'overlap',
    });
  });
});
