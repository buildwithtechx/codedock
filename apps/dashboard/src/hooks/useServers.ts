import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateServerRequest } from '#/interfaces/server';
import { serverService } from '#/services/servers';

export const useListServers = () =>
  useQuery({ queryKey: ['servers'], queryFn: () => serverService.list() });

export const useCreateServer = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateServerRequest) => serverService.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['servers'] }),
  });
};
