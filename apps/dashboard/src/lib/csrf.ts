type CsrfEntry = { token?: string; expiresAt: number; pending?: Promise<string> };
const csrfTokens = new Map<string, CsrfEntry>();

export function clearCsrfToken(apiBaseUrl: string): void {
  csrfTokens.delete(apiBaseUrl);
}

export async function bootstrapCsrf(apiBaseUrl: string): Promise<string> {
  const cached = csrfTokens.get(apiBaseUrl);
  if (cached?.pending) return cached.pending;
  if (
    typeof window !== 'undefined' &&
    typeof document !== 'undefined' &&
    new URL(apiBaseUrl, window.location.href).hostname === window.location.hostname
  ) {
    const token = document.cookie.match(/(?:^|;\s*)csrf_token=([^;]*)/)?.[1];
    const decoded = decodeCsrfCookie(token);
    if (decoded) return decoded;
  }
  if (cached?.token && cached.expiresAt > Date.now()) return cached.token;
  const entry: CsrfEntry = { expiresAt: 0 };
  const pending = (async () => {
    try {
      const response = await fetch(`${apiBaseUrl}/auth/csrf`, { credentials: 'include' });
      if (!response.ok) throw new Error('Could not initialize CSRF protection');
      const data = (await response.json()) as { token?: string };
      if (!data.token || data.token === '_echo_csrf_using_sec_fetch_site_')
        throw new Error('Server did not provide a CSRF token');
      entry.token = data.token;
      entry.expiresAt = Date.now() + 23 * 60 * 60 * 1000;
      return data.token;
    } catch (error) {
      csrfTokens.delete(apiBaseUrl);
      throw error;
    } finally {
      entry.pending = undefined;
    }
  })();
  entry.pending = pending;
  csrfTokens.set(apiBaseUrl, entry);
  return pending;
}

export async function prepareCsrfHeaders(
  baseUrl: string,
  endpoint: string,
  method: string,
  headers: Headers
): Promise<void> {
  if (['GET', 'HEAD', 'OPTIONS'].includes(method.toUpperCase())) return;
  if (!headers.has('Authorization') || endpoint.startsWith('/auth/')) {
    headers.set('X-CSRF-Token', await bootstrapCsrf(baseUrl));
  }
}

function decodeCsrfCookie(token?: string): string | undefined {
  if (!token) return undefined;
  try {
    const decoded = decodeURIComponent(token);
    return decoded && decoded !== '_echo_csrf_using_sec_fetch_site_' ? decoded : undefined;
  } catch {
    return undefined;
  }
}
