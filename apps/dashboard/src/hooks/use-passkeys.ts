import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import {
  type AuthenticationResponseJSON,
  browserSupportsWebAuthn,
  startAuthentication,
  startRegistration,
} from '@simplewebauthn/browser';
import { toast } from 'sonner';
import { passkeyService } from '#/services/passkey';
import { useAuthStore } from '#/stores/auth-store';

export const passkeyKeys = {
  all: ['passkeys'] as const,
};

export function usePasskeys() {
  return useQuery({
    queryKey: passkeyKeys.all,
    queryFn: () => passkeyService.listPasskeys(),
  });
}

export function useRegisterPasskey() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ name }: { name?: string } = {}) => {
      if (!browserSupportsWebAuthn()) {
        throw new Error('WebAuthn is not supported in this browser.');
      }
      const options = await passkeyService.getRegisterOptions();
      const response = await startRegistration({ optionsJSON: options });
      return passkeyService.verifyRegistration({
        name: name?.trim() || 'Passkey',
        response,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: passkeyKeys.all });
      toast.success('Passkey added successfully');
    },
    onError: (error: Error) => {
      if (error.name === 'NotAllowedError') {
        toast.info('Passkey registration cancelled');
        return;
      }
      toast.error(error.message || 'Failed to register passkey');
    },
  });
}

export function useRenamePasskey() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, name }: { id: string; name: string }) =>
      passkeyService.renamePasskey(id, name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: passkeyKeys.all });
      toast.success('Passkey renamed');
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to rename passkey');
    },
  });
}

export function useDeletePasskey() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => passkeyService.deletePasskey(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: passkeyKeys.all });
      toast.success('Passkey removed');
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to delete passkey');
    },
  });
}

export function usePasskeyLogin() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: async (args?: { totpCode?: string; assertion?: AuthenticationResponseJSON }) => {
      if (!browserSupportsWebAuthn()) {
        throw new Error('WebAuthn is not supported in this browser.');
      }

      let assertion = args?.assertion;
      if (!assertion) {
        const options = await passkeyService.getAuthenticateOptions();
        assertion = await startAuthentication({ optionsJSON: options });
      }

      try {
        const authData = await passkeyService.verifyAuthentication({
          totpCode: args?.totpCode,
          response: assertion,
        });
        return { authData, assertion, requires2FA: false };
      } catch (err) {
        const message = err instanceof Error ? err.message : '';
        if (message.toLowerCase().includes('2fa code required')) {
          return { authData: null, assertion, requires2FA: true };
        }
        throw err;
      }
    },
    onSuccess: async (result) => {
      if (result.requires2FA || !result.authData) {
        return;
      }
      const data = result.authData;
      if (!data?.token || !data?.user) {
        toast.error('Passkey login failed: invalid response from server');
        return;
      }

      useAuthStore.getState().setAuth(data.token, data.refreshToken || null, data.user);
      queryClient.clear();
      await router.navigate({ to: '/' });
      toast.success('Welcome back!');
    },
    onError: (error: Error) => {
      if (error.name === 'NotAllowedError') {
        toast.info('Passkey authentication cancelled');
        return;
      }
      if (!error.message?.toLowerCase().includes('2fa code required')) {
        toast.error(error.message || 'Passkey sign-in failed');
      }
    },
  });
}
