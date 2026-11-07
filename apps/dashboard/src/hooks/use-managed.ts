import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type {
  CreateManagedCredentialRequest,
  ManagedCatalog,
  ManagedCredential,
  ManagedOffers,
  ManagedOperation,
  ManagedQuota,
  ManagedReview,
  ReviewManagedProvisionRequest,
} from '#/interfaces/managed';
import { managedService } from '#/services/managed';

export const useManagedCatalog = () => {
  return useQuery<ManagedCatalog>({
    queryKey: ['managed', 'catalog'],
    queryFn: () => managedService.catalog(),
  });
};

export const useManagedCredentials = (organizationId?: string) => {
  return useQuery<ManagedCredential[]>({
    queryKey: ['managed', 'credentials', organizationId].filter(Boolean),
    queryFn: () => (organizationId ? managedService.listCredentials(organizationId) : []),
    enabled: Boolean(organizationId),
  });
};

export const useCreateManagedCredential = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      organizationId,
      payload,
    }: {
      organizationId: string;
      payload: CreateManagedCredentialRequest;
    }) => managedService.createCredential(organizationId, payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['managed', 'credentials', variables.organizationId],
      });
    },
  });
};

export const useDeleteManagedCredential = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      organizationId,
      credentialId,
    }: {
      organizationId: string;
      credentialId: string;
    }) => managedService.deleteCredential(organizationId, credentialId),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['managed', 'credentials', variables.organizationId],
      });
    },
  });
};

export const useManagedOffers = (organizationId?: string, credentialId?: string) => {
  return useQuery<ManagedOffers>({
    queryKey: ['managed', 'offers', organizationId, credentialId].filter(Boolean),
    queryFn: () =>
      organizationId && credentialId
        ? managedService.offers(organizationId, credentialId)
        : Promise.reject(new Error('Missing parameters')),
    enabled: Boolean(organizationId && credentialId),
  });
};

export const useManagedQuota = (organizationId?: string) => {
  return useQuery<ManagedQuota | null>({
    queryKey: ['managed', 'quota', organizationId].filter(Boolean),
    queryFn: () => (organizationId ? managedService.getQuota(organizationId) : null),
    enabled: Boolean(organizationId),
  });
};

export const useSetManagedQuota = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      organizationId,
      payload,
    }: {
      organizationId: string;
      payload: { maxServers: number; maxMemoryGB: number };
    }) => managedService.setQuota(organizationId, payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['managed', 'quota', variables.organizationId],
      });
    },
  });
};

export const useReviewManagedProvision = () => {
  return useMutation<
    ManagedReview,
    Error,
    { organizationId: string; payload: ReviewManagedProvisionRequest }
  >({
    mutationFn: ({ organizationId, payload }) =>
      managedService.reviewProvision(organizationId, payload),
  });
};

export const useApplyManagedOperation = () => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, { operationId: string; confirmation: string }>({
    mutationFn: ({ operationId, confirmation }) =>
      managedService.applyOperation(operationId, confirmation),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
  });
};

export const useManagedOperation = (operationId?: string) => {
  return useQuery<ManagedOperation>({
    queryKey: ['managed', 'operation', operationId].filter(Boolean),
    queryFn: () =>
      operationId
        ? managedService.getOperation(operationId)
        : Promise.reject(new Error('No operation ID')),
    enabled: Boolean(operationId),
  });
};
