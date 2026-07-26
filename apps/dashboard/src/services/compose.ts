import type { CreateDatabaseRequest } from '#/features/databases/types';

import type { CreateAppServiceRequest } from '#/features/services/types';
import { apiClient } from '#/lib/api-client';

export interface ComposeAnalyzeRequest {
  projectId: string;
  composeContent: string;
}

export interface ComposeAnalyzeResponse {
  appServices: CreateAppServiceRequest[];
  databases: CreateDatabaseRequest[];
}

class ComposeService {
  async analyze(req: ComposeAnalyzeRequest): Promise<ComposeAnalyzeResponse> {
    const response = await apiClient.post<{ data: ComposeAnalyzeResponse }>(
      `/projects/${req.projectId}/compose/analyze`,
      req
    );
    return response.data;
  }

  async deploy(projectId: string, composeContent: string): Promise<any> {
    const formData = new FormData();
    formData.append('projectId', projectId);

    const blob = new Blob([composeContent], { type: 'text/yaml' });
    formData.append('file', blob, 'docker-compose.yml');

    const response = await apiClient.post<{ data: any }>(
      `/projects/${projectId}/compose/deploy`,
      formData
    );
    return response.data;
  }
}

export const composeService = new ComposeService();
