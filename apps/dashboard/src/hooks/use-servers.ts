import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateServerRequest, TestSSHRequest, UpdateServerRequest } from '#/interfaces/server';
import { serverService } from '#/services/servers';

export const useListServers = () =>
  useQuery({ queryKey: ['servers'], queryFn: () => serverService.list() });

export const useServer = (id?: string) =>
  useQuery({
    queryKey: ['servers', id],
    queryFn: () => (id ? serverService.get(id) : Promise.reject(new Error('No server ID'))),
    enabled: Boolean(id),
  });

export const useCreateServer = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateServerRequest) => serverService.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['servers'] }),
  });
};

export const useUpdateServer = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateServerRequest }) =>
      serverService.update(id, payload),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: ['servers'] });
      qc.invalidateQueries({ queryKey: ['servers', vars.id] });
    },
  });
};

export const useDeleteServer = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => serverService.delete(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['servers'] }),
  });
};

export const useTestSSH = () =>
  useMutation({
    mutationFn: (payload: TestSSHRequest) => serverService.testSSH(payload),
  });
