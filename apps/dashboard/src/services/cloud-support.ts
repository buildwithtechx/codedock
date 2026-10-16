import { apiClient } from '#/lib/api-client';

export type SupportCategory = 'deployment' | 'billing' | 'account' | 'general';
export type SupportStatus = 'open' | 'closed' | 'resolved';

export interface SupportMessage {
  id: string;
  sender: 'user' | 'staff';
  senderName: string;
  content: string;
  createdAt: string;
}

export interface SupportTicket {
  id: string;
  subject: string;
  category: SupportCategory;
  status: SupportStatus;
  hasDiagnosticBundle?: boolean;
  messages: SupportMessage[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateTicketInput {
  subject: string;
  category: SupportCategory;
  message: string;
  attachDiagnostics?: boolean;
}

export async function getTickets(): Promise<SupportTicket[]> {
  try {
    const res = await apiClient.get<{ tickets: SupportTicket[] }>('/api/support/tickets');
    return res.tickets || [];
  } catch {
    return [];
  }
}

export async function createTicket(input: CreateTicketInput): Promise<SupportTicket> {
  const res = await apiClient.post<{ ticket: SupportTicket }>('/api/support/tickets', input);
  return res.ticket;
}

export async function getTicket(id: string): Promise<SupportTicket> {
  const res = await apiClient.get<{ ticket: SupportTicket }>(`/api/support/tickets/${id}`);
  return res.ticket;
}

export async function replyTicket(ticketId: string, message: string): Promise<SupportTicket> {
  const res = await apiClient.post<{ ticket: SupportTicket }>(
    `/api/support/tickets/${ticketId}/reply`,
    { message }
  );
  return res.ticket;
}
