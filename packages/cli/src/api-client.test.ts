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
    vi.fn().mockImplementation(async (_url: string, init: RequestInit) => {
      const reader = (init.body as ReadableStream).getReader();
      try {
        await reader.read();
      } catch (cause) {
        throw new TypeError('fetch failed', { cause });
      }
    })
  );
  const source = new ReadableStream<Uint8Array>({
    pull(controller) {
      controller.error(archiveError);
    },
  });
  const client = new ApiClient({
    config: {},
    serverUrl: 'http://localhost:8080',
    token: 'test',
    json: false,
  });
  await expect(client.postStream('/api/deploy/archive', source, 'application/gzip')).rejects.toBe(
    archiveError
  );
});

it('preserves server context for streaming transport errors', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockRejectedValue(new TypeError('fetch failed', { cause: new Error('ECONNREFUSED') }))
  );
  const client = new ApiClient({
    config: {},
    serverUrl: 'http://localhost:8080',
    token: 'test',
    json: false,
  });
  await expect(
    client.postStream('/api/deploy/archive', new ReadableStream(), 'application/gzip')
  ).rejects.toThrow('Failed to connect to Codedock server at http://localhost:8080');
});
