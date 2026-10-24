import type {
  CodedockServer,
  CreateServerRequest,
  TestSSHRequest,
  UpdateServerRequest,
} from '@codedock/contracts';
import type { HttpClient } from '../http';

export function createServersOperations(http: HttpClient) {
  return {
    list: (): Promise<CodedockServer[]> => http.get<CodedockServer[]>('/servers'),

    get: (id: string): Promise<CodedockServer> =>
      http.get<CodedockServer>(`/servers/${encodeURIComponent(id)}`),

    create: (payload: CreateServerRequest): Promise<CodedockServer> =>
      http.post<CodedockServer>('/servers', payload),

    update: (id: string, payload: UpdateServerRequest): Promise<CodedockServer> =>
      http.patch<CodedockServer>(`/servers/${encodeURIComponent(id)}`, payload),

    delete: (id: string): Promise<void> => http.delete<void>(`/servers/${encodeURIComponent(id)}`),

    testSSH: (payload: TestSSHRequest): Promise<{ success: boolean; message: string }> =>
      http.post<{ success: boolean; message: string }>('/servers/test-ssh', payload),
  };
}
