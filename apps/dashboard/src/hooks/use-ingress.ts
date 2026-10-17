import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type {
  IngressConnectInput,
  IngressDomainInput,
  IngressUpdateInput,
} from '#/interfaces/ingress';
import { ingressService } from '#/services/ingress';

export function useIngressConnections(serverId?: string) {
  return useQuery({
    queryKey: ['dns', 'ingress', serverId],
    queryFn: () => ingressService.list(serverId),
  });
}

export function useIngressHostnames(domainId?: string) {
  return useQuery({
    queryKey: ['dns', 'ingress', 'hostnames', domainId],
    queryFn: () => ingressService.listHostnames(domainId),
  });
}

export function useDiscoverIngress() {
  return useMutation({
    mutationFn: (input: { credentialId: string; accountId?: string; includePolicies?: boolean }) =>
      ingressService.discover(input),
  });
}

export function useConnectIngress() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: IngressConnectInput) => ingressService.connect(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dns', 'ingress'] });
    },
  });
}

export function useUpdateIngress() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: IngressUpdateInput) => ingressService.update(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dns', 'ingress'] });
    },
  });
}

export function useRetryIngress() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => ingressService.retry(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dns', 'ingress'] });
    },
  });
}

export function useDisconnectIngress() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => ingressService.disconnect(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dns', 'ingress'] });
    },
  });
}

export function useBindIngressDomain() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: IngressDomainInput) => ingressService.bindDomain(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dns', 'ingress'] });
    },
  });
}
