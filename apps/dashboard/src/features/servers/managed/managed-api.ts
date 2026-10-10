import type { BaseResponse } from '#/interfaces/base';
import type {
  ManagedRuntimeCredential,
  ManagedRuntimeStatus,
  ManagedServerInfo,
  ManagedSshStatus,
  ManagedWorkload,
  ServerNetworkSettings,
} from '#/interfaces/server';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export async function getServerNetworkSettings(serverId: string): Promise<ServerNetworkSettings> {
  try {
    const res = await apiClient.get<BaseResponse<ServerNetworkSettings>>(
      `/servers/${serverId}/network-settings`
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function updateServerNetworkSettings(
  serverId: string,
  payload: {
    internetAccess: boolean;
    expectedInternetAccess: boolean;
    egress?: string[];
    expectedRevision?: string;
    confirm: true;
  }
): Promise<ServerNetworkSettings> {
  try {
    const res = await apiClient.patch<BaseResponse<ServerNetworkSettings>>(
      `/servers/${serverId}/network-settings`,
      payload
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getManagedInfo(serverId: string): Promise<ManagedServerInfo> {
  try {
    const res = await apiClient.get<BaseResponse<ManagedServerInfo>>(
      `/servers/${serverId}/managed/info`
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getManagedBootLogs(
  serverId: string,
  tail = 200
): Promise<{ logs: string; truncated: boolean }> {
  try {
    const res = await apiClient.post<BaseResponse<{ logs: string; truncated: boolean }>>(
      `/servers/${serverId}/managed/boot-logs`,
      { tail }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getManagedSshStatus(serverId: string): Promise<ManagedSshStatus> {
  try {
    const res = await apiClient.get<BaseResponse<ManagedSshStatus>>(
      `/servers/${serverId}/managed/ssh`
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function setManagedSsh(
  serverId: string,
  payload: { enabled: boolean; expectedEnabled: boolean; confirm: true }
): Promise<{ status: ManagedSshStatus; initialPassword?: string }> {
  try {
    const res = await apiClient.patch<
      BaseResponse<{ status: ManagedSshStatus; initialPassword?: string }>
    >(`/servers/${serverId}/managed/ssh`, payload);
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function setManagedSshKey(
  serverId: string,
  publicKey: string
): Promise<{ ok: boolean }> {
  try {
    const res = await apiClient.put<BaseResponse<{ ok: boolean }>>(
      `/servers/${serverId}/managed/ssh/key`,
      { publicKey, confirm: true }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function setManagedSshPassword(
  serverId: string,
  password: string
): Promise<{ ok: boolean }> {
  try {
    const res = await apiClient.put<BaseResponse<{ ok: boolean }>>(
      `/servers/${serverId}/managed/ssh/password`,
      { password, confirm: true }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getManagedRuntimeStatus(serverId: string): Promise<ManagedRuntimeStatus> {
  try {
    const res = await apiClient.get<BaseResponse<ManagedRuntimeStatus>>(
      `/servers/${serverId}/managed/runtime-api`
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function enableManagedRuntime(serverId: string): Promise<ManagedRuntimeStatus> {
  try {
    const res = await apiClient.post<BaseResponse<ManagedRuntimeStatus>>(
      `/servers/${serverId}/managed/runtime-api/enable`,
      { confirm: true }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getManagedRuntimeCredential(
  serverId: string
): Promise<ManagedRuntimeCredential> {
  try {
    const res = await apiClient.post<BaseResponse<ManagedRuntimeCredential>>(
      `/servers/${serverId}/managed/runtime-api/credential`,
      { confirm: true }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function rotateManagedRuntimeCredential(
  serverId: string
): Promise<ManagedRuntimeCredential> {
  try {
    const res = await apiClient.post<BaseResponse<ManagedRuntimeCredential>>(
      `/servers/${serverId}/managed/runtime-api/rotate`,
      { confirm: true }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getManagedWorkloads(
  serverId: string
): Promise<{ workloads: ManagedWorkload[]; truncated: boolean }> {
  try {
    const res = await apiClient.get<
      BaseResponse<{ workloads: ManagedWorkload[]; truncated: boolean }>
    >(`/servers/${serverId}/managed/workloads`);
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function createManagedWorkload(
  serverId: string,
  payload: {
    name: string;
    command: string;
    workingDirectory?: string;
    environment?: string[];
    restartPolicy?: string;
    confirm: true;
  }
): Promise<ManagedWorkload> {
  try {
    const res = await apiClient.post<BaseResponse<ManagedWorkload>>(
      `/servers/${serverId}/managed/workloads`,
      {
        ...payload,
        idempotencyKey: crypto.randomUUID(),
      }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function controlManagedWorkload(
  serverId: string,
  payload: {
    workloadId: string;
    action: 'start' | 'stop' | 'restart' | 'delete';
    confirm: true;
  }
): Promise<{ ok: boolean; workload?: ManagedWorkload }> {
  try {
    const res = await apiClient.post<BaseResponse<{ ok: boolean; workload?: ManagedWorkload }>>(
      `/servers/${serverId}/managed/workloads/control`,
      payload
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}

export async function getManagedWorkloadLogs(
  serverId: string,
  workloadId: string,
  tail = 100
): Promise<{ logs: string; truncated: boolean }> {
  try {
    const res = await apiClient.post<BaseResponse<{ logs: string; truncated: boolean }>>(
      `/servers/${serverId}/managed/workloads/logs`,
      { workloadId, tail }
    );
    return res.data;
  } catch (err) {
    throw handleApiError(err);
  }
}
