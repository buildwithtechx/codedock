import axios from 'axios';
import type {
  CodedockApp,
  CodedockDatabase,
  CodedockDeployment,
  CodedockMeta,
  CodedockServer,
  UserProfile,
} from '@/interfaces';
import { cleanUrl } from '@/utils';

function createClient(url: string, token: string = '') {
  return axios.create({
    baseURL: cleanUrl(url),
    headers: {
      Accept: 'application/json',
      ...(token ? { Authorization: `Bearer ${token.trim()}` } : {}),
    },
    timeout: 15000,
  });
}

export const MobileApi = {
  async getMeta(url: string): Promise<CodedockMeta> {
    const client = createClient(url);
    const { data } = await client.get<CodedockMeta>('/api/meta');
    return data;
  },

  async getMe(url: string, token: string): Promise<UserProfile> {
    const client = createClient(url, token);
    const { data } = await client.get<UserProfile>('/api/auth/me');
    return data;
  },

  async listApps(url: string, token: string): Promise<CodedockApp[]> {
    const client = createClient(url, token);
    const { data } = await client.get<CodedockApp[]>('/api/apps');
    return Array.isArray(data) ? data : [];
  },

  async redeployApp(url: string, token: string, appId: string): Promise<void> {
    const client = createClient(url, token);
    await client.post(`/api/apps/${encodeURIComponent(appId)}/redeploy`);
  },

  async restartApp(url: string, token: string, appId: string): Promise<void> {
    const client = createClient(url, token);
    await client.post(`/api/apps/${encodeURIComponent(appId)}/restart`);
  },

  async stopApp(url: string, token: string, appId: string): Promise<void> {
    const client = createClient(url, token);
    await client.post(`/api/apps/${encodeURIComponent(appId)}/stop`);
  },

  async listDatabases(url: string, token: string): Promise<CodedockDatabase[]> {
    const client = createClient(url, token);
    const { data } = await client.get<CodedockDatabase[]>('/api/databases');
    return Array.isArray(data) ? data : [];
  },

  async restartDatabase(url: string, token: string, dbId: string): Promise<void> {
    const client = createClient(url, token);
    await client.post(`/api/databases/${encodeURIComponent(dbId)}/restart`);
  },

  async startDatabase(url: string, token: string, dbId: string): Promise<void> {
    const client = createClient(url, token);
    await client.post(`/api/databases/${encodeURIComponent(dbId)}/start`);
  },

  async stopDatabase(url: string, token: string, dbId: string): Promise<void> {
    const client = createClient(url, token);
    await client.post(`/api/databases/${encodeURIComponent(dbId)}/stop`);
  },

  async revealDatabaseCredentials(
    url: string,
    token: string,
    dbId: string
  ): Promise<{ password?: string }> {
    const client = createClient(url, token);
    const { data } = await client.post<{ password?: string }>(
      `/api/databases/${encodeURIComponent(dbId)}/credentials/reveal`
    );
    return data;
  },

  async listServers(url: string, token: string): Promise<CodedockServer[]> {
    const client = createClient(url, token);
    const { data } = await client.get<CodedockServer[]>('/api/servers');
    return Array.isArray(data) ? data : [];
  },

  async listDeployments(
    url: string,
    token: string,
    limit: number = 20
  ): Promise<CodedockDeployment[]> {
    const client = createClient(url, token);
    const { data } = await client.get<CodedockDeployment[]>(`/api/deployments?limit=${limit}`);
    return Array.isArray(data) ? data : [];
  },

  async getDeploymentLogs(url: string, token: string, id: string): Promise<string> {
    const client = createClient(url, token);
    const { data } = await client.get<string>(`/api/deployments/${encodeURIComponent(id)}/logs`);
    return typeof data === 'string' ? data : JSON.stringify(data);
  },

  async registerDeviceToken(
    url: string,
    token: string,
    deviceToken: string,
    platform: string
  ): Promise<void> {
    const client = createClient(url, token);
    await client.post('/api/notifications/devices', {
      token: deviceToken,
      platform,
    });
  },
};
