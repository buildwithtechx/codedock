import { describe, expect, it, vi } from 'vitest';
import { websocketUrl } from './websocket';

vi.mock('#/env', () => ({
  env: { VITE_API_URL: '/api' },
}));

describe('websocketUrl', () => {
  it('builds a ws url from the relative api base', () => {
    expect(websocketUrl('/ws/servers/abc/metrics')).toBe(
      `ws://${window.location.host}/api/ws/servers/abc/metrics`
    );
  });

  it('keeps service log paths under the api base', () => {
    expect(websocketUrl('/services/svc-1/logs')).toBe(
      `ws://${window.location.host}/api/services/svc-1/logs`
    );
  });
});
