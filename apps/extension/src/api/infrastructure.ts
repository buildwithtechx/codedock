import type { CodedockMeta, CodedockServer, SystemStats, UserProfile } from '../interfaces';
import { apiRequest } from './client';

export const InfrastructureApi = {
  async getMeta(url: string): Promise<CodedockMeta> {
    return await apiRequest<CodedockMeta>(url, '', '/api/meta');
  },

  async getMe(url: string, token: string): Promise<UserProfile> {
    return await apiRequest<UserProfile>(url, token, '/api/auth/me');
  },

  async startDeviceAuth(
    url: string,
    clientName = 'VS Code Extension'
  ): Promise<{
    deviceCode: string;
    userCode: string;
    verificationUri: string;
    verificationUriComplete: string;
    expiresIn: number;
    interval: number;
  }> {
    const res = await apiRequest<{
      data: {
        deviceCode: string;
        userCode: string;
        verificationUri: string;
        verificationUriComplete: string;
        expiresIn: number;
        interval: number;
      };
    }>(url, '', '/api/auth/device/code', {
      method: 'POST',
      body: JSON.stringify({ clientName }),
    });
    return res.data;
  },

  async pollDeviceToken(
    url: string,
    deviceCode: string
  ): Promise<{
    status: 'pending' | 'approved' | 'denied';
    token?: string;
    email?: string;
    name?: string;
  }> {
    const res = await apiRequest<{
      data: {
        status: 'pending' | 'approved' | 'denied';
        token?: string;
        email?: string;
        name?: string;
      };
    }>(url, '', '/api/auth/device/token', {
      method: 'POST',
      body: JSON.stringify({ deviceCode }),
    });
    return res.data;
  },

  async listServers(url: string, token: string): Promise<CodedockServer[]> {
    const res = await apiRequest<CodedockServer[]>(url, token, '/api/servers');
    return Array.isArray(res) ? res : [];
  },

  async testServerSSH(
    url: string,
    token: string,
    serverId: string
  ): Promise<{ success: boolean; message?: string }> {
    return await apiRequest(url, token, '/api/servers/test-ssh', {
      method: 'POST',
      body: JSON.stringify({ serverId }),
    });
  },

  async getSystemStats(url: string, token: string): Promise<SystemStats> {
    return await apiRequest<SystemStats>(url, token, '/api/system/stats');
  },
};
