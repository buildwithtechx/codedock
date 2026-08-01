import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { dnsService, domainsService } from './api';
import type { CreateDNSRecordRequest, UpdateDNSRecordRequest } from './interfaces';

export const useListDNS = () => {
  return useQuery({
    queryKey: ['dns'],
    queryFn: () => dnsService.list(),
  });
};

export const useCreateDNS = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateDNSRecordRequest) => dnsService.create(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['dns'] });
    },
  });
};

export const useUpdateDNS = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateDNSRecordRequest }) =>
      dnsService.update(id, payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['dns'] });
    },
  });
};

export const useDeleteDNS = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => dnsService.delete(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['dns'] });
    },
  });
};

export const useListAllDomains = () => {
  return useQuery({
    queryKey: ['domains', 'all'],
    queryFn: () => domainsService.listAll(),
  });
};

export const useListByService = (serviceId: string) => {
  return useQuery({
    queryKey: ['domains', 'listByService', serviceId].filter(Boolean),
    queryFn: () => domainsService.listByService(serviceId),
  });
};

export const useCreateDomain = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      serviceId: string;
      payload: { domainName: string; redirectTo?: string; pathPrefix?: string };
    }) => domainsService.create(payload.serviceId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['domains'] });
    },
  });
};

export const useDeleteDomain = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => domainsService.delete(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['domains'] });
    },
  });
};

export const useVerifyDomain = () => {
  return useMutation({
    mutationFn: (domainId: string) => domainsService.verifyDomain(domainId),
  });
};

export const useCreate = useCreateDomain;
export const useDelete = useDeleteDomain;
export const useVerify = useVerifyDomain;
