import type { CreateTicketInput, SupportTicket } from '#/interfaces/support';
import { apiClient } from '#/lib/api-client';

export async function getTickets(): Promise<SupportTicket[]> {
  try {
    const res = await apiClient.get<{ tickets: SupportTicket[] }>('/cloud/support/mine');
    return res.tickets || [];
  } catch {
    return [];
  }
}

export async function createTicket(input: CreateTicketInput): Promise<SupportTicket> {
  const res = await apiClient.post<{ ticket: SupportTicket }>('/cloud/support/mine', input);
  return res.ticket;
}

export async function getTicket(id: string): Promise<SupportTicket> {
  const res = await apiClient.get<{ ticket: SupportTicket }>(
    `/cloud/support/mine/${encodeURIComponent(id)}`
  );
  return res.ticket;
}

export async function replyTicket(ticketId: string, message: string): Promise<SupportTicket> {
  const res = await apiClient.post<{ ticket: SupportTicket }>(
    `/cloud/support/mine/${encodeURIComponent(ticketId)}/replies`,
    { message }
  );
  return res.ticket;
}
