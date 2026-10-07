import { afterEach, describe, expect, it, vi } from 'vitest';
import {
  clearToken,
  defaultConnection,
  getConnection,
  isDefaultConnection,
  loadToken,
  saveToken,
  setDaemonPort,
  setDaemonUrl,
} from './connection-client.js';
import {
  getSidecarStatus,
  restartSidecar,
  setSidecarPort,
  startSidecar,
  stopSidecar,
} from './sidecar-client.js';
import { setInvokeForTesting } from './tauri-ipc.js';
import { checkUpdate, downloadAndInstall, restartApp } from './update-client.js';

afterEach(() => {
  setInvokeForTesting(null);
});

describe('sidecar client', () => {
  it('requests status through the sidecar_status command', async () => {
    const status = {
      running: true,
      healthy: true,
      restarts: 2,
      pid: 4242,
      lastError: null,
      port: 8080,
    };
    const invoke = vi.fn(async () => status);
    setInvokeForTesting(invoke);

    await expect(getSidecarStatus()).resolves.toEqual(status);
    expect(invoke).toHaveBeenCalledWith('sidecar_status', undefined);
  });

  it('maps lifecycle commands', async () => {
    const invoke = vi.fn(async () => undefined);
    setInvokeForTesting(invoke);

    await restartSidecar();
    await stopSidecar();
    await startSidecar();
    await setSidecarPort(9090);

    expect(invoke).toHaveBeenNthCalledWith(1, 'sidecar_restart', undefined);
    expect(invoke).toHaveBeenNthCalledWith(2, 'sidecar_stop', undefined);
    expect(invoke).toHaveBeenNthCalledWith(3, 'sidecar_start', undefined);
    expect(invoke).toHaveBeenNthCalledWith(4, 'sidecar_set_port', { port: 9090 });
  });

  it('propagates sidecar errors', async () => {
    setInvokeForTesting(vi.fn(async () => Promise.reject(new Error('spawn daemon: denied'))));

    await expect(getSidecarStatus()).rejects.toThrow('spawn daemon: denied');
  });
});

describe('connection client', () => {
  it('exposes the local daemon defaults', () => {
    expect(defaultConnection()).toEqual({
      daemonUrl: 'http://localhost:8080',
      daemonPort: 8080,
      hasToken: false,
    });
    expect(isDefaultConnection(defaultConnection())).toBe(true);
    expect(
      isDefaultConnection({ daemonUrl: 'https://ops.example', daemonPort: 443, hasToken: true })
    ).toBe(false);
  });

  it('maps connection and token commands', async () => {
    const invoke = vi.fn(async (cmd: string) => {
      if (cmd === 'get_connection') {
        return defaultConnection();
      }
      if (cmd === 'load_token') {
        return 'secret';
      }
      return undefined;
    });
    setInvokeForTesting(invoke);

    await expect(getConnection()).resolves.toEqual(defaultConnection());
    await setDaemonUrl('https://ops.example/');
    await setDaemonPort(443);
    await saveToken('secret');
    await expect(loadToken()).resolves.toBe('secret');
    await clearToken();

    expect(invoke).toHaveBeenCalledWith('set_daemon_url', { url: 'https://ops.example/' });
    expect(invoke).toHaveBeenCalledWith('set_daemon_port', { port: 443 });
    expect(invoke).toHaveBeenCalledWith('save_token', { token: 'secret' });
    expect(invoke).toHaveBeenCalledWith('clear_token', undefined);
  });
});

describe('update client', () => {
  it('maps update commands', async () => {
    const info = { available: true, version: '0.2.0', notes: 'fixes', date: '2026-10-01' };
    const invoke = vi.fn(async (cmd: string) => (cmd === 'check_update' ? info : undefined));
    setInvokeForTesting(invoke);

    await expect(checkUpdate()).resolves.toEqual(info);
    await downloadAndInstall();
    await restartApp();

    expect(invoke).toHaveBeenCalledWith('download_and_install', undefined);
    expect(invoke).toHaveBeenCalledWith('restart_app', undefined);
  });
});
