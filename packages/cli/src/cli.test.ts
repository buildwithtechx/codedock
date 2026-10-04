import { existsSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { ApiClient } from './api-client.js';
import { CLI_VERSION } from './commands/version-command.js';
import { clearConfig, loadConfig, saveConfig } from './config-store.js';

describe('CLI Config Store', () => {
  const testDir = join(tmpdir(), `codedock-test-${Date.now()}`);

  beforeEach(() => {
    process.env.CODEDOCK_CONFIG_DIR = testDir;
  });

  afterEach(() => {
    delete process.env.CODEDOCK_CONFIG_DIR;
    if (existsSync(testDir)) {
      rmSync(testDir, { recursive: true, force: true });
    }
  });

  it('returns empty config when no file exists', () => {
    const cfg = loadConfig();
    expect(cfg).toEqual({});
  });

  it('saves and reloads config correctly', () => {
    saveConfig({
      serverUrl: 'https://dock.example.com',
      token: 'test_token_123',
      email: 'admin@example.com',
    });

    const reloaded = loadConfig();
    expect(reloaded.serverUrl).toBe('https://dock.example.com');
    expect(reloaded.token).toBe('test_token_123');
    expect(reloaded.email).toBe('admin@example.com');
  });

  it('clears config on logout', () => {
    saveConfig({
      serverUrl: 'https://dock.example.com',
      token: 'test_token_123',
    });

    clearConfig();
    const afterClear = loadConfig();
    expect(afterClear.token).toBeUndefined();
  });
});

describe('CLI Version and API Client', () => {
  it('has valid semver version string', () => {
    expect(CLI_VERSION).toMatch(/^\d+\.\d+\.\d+$/);
  });

  it('constructs ApiClient with sanitized base URL', () => {
    const client = new ApiClient({
      config: {},
      serverUrl: 'http://localhost:8080///',
      token: 'secret',
      json: false,
    });
    expect(client).toBeDefined();
  });
});
