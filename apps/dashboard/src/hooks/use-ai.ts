import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useGetAISettings } from '#/hooks/use-settings';
import type { AiChatRequestPayload, AiProviderStatusSummary } from '#/interfaces/ai';
import { aiAssistantService, evaluateAiProviderStatus } from '#/services/ai';

export const useAiResources = () => {
  return useQuery({
    queryKey: ['ai', 'resources'],
    queryFn: () => aiAssistantService.listResources(),
  });
};

export const useSendAiChat = () => {
  return useMutation({
    mutationFn: (payload: AiChatRequestPayload) => aiAssistantService.sendChat(payload),
  });
};

export const useConfirmAiAction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (actionId: string) => aiAssistantService.confirmAction(actionId),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['projects'] }),
        queryClient.invalidateQueries({ queryKey: ['apps'] }),
        queryClient.invalidateQueries({ queryKey: ['deployments'] }),
        queryClient.invalidateQueries({ queryKey: ['databases'] }),
      ]);
    },
  });
};

export const useAiProviderStatus = (): {
  status: AiProviderStatusSummary;
  isLoading: boolean;
} => {
  const { data: settingsRes, isLoading } = useGetAISettings();
  const status = evaluateAiProviderStatus(
    (settingsRes?.data as Record<string, unknown> | undefined) ?? null
  );
  return { status, isLoading };
};
