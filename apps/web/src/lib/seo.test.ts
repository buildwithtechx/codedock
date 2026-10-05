import { afterEach, expect, it, vi } from 'vitest';

afterEach(() => {
  vi.unstubAllEnvs();
  vi.resetModules();
});
it('resolves canonical and Open Graph URLs with a trailing-slash site origin', async () => {
  vi.stubEnv('VITE_SITE_URL', 'https://example.com/');
  vi.resetModules();
  const { pageHead } = await import('./seo');
  const head = pageHead('/features/databases', 'Databases', 'Database operations');
  expect(head.links).toEqual([
    { rel: 'canonical', href: 'https://example.com/features/databases' },
  ]);
  expect(head.meta).toContainEqual({
    property: 'og:url',
    content: 'https://example.com/features/databases',
  });
});
