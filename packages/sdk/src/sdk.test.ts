import { describe, expect, it, vi } from 'vitest';
import { CodedockApiError, createCodedockClient } from './index';

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
});
