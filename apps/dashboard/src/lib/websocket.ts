import { env } from '#/env';

export function websocketUrl(path: string): string {
  const apiUrl = new URL(env.VITE_API_URL || '/api', window.location.origin);
  const protocol = apiUrl.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${protocol}//${apiUrl.host}${apiUrl.pathname.replace(/\/$/, '')}${path}`;
}
