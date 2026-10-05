import { createInterface } from 'node:readline/promises';
import { ApiClient } from '../api-client.js';
import { saveConfig } from '../config-store.js';
import { printError, printInfo, printSuccess } from '../output-format.js';
import { readPassword } from '../password-prompt.js';
import type { CliContext, UserProfile } from '../types.js';

interface LoginOptions {
  server?: string;
  email?: string;
  token?: string;
  totp?: string;
}

export async function loginCommand(
  ctx: CliContext,
  _args: string[],
  options: LoginOptions
): Promise<void> {
  let serverUrl = options.server || ctx.serverUrl;
  let email = options.email;

  const token = options.token;

  let rl = createInterface({
    input: process.stdin,
    output: process.stdout,
  });

  try {
    if (token) {
      if (!serverUrl) {
        serverUrl = await rl.question('Codedock Server URL (e.g. http://localhost:8080): ');
      }
      serverUrl = serverUrl.trim().replace(/\/+$/, '');
      const client = new ApiClient({ ...ctx, serverUrl, token });
      const userRes = await client.get<UserProfile>('/api/auth/me');
      saveConfig({
        serverUrl,
        token,
        email: userRes.data?.email || 'api-token',
      });
      printSuccess(`Authenticated via API token as ${userRes.data?.email || 'user'}`);
      return;
    }

    if (!serverUrl) {
      serverUrl = await rl.question('Codedock Server URL (e.g. http://localhost:8080): ');
    }
    serverUrl = serverUrl.trim().replace(/\/+$/, '');

    if (!email) {
      email = await rl.question('Email: ');
    }
    email = email.trim();
    if (!email) {
      throw new Error('Email is required');
    }

    rl.close();
    const password = await readPassword();
    rl = createInterface({ input: process.stdin, output: process.stdout });
    if (!password) {
      throw new Error('Password is required');
    }

    printInfo('Authenticating with Codedock...');
    const client = new ApiClient({ ...ctx, serverUrl });

    let authRes: { data?: { token?: string; user?: UserProfile } };
    try {
      authRes = await client.post<{ token?: string; user?: UserProfile }>('/api/auth/signin', {
        email,
        password,
        totpCode: options.totp,
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      if (msg.includes('2FA') || msg.includes('totp') || msg.includes('code required')) {
        const totp = await rl.question('Enter 2FA Code: ');
        authRes = await client.post<{ token?: string; user?: UserProfile }>('/api/auth/signin', {
          email,
          password,
          totpCode: totp.trim(),
        });
      } else {
        throw err;
      }
    }

    const receivedToken = authRes.data?.token;
    if (!receivedToken) {
      throw new Error('Authentication succeeded but server returned no token');
    }

    saveConfig({
      serverUrl,
      token: receivedToken,
      email: authRes.data?.user?.email || email,
    });

    printSuccess(`Successfully logged in as ${authRes.data?.user?.email || email} (${serverUrl})`);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    printError(`Login failed: ${msg}`);
    process.exitCode = 1;
  } finally {
    rl.close();
  }
}
