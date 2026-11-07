import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateTicketInput } from '#/interfaces/support';
import { createTicket, getTicket, getTickets, replyTicket } from '#/services/cloud-support';

export function useSupportTickets() {
  return useQuery({
    queryKey: ['cloud-support', 'tickets'],
    queryFn: () => getTickets(),
  });
}

export function useSupportTicket(id?: string) {
  return useQuery({
    queryKey: ['cloud-support', 'tickets', id],
    queryFn: () => (id ? getTicket(id) : null),
    enabled: Boolean(id),
  });
}

export function useCreateSupportTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateTicketInput) => createTicket(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cloud-support', 'tickets'] });
    },
  });
}

export function useReplySupportTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ ticketId, message }: { ticketId: string; message: string }) =>
      replyTicket(ticketId, message),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['cloud-support', 'tickets', variables.ticketId],
      });
      queryClient.invalidateQueries({ queryKey: ['cloud-support', 'tickets'] });
    },
  });
}
