import type { DomainConfig } from '#/features/projects/types';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';
import type { CreateDNSRecordRequest, DNSRecord, UpdateDNSRecordRequest } from './types';

export const dnsService = {
  list: async (): Promise<BaseResponse<DNSRecord[]>> => {
    try {
      return await apiClient.get<BaseResponse<DNSRecord[]>>(`/dns`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (payload: CreateDNSRecordRequest): Promise<BaseResponse<DNSRecord>> => {
    try {
      return await apiClient.post<BaseResponse<DNSRecord>>(`/dns`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  update: async (id: string, payload: UpdateDNSRecordRequest): Promise<BaseResponse<DNSRecord>> => {
    try {
      return await apiClient.put<BaseResponse<DNSRecord>>(`/dns/${id}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/dns/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};

export const domainsService = {
  listByService: async (serviceId: string): Promise<BaseResponse<DomainConfig[]>> => {
    try {
      return await apiClient.get<BaseResponse<DomainConfig[]>>(`/services/${serviceId}/domains`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (
    serviceId: string,
    payload: { domainName: string; redirectTo?: string; pathPrefix?: string }
  ): Promise<BaseResponse<DomainConfig>> => {
    try {
      return await apiClient.post<BaseResponse<DomainConfig>>(
        `/services/${serviceId}/domains`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/domains/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
