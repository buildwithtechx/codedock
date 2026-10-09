import type {
  MaskedStack,
  MigrationPreview,
  MigrationReview,
  MigrationRun,
  MigrationSource,
} from '#/features/servers/migration-types';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const migrationService = {
  listSources: async (organizationId: string): Promise<MigrationSource[]> => {
    try {
      const res = await apiClient.get<BaseResponse<MigrationSource[]>>(
        `/organizations/${organizationId}/migration/sources`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  testSource: async (
    organizationId: string,
    payload: Record<string, unknown>
  ): Promise<{ fingerprint: string }> => {
    try {
      const res = await apiClient.post<BaseResponse<{ fingerprint: string }>>(
        `/organizations/${organizationId}/migration/sources/test`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  createSource: async (
    organizationId: string,
    payload: Record<string, unknown>
  ): Promise<MigrationSource> => {
    try {
      const res = await apiClient.post<BaseResponse<MigrationSource>>(
        `/organizations/${organizationId}/migration/sources`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  deleteSource: async (organizationId: string, sourceId: string): Promise<void> => {
    try {
      await apiClient.delete<BaseResponse<null>>(
        `/organizations/${organizationId}/migration/sources/${sourceId}`
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },
  scan: async (organizationId: string, sourceId: string): Promise<MaskedStack> => {
    try {
      const res = await apiClient.post<BaseResponse<MaskedStack>>(
        `/organizations/${organizationId}/migration/scan`,
        { sourceId }
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  revealEnv: async (
    organizationId: string,
    payload: { sourceId: string; containerId: string }
  ): Promise<Record<string, string>> => {
    try {
      const res = await apiClient.post<BaseResponse<Record<string, string>>>(
        `/organizations/${organizationId}/migration/reveal-env`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  preview: async (
    organizationId: string,
    payload: Record<string, unknown>
  ): Promise<MigrationPreview> => {
    try {
      const res = await apiClient.post<BaseResponse<MigrationPreview>>(
        `/organizations/${organizationId}/migration/preview`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  migrate: async (
    organizationId: string,
    payload: Record<string, unknown>
  ): Promise<MigrationReview> => {
    try {
      const res = await apiClient.post<BaseResponse<MigrationReview>>(
        `/organizations/${organizationId}/migration/migrate`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  moveProject: async (
    organizationId: string,
    payload: Record<string, unknown>
  ): Promise<MigrationReview> => {
    try {
      const res = await apiClient.post<BaseResponse<MigrationReview>>(
        `/organizations/${organizationId}/migration/project`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  getRun: async (runId: string): Promise<Record<string, unknown>> => {
    try {
      const res = await apiClient.get<BaseResponse<Record<string, unknown>>>(
        `/migrations/${runId}`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  cutover: async (
    runId: string,
    payload: { confirmation: string; kill: boolean }
  ): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(`/migrations/${runId}/cutover`, payload);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  cancel: async (runId: string): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(`/migrations/${runId}/cancel`);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  respond: async (
    runId: string,
    payload: { promptId: string; optionId: string; value?: string }
  ): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(`/migrations/${runId}/respond`, payload);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  resume: async (
    runId: string,
    payload: { overrides?: Record<string, string>; skips?: string[] }
  ): Promise<MigrationReview> => {
    try {
      const res = await apiClient.post<BaseResponse<MigrationReview>>(
        `/migrations/${runId}/resume`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  cleanupTarget: async (runId: string): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(`/migrations/${runId}/cleanup-target`);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  listRuns: async (organizationId: string, sourceId?: string): Promise<MigrationRun[]> => {
    try {
      const res = await apiClient.get<BaseResponse<MigrationRun[]>>(
        `/organizations/${organizationId}/migration/runs${sourceId ? `?sourceId=${sourceId}` : ''}`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
