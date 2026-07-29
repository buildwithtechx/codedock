import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { usersService } from './api';

export const useListUsers = () =>
  useQuery({ queryKey: ['users'], queryFn: () => usersService.list() });

export const useDeleteUser = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => usersService.delete(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['users'] }),
  });
};

export const useInviteUser = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: { email: string; role: string }) => usersService.invite(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['users'] }),
  });
};
