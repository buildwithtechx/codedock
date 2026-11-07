export function cleanBaseUrl(url: string): string {
  return url.trim().replace(/\/+$/, '');
}

export async function apiRequest<T>(
  url: string,
  token: string,
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const base = cleanBaseUrl(url);
  const path = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const fullUrl = `${base}${path}`;

  const headers = new Headers(options.headers || {});
  headers.set('Accept', 'application/json');
  if (token) {
    headers.set('Authorization', `Bearer ${token.trim()}`);
  }
  if (options.body && typeof options.body === 'string' && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const res = await fetch(fullUrl, { ...options, headers });

  if (!res.ok) {
    const errorText = await res.text().catch(() => res.statusText);
    let parsedMessage = errorText;
    try {
      const parsed = JSON.parse(errorText);
      parsedMessage = parsed.message || parsed.error || errorText;
    } catch {
      // Keep plain text
    }
    throw new Error(`[${res.status}] ${parsedMessage || res.statusText}`);
  }

  if (res.status === 204) {
    return {} as T;
  }

  const contentType = res.headers.get('content-type') || '';
  if (contentType.includes('application/json')) {
    return (await res.json()) as T;
  }

  return (await res.text()) as unknown as T;
}
