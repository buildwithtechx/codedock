import type {
  AppService,
  CreateAppServiceRequest,
  CreateServiceVarRequest,
  EnvExampleVariableSuggestion,
  UpdateAppServiceRequest,
  UpdateServiceVarRequest,
  Variable,
} from '@codedock/contracts';
import type { HttpClient } from '../http';

export function createServicesOperations(http: HttpClient) {
  return {
    listByOrganization: (organizationId: string): Promise<AppService[]> =>
      http.get<AppService[]>(`/apps?organizationId=${encodeURIComponent(organizationId)}`),

    listByProject: (projectId: string, environmentId?: string): Promise<AppService[]> => {
      const url = environmentId
        ? `/projects/${encodeURIComponent(projectId)}/apps?environmentId=${encodeURIComponent(environmentId)}`
        : `/projects/${encodeURIComponent(projectId)}/apps`;
      return http.get<AppService[]>(url);
    },

    get: (id: string): Promise<AppService> =>
      http.get<AppService>(`/apps/${encodeURIComponent(id)}`),

    create: (payload: CreateAppServiceRequest): Promise<AppService> =>
      http.post<AppService>('/apps', payload),

    update: (id: string, payload: UpdateAppServiceRequest): Promise<AppService> =>
      http.put<AppService>(`/apps/${encodeURIComponent(id)}`, payload),

    delete: (id: string): Promise<void> => http.delete<void>(`/apps/${encodeURIComponent(id)}`),

    restart: (id: string): Promise<void> =>
      http.post<void>(`/apps/${encodeURIComponent(id)}/restart`),

    stop: (id: string): Promise<void> => http.post<void>(`/apps/${encodeURIComponent(id)}/stop`),

    rebuild: (id: string): Promise<void> =>
      http.post<void>(`/apps/${encodeURIComponent(id)}/rebuild`),

    listVariables: (serviceId: string): Promise<Variable[]> =>
      http.get<Variable[]>(`/services/${encodeURIComponent(serviceId)}/variables`),

    getEnvSuggestions: (serviceId: string): Promise<EnvExampleVariableSuggestion[]> =>
      http.get<EnvExampleVariableSuggestion[]>(
        `/services/${encodeURIComponent(serviceId)}/env-suggestions`
      ),

    createVariable: (serviceId: string, payload: CreateServiceVarRequest): Promise<Variable> =>
      http.post<Variable>(`/services/${encodeURIComponent(serviceId)}/variables`, payload),

    updateVariable: (
      serviceId: string,
      id: string,
      payload: UpdateServiceVarRequest
    ): Promise<Variable> =>
      http.put<Variable>(
        `/services/${encodeURIComponent(serviceId)}/variables/${encodeURIComponent(id)}`,
        payload
      ),

    deleteVariable: (serviceId: string, id: string): Promise<void> =>
      http.delete<void>(
        `/services/${encodeURIComponent(serviceId)}/variables/${encodeURIComponent(id)}`
      ),
  };
}
