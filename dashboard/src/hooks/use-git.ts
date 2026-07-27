import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { ConnectGitRequest } from '#/interfaces/git';
import { gitService } from '#/services/git';

export const useGitStatus = () => {
  return useQuery({
    queryKey: ['git', 'status'],
    queryFn: () => gitService.getStatus(),
  });
};

export const useListGitRepos = (provider: string) => {
  return useQuery({
    queryKey: ['git', 'repos', provider],
    queryFn: () => gitService.listRepos(provider),
    enabled: !!provider && provider !== 'public',
  });
};

export const useListGitBranches = (provider: string, repo: string) => {
  return useQuery({
    queryKey: ['git', 'branches', provider, repo],
    queryFn: () => gitService.listBranches(provider, repo),
    enabled: !!provider && !!repo && provider !== 'public',
  });
};

export const useConnectGit = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ConnectGitRequest) => gitService.connect(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['git'] });
    },
  });
};

export const useDisconnectGit = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (provider: string) => gitService.disconnect(provider),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['git'] });
    },
  });
};

export const useConnect = useConnectGit;
export const useDisconnect = useDisconnectGit;
export const useGetStatus = useGitStatus;
