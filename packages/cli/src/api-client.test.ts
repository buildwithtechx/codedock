import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiClient } from './api-client.js';

afterEach(() => vi.unstubAllGlobals());
describe('CSRF requests', () => {
  it('reuses a single bootstrap for simultaneous auth requests', async () => {
    const fetchMock = vi
      .fn()
      .mockImplementation(
        async (url: string) =>
          new Response(
            JSON.stringify(url.endsWith('/csrf') ? { token: 'csrf-token' } : { data: {} }),
            { status: 200 }
          )
      );
    vi.stubGlobal('fetch', fetchMock);
    const client = new ApiClient({ config: {}, serverUrl: 'http://localhost:8080', json: false });
    await Promise.all([client.post('/api/auth/signin', {}), client.post('/api/auth/signin', {})]);
    expect(fetchMock.mock.calls.filter(([url]) => url.endsWith('/csrf'))).toHaveLength(1);
    expect(fetchMock.mock.calls).toHaveLength(3);
  });
  it('reports connection context when bootstrap cannot connect', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('fetch failed')));
    const client = new ApiClient({ config: {}, serverUrl: 'http://localhost:8080', json: false });
    await expect(client.post('/api/auth/signin', {})).rejects.toThrow(
      'Failed to connect to Codedock server at http://localhost:8080'
    );
  });
});

it('preserves archive stream errors reported by fetch', async () => {
  const archiveError = new Error('Deployment archive exceeds the 500 MB limit');
  vi.stubGlobal(
    'fetch',
    vi.fn().mockRejectedValue(new TypeError('fetch failed', { cause: archiveError }))
  );
  const client = new ApiClient({
    config: {},
    serverUrl: 'http://localhost:8080',
    token: 'test',
    json: false,
  });
  await expect(
    client.postStream('/api/deploy/archive', new ReadableStream(), 'application/gzip')
  ).rejects.toBe(archiveError);
});
