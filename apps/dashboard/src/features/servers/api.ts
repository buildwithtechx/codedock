import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateServerRequest,
  Server,
  TestSSHRequest,
  UpdateServerRequest,
} from '#/interfaces/server';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export interface ServerConnectionVerdict {
  ok: boolean;
  message: string;
}

export async function listServers(): Promise<Server[]> {
  try {
    const res = await apiClient.get<BaseResponse<Server[]>>('/servers');
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getServer(id: string): Promise<Server> {
  try {
    const res = await apiClient.get<BaseResponse<Server>>(`/servers/${id}`);
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function createServer(payload: CreateServerRequest): Promise<Server> {
  try {
    const res = await apiClient.post<BaseResponse<Server>>('/servers', payload);
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function updateServer(id: string, payload: UpdateServerRequest): Promise<Server> {
  try {
    const res = await apiClient.patch<BaseResponse<Server>>(`/servers/${id}`, payload);
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function removeServer(id: string): Promise<void> {
  try {
    await apiClient.delete<BaseResponse<null>>(`/servers/${id}`);
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function testServerConnection(
  payload: TestSSHRequest
): Promise<ServerConnectionVerdict> {
  try {
    const res = await apiClient.post<BaseResponse<null>>('/servers/test-ssh', payload);
    return { ok: true, message: res.message || 'SSH connection successful' };
  } catch (err) {
    const error = handleApiError(err);
    return { ok: false, message: error.message };
  }
}
