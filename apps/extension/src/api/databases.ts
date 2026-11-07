import type { CodedockDatabase } from '../interfaces';
import { apiRequest } from './client';

export const DatabasesApi = {
  async listDatabases(url: string, token: string): Promise<CodedockDatabase[]> {
    const res = await apiRequest<CodedockDatabase[]>(url, token, '/api/databases');
    return Array.isArray(res) ? res : [];
  },

  async restartDatabase(url: string, token: string, databaseId: string): Promise<void> {
    await apiRequest<void>(url, token, `/api/databases/${encodeURIComponent(databaseId)}/restart`, {
      method: 'POST',
    });
  },

  async stopDatabase(url: string, token: string, databaseId: string): Promise<void> {
    await apiRequest<void>(url, token, `/api/databases/${encodeURIComponent(databaseId)}/stop`, {
      method: 'POST',
    });
  },

  async startDatabase(url: string, token: string, databaseId: string): Promise<void> {
    await apiRequest<void>(url, token, `/api/databases/${encodeURIComponent(databaseId)}/start`, {
      method: 'POST',
    });
  },

  async revealCredentials(
    url: string,
    token: string,
    databaseId: string
  ): Promise<{ connectionString?: string; password?: string; username?: string }> {
    return await apiRequest(
      url,
      token,
      `/api/databases/${encodeURIComponent(databaseId)}/credentials/reveal`,
      { method: 'POST' }
    );
  },
};
