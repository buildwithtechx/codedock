import { describe, expect, it, vi } from 'vitest';
import {
  CodedockApiError,
  CodedockBillingGateError,
  CodedockRateLimitError,
  createCodedockClient,
} from './index';

describe('CodedockClient SDK', () => {
  it('instantiates client with base options and sends authenticated requests', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        data: [
          {
            id: 'ing-1',
            serverId: 'srv-1',
            tunnelId: 'tun-1',
            name: 'Main Ingress',
            status: 'ready',
          },
        ],
      }),
    });

    const client = createCodedockClient({
      baseUrl: 'https://api.codedock.example/api',
      token: 'cdk_secret_token_123',
      fetch: mockFetch,
    });

    const connections = await client.ingress.list();
    expect(connections).toHaveLength(1);
    expect(connections[0].name).toBe('Main Ingress');

    expect(mockFetch).toHaveBeenCalledWith(
      'https://api.codedock.example/api/dns/ingress',
      expect.objectContaining({
        headers: expect.any(Headers),
      })
    );
  });

  it('handles action workflows and runs operations', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        data: {
          id: 'run-1',
          workflowId: 'wf-1',
          name: 'Build',
          number: 1,
          status: 'running',
        },
      }),
    });

    const client = createCodedockClient({ fetch: mockFetch });
    const run = await client.actions.retryRun('run-1');

    expect(run.id).toBe('run-1');
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/actions/runs/run-1/retry',
      expect.objectContaining({ method: 'POST' })
    );
  });

  it('handles networks and watchdog operations', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        data: {
          operationId: 'op-99',
          armed: true,
          remainingSeconds: 120,
          status: 'pending_commit',
        },
      }),
    });

    const client = createCodedockClient({ fetch: mockFetch });
    const watchdog = await client.networks.watchdog('op-99');

    expect(watchdog.armed).toBe(true);
    expect(watchdog.remainingSeconds).toBe(120);
  });

  it('throws CodedockApiError on unsuccessful responses', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 404,
      statusText: 'Not Found',
      json: async () => ({ message: 'resource not found' }),
    });

    const client = createCodedockClient({ fetch: mockFetch });

    await expect(client.ingress.disconnect('invalid-id')).rejects.toThrow(CodedockApiError);
  });

  it('handles account auth and pairing handshake', async () => {
    const mockFetch = vi.fn().mockImplementation(async (url: string, init?: RequestInit) => {
      if (url.endsWith('/auth/account')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({ data: { id: 'usr-1', email: 'owner@example.com', role: 'admin' } }),
        };
      }
      if (url.endsWith('/auth/pairing/code')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({ data: { code: 'pair-code-abc', expiresAt: '2026-10-18T00:00:00Z' } }),
        };
      }
      if (url.endsWith('/auth/pair')) {
        const body = JSON.parse(String(init?.body));
        expect(body.code).toBe('pair-code-abc');
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: {
              token: 'cdk_paired_jwt_456',
              account: { id: 'usr-1', email: 'owner@example.com' },
            },
          }),
        };
      }
      return { ok: false, status: 404, json: async () => ({ error: 'not found' }) };
    });

    const client = createCodedockClient({ fetch: mockFetch });
    const handshake = await client.auth.createPairingCode();
    expect(handshake.code).toBe('pair-code-abc');

    const claimed = await client.auth.pair('pair-code-abc');
    expect(claimed.token).toBe('cdk_paired_jwt_456');
    expect(claimed.account.email).toBe('owner@example.com');

    const account = await client.auth.getAccount();
    expect(account.id).toBe('usr-1');
  });

  it('retains baseUrl and token when scoping to a tenant, forwarding X-Tenant-ID', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ data: [] }),
    });

    const baseClient = createCodedockClient({
      baseUrl: 'https://cloud.codedock.example/api',
      token: 'bearer_token_xyz',
      fetch: mockFetch,
    });

    const scopedClient = baseClient.scope('tenant-org-789');
    await scopedClient.ingress.list();

    expect(mockFetch).toHaveBeenCalledWith(
      'https://cloud.codedock.example/api/dns/ingress',
      expect.objectContaining({
        headers: expect.any(Headers),
      })
    );

    const callHeaders = mockFetch.mock.calls[0][1].headers as Headers;
    expect(callHeaders.get('Authorization')).toBe('Bearer bearer_token_xyz');
    expect(callHeaders.get('X-Tenant-ID')).toBe('tenant-org-789');
  });

  it('parses Retry-After and throws CodedockRateLimitError on 429', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 429,
      statusText: 'Too Many Requests',
      headers: new Headers({ 'Retry-After': '120' }),
      json: async () => ({ code: 'RATE_LIMIT_EXCEEDED', planTier: 'starter' }),
    });

    const client = createCodedockClient({ fetch: mockFetch });
    await expect(client.actions.listRuns()).rejects.toThrow(CodedockRateLimitError);

    try {
      await client.actions.listRuns();
    } catch (err) {
      expect(err).toBeInstanceOf(CodedockRateLimitError);
      const rateErr = err as import('./client').CodedockRateLimitError;
      expect(rateErr.retryAfterSeconds).toBe(120);
      expect(rateErr.planTier).toBe('starter');
    }
  });

  it('detects billing gate errors and throws CodedockBillingGateError with upgrade details', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 402,
      statusText: 'Payment Required',
      headers: new Headers(),
      json: async () => ({
        code: 'BILLING_GATE_REQUIRED',
        message: 'Plan upgrade required for distributed clusters',
        upgradeUrl: 'https://codedock.example/billing/plans',
        requiredPlan: 'team',
      }),
    });

    const client = createCodedockClient({ fetch: mockFetch });
    await expect(client.clusters.setupScaling('cl-1')).rejects.toThrow(CodedockBillingGateError);

    try {
      await client.clusters.setupScaling('cl-1');
    } catch (err) {
      expect(err).toBeInstanceOf(CodedockBillingGateError);
      const gateErr = err as import('./client').CodedockBillingGateError;
      expect(gateErr.upgradeUrl).toBe('https://codedock.example/billing/plans');
      expect(gateErr.requiredPlan).toBe('team');
    }
  });
});
