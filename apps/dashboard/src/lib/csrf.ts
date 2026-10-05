export async function bootstrapCsrf(apiBaseUrl: string): Promise<string> {
  const response = await fetch(`${apiBaseUrl}/auth/csrf`, { credentials: 'include' });
  if (!response.ok) throw new Error('Could not initialize CSRF protection');
  const data = (await response.json()) as { token?: string };
  if (!data.token) throw new Error('Server did not provide a CSRF token');
  return data.token;
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
