import type {
  CreateChatAppInput,
  UpdateChatAppInput,
  VerifyPairingInput,
} from '@codedock/contracts';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { chatAppsService } from '#/services/chat-apps';

export const useChatConnections = () => {
  return useQuery({
    queryKey: ['chat-apps', 'connections'],
    queryFn: async () => {
      const res = await chatAppsService.listConnections();
      return res.data || [];
    },
    refetchInterval: 30_000,
  });
};

export const useChatIdentities = (connectionId?: string) => {
  return useQuery({
    queryKey: ['chat-apps', 'identities', connectionId],
    queryFn: async () => {
      if (!connectionId) return [];
      const res = await chatAppsService.listIdentities(connectionId);
      return res.data || [];
    },
    enabled: Boolean(connectionId),
  });
};

export const usePairingChallenge = (token?: string) => {
  return useQuery({
    queryKey: ['chat-apps', 'challenge', token],
    queryFn: async () => {
      if (!token) return null;
      const res = await chatAppsService.getChallenge(token);
      return res.data;
    },
    enabled: Boolean(token),
  });
};

export const useCreateChatConnection = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateChatAppInput) => {
      const res = await chatAppsService.createConnection(input);
      return res.data;
    },
    onSuccess: (conn) => {
      queryClient.invalidateQueries({ queryKey: ['chat-apps', 'connections'] });
      toast.success(`Connected ${conn.provider} integration`);
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to create chat connection');
    },
  });
};

export const useUpdateChatConnection = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, input }: { id: string; input: UpdateChatAppInput }) => {
      const res = await chatAppsService.updateConnection(id, input);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['chat-apps', 'connections'] });
      toast.success('Chat connection updated');
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to update chat connection');
    },
  });
};

export const useDeleteChatConnection = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await chatAppsService.deleteConnection(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['chat-apps', 'connections'] });
      toast.success('Chat connection removed');
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to delete chat connection');
    },
  });
};

export const useTestChatConnection = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const res = await chatAppsService.testConnection(id);
      return res.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['chat-apps', 'connections'] });
      if (data?.ok) {
        toast.success('Connection test passed!');
      } else {
        toast.error('Connection test failed: check credentials');
      }
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Connection test failed');
    },
  });
};

export const useDeleteChatIdentity = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      connectionId,
      identityId,
    }: {
      connectionId: string;
      identityId: string;
    }) => {
      await chatAppsService.deleteIdentity(identityId);
      return connectionId;
    },
    onSuccess: (connectionId) => {
      queryClient.invalidateQueries({ queryKey: ['chat-apps', 'identities', connectionId] });
      toast.success('Unlinked chat identity');
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to unlink identity');
    },
  });
};

export const useVerifyPairing = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: VerifyPairingInput) => {
      const res = await chatAppsService.verifyPairing(input);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['chat-apps', 'connections'] });
      toast.success('Account successfully paired!');
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Failed to verify pairing code');
    },
  });
};
