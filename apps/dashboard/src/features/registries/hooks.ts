import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { registryService } from './api';
import type { CreateRegistryRequest } from './interfaces';

export const useListRegistries = (projectId: string) =>
  useQuery({
    queryKey: ['projects', projectId, 'registries'],
    queryFn: () => registryService.listByProject(projectId),
    enabled: !!projectId,
  });

export const useCreateRegistry = (projectId: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateRegistryRequest) => registryService.create(projectId, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['projects', projectId, 'registries'] }),
  });
};

export const useDeleteRegistry = (projectId: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => registryService.delete(projectId, id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['projects', projectId, 'registries'] }),
  });
};
