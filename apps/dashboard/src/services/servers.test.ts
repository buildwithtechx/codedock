import { beforeEach, describe, expect, it, vi } from 'vitest';
import { apiClient } from '#/lib/api-client';
import { getServerMonitorStreamUrl, serverService } from './servers';

vi.mock('#/lib/api-client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('serverService and monitor stream', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('builds standard server monitor stream URL without format query', () => {
    const url = getServerMonitorStreamUrl('server-123');
    expect(url).toBe('/api/system/monitor/stream?serverId=server-123');
  });

  it('builds server monitor stream URL with opt-in usage format', () => {
    const url = getServerMonitorStreamUrl('server-123', { format: 'usage' });
    expect(url).toBe('/api/system/monitor/stream?serverId=server-123&format=usage');
  });

  it('builds server monitor stream URL with stats format', () => {
    const url = getServerMonitorStreamUrl('server-abc', { format: 'stats' });
    expect(url).toBe('/api/system/monitor/stream?serverId=server-abc&format=stats');
  });

  it('delegates through serverService.getMonitorStreamUrl', () => {
    const url = serverService.getMonitorStreamUrl('srv-456', { format: 'usage' });
    expect(url).toBe('/api/system/monitor/stream?serverId=srv-456&format=usage');
  });

  it('fetches server list via apiClient', async () => {
    const mockServers = [{ id: 'srv-1', name: 'Node 1' }];
    vi.mocked(apiClient.get).mockResolvedValueOnce({
      success: true,
      message: 'OK',
      data: mockServers,
    });

    const result = await serverService.list();
    expect(apiClient.get).toHaveBeenCalledWith('/servers');
    expect(result).toEqual(mockServers);
  });
});
