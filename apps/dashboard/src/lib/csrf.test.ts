import { afterEach, describe, expect, it, vi } from 'vitest';
import { bootstrapCsrf, clearCsrfToken } from './csrf';

const base = 'https://remote.example.com/api';
afterEach(() => {
  clearCsrfToken(base);
  vi.unstubAllGlobals();
});
describe('CSRF bootstrap', () => {
  it('deduplicates concurrent requests and reuses the returned token', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValue({ ok: true, json: async () => ({ token: 'token' }) });
    vi.stubGlobal('fetch', fetchMock);
    expect(await Promise.all([bootstrapCsrf(base), bootstrapCsrf(base)])).toEqual([
      'token',
      'token',
    ]);
    expect(await bootstrapCsrf(base)).toBe('token');
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });
  it('allows retry after failed initialization', async () => {
    const fetchMock = vi
      .fn()
      .mockRejectedValueOnce(new Error('offline'))
      .mockResolvedValue({ ok: true, json: async () => ({ token: 'token' }) });
    vi.stubGlobal('fetch', fetchMock);
    await expect(bootstrapCsrf(base)).rejects.toThrow('offline');
    expect(await bootstrapCsrf(base)).toBe('token');
  });
});

it.each(['_echo_csrf_using_sec_fetch_site_', '%invalid', ''])(
  'bootstraps instead of returning an invalid cookie %s',
  async (cookie) => {
    vi.stubGlobal('window', { location: { href: base, hostname: 'remote.example.com' } });
    vi.stubGlobal('document', { cookie: `csrf_token=${cookie}` });
    const fetchMock = vi
      .fn()
      .mockResolvedValue({ ok: true, json: async () => ({ token: 'fresh' }) });
    vi.stubGlobal('fetch', fetchMock);
    expect(await bootstrapCsrf(base)).toBe('fresh');
    expect(fetchMock).toHaveBeenCalledOnce();
  }
);
