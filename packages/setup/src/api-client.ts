import type { ServerSshSettings, SetupServerPayload } from './types';

export function buildSetupPayload(options: {
  name: string;
  host: string;
  sshSettings?: Partial<ServerSshSettings>;
}): SetupServerPayload {
  const { name, host, sshSettings } = options;
  return {
    name: name.trim() || host.trim(),
    ipAddress: host.trim(),
    sshPort: sshSettings?.sshPort ?? 22,
    sshUser: sshSettings?.sshUser?.trim() || 'root',
    sshAuthMethod: sshSettings?.sshAuthMethod || 'key',
    sshKey: sshSettings?.sshKey,
    sshPassword: sshSettings?.sshPassword,
    sshJumpHost: sshSettings?.sshJumpHost,
  };
}

export async function pushServerSetup(
  baseUrl: string,
  payload: SetupServerPayload,
  token?: string,
  fetchImpl: typeof fetch = globalThis.fetch
): Promise<{ id: string; name: string }> {
  const normalizedBase = baseUrl.replace(/\/+$/, '');
  const url = `${normalizedBase}/servers`;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetchImpl(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    let errorMsg = `Server setup failed: ${response.statusText}`;
    try {
      const errBody = (await response.json()) as { error?: string; message?: string };
      errorMsg = errBody.message || errBody.error || errorMsg;
    } catch {
      // ignore json parse failure
    }
    throw new Error(errorMsg);
  }

  const data = (await response.json()) as {
    id?: string;
    name?: string;
    data?: { id: string; name: string };
  };
  return data.data || { id: data.id || '', name: data.name || payload.name };
}

export async function waitForApi(
  url: string,
  timeoutMs = 15000,
  intervalMs = 1000,
  fetchImpl: typeof fetch = globalThis.fetch
): Promise<boolean> {
  const deadline = Date.now() + timeoutMs;

  while (Date.now() < deadline) {
    try {
      const res = await fetchImpl(url, { method: 'GET' });
      if (res.ok) {
        return true;
      }
    } catch {
      // transient network or connection error while polling
    }
    await new Promise((resolve) => setTimeout(resolve, intervalMs));
  }

  return false;
}
