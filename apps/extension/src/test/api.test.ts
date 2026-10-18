import * as assert from 'node:assert';
import { afterEach, describe, test } from 'vitest';
import { DatabasesApi, DeploymentsApi, InfrastructureApi, WorkloadsApi } from '../api';

describe('Codedock API Client Suite', () => {
  const originalFetch = globalThis.fetch;

  afterEach(() => {
    globalThis.fetch = originalFetch;
  });

  test('getMeta returns metadata', async () => {
    globalThis.fetch = (async (url: string | URL | Request) => {
      assert.strictEqual(url.toString(), 'http://localhost:8080/api/meta');
      return new Response(
        JSON.stringify({
          name: 'codedock',
          version: 'v1.0.0',
          capabilities: { sseLogs: true, push: true, actions: true, apps: true },
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      );
    }) as typeof fetch;

    const meta = await InfrastructureApi.getMeta('http://localhost:8080/');
    assert.strictEqual(meta.name, 'codedock');
    assert.strictEqual(meta.version, 'v1.0.0');
    assert.strictEqual(meta.capabilities.sseLogs, true);
  });

  test('getMe handles authenticated user', async () => {
    globalThis.fetch = (async (url: string | URL | Request, init?: RequestInit) => {
      assert.strictEqual(url.toString(), 'http://localhost:8080/api/auth/me');
      assert.strictEqual(
        init?.headers && (init.headers as any).get?.('Authorization'),
        'Bearer test-pat'
      );
      return new Response(
        JSON.stringify({
          id: 'user-1',
          email: 'admin@codedock.io',
          name: 'Admin',
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      );
    }) as typeof fetch;

    const me = await InfrastructureApi.getMe('http://localhost:8080', 'test-pat');
    assert.strictEqual(me.id, 'user-1');
    assert.strictEqual(me.email, 'admin@codedock.io');
  });

  test('listApps returns array of apps', async () => {
    globalThis.fetch = (async () => {
      return new Response(
        JSON.stringify([
          { id: 'app-1', name: 'web-service', status: 'running' },
          { id: 'app-2', name: 'api-service', status: 'stopped' },
        ]),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      );
    }) as typeof fetch;

    const apps = await WorkloadsApi.listApps('http://localhost:8080', 'test-pat');
    assert.strictEqual(apps.length, 2);
    assert.strictEqual(apps[0].name, 'web-service');
    assert.strictEqual(apps[1].status, 'stopped');
  });

  test('redeployApp triggers deployment', async () => {
    let called = false;
    globalThis.fetch = (async (url: string | URL | Request) => {
      if (url.toString().includes('/api/apps/app-1/redeploy')) {
        called = true;
        return new Response(JSON.stringify({ success: true }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        });
      }
      return new Response(null, { status: 404 });
    }) as typeof fetch;

    await WorkloadsApi.redeployApp('http://localhost:8080', 'test-pat', 'app-1');
    assert.strictEqual(called, true);
  });

  test('listDatabases returns databases', async () => {
    globalThis.fetch = (async () => {
      return new Response(
        JSON.stringify([{ id: 'db-1', name: 'pg-main', type: 'postgres', status: 'running' }]),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      );
    }) as typeof fetch;

    const dbs = await DatabasesApi.listDatabases('http://localhost:8080', 'test-pat');
    assert.strictEqual(dbs.length, 1);
    assert.strictEqual(dbs[0].name, 'pg-main');
    assert.strictEqual(dbs[0].type, 'postgres');
  });

  test('listDeployments returns deployments', async () => {
    globalThis.fetch = (async () => {
      return new Response(
        JSON.stringify([{ id: 'dep-1', serviceName: 'web', status: 'succeeded' }]),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      );
    }) as typeof fetch;

    const deps = await DeploymentsApi.listDeployments('http://localhost:8080', 'test-pat');
    assert.strictEqual(deps.length, 1);
    assert.strictEqual(deps[0].status, 'succeeded');
  });
});
