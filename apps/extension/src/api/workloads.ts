import type { CodedockApp, CodedockProject, ServiceVariable } from '../interfaces';
import { apiRequest } from './client';

export const WorkloadsApi = {
  async listProjects(url: string, token: string): Promise<CodedockProject[]> {
    const res = await apiRequest<CodedockProject[]>(url, token, '/api/projects');
    return Array.isArray(res) ? res : [];
  },

  async listApps(url: string, token: string): Promise<CodedockApp[]> {
    const res = await apiRequest<CodedockApp[]>(url, token, '/api/apps');
    return Array.isArray(res) ? res : [];
  },

  async redeployApp(url: string, token: string, appId: string): Promise<void> {
    await apiRequest<void>(url, token, `/api/apps/${encodeURIComponent(appId)}/redeploy`, {
      method: 'POST',
    });
  },

  async restartApp(url: string, token: string, appId: string): Promise<void> {
    await apiRequest<void>(url, token, `/api/apps/${encodeURIComponent(appId)}/restart`, {
      method: 'POST',
    });
  },

  async stopApp(url: string, token: string, appId: string): Promise<void> {
    await apiRequest<void>(url, token, `/api/apps/${encodeURIComponent(appId)}/stop`, {
      method: 'POST',
    });
  },

  async getAppVariables(url: string, token: string, serviceId: string): Promise<ServiceVariable[]> {
    const res = await apiRequest<ServiceVariable[]>(
      url,
      token,
      `/api/services/${encodeURIComponent(serviceId)}/variables`
    );
    return Array.isArray(res) ? res : [];
  },
};
